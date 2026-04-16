package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v81"
	billingportalsession "github.com/stripe/stripe-go/v81/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/subscription"
	"github.com/stripe/stripe-go/v81/webhook"

	"job-applier-backend/internal/models"
)

// CreateCheckoutSession creates a Stripe Checkout Session for upgrading to Pro or Enterprise.
// POST /api/v1/payments/checkout
func (h *Handlers) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)

	var input struct {
		Tier string `json:"tier"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Determine the price ID based on the tier.
	var priceID string
	switch input.Tier {
	case "pro":
		priceID = h.cfg.StripePriceProMonthly
	case "enterprise":
		priceID = h.cfg.StripePriceEnterpriseMonthly
	default:
		h.respondError(w, http.StatusBadRequest, "Invalid tier. Must be 'pro' or 'enterprise'")
		return
	}

	if priceID == "" {
		h.respondError(w, http.StatusInternalServerError, "Stripe price ID not configured for this tier")
		return
	}

	// Get the user's email for creating/finding the Stripe customer.
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		h.respondError(w, http.StatusNotFound, "User not found")
		return
	}

	// Get or create the Stripe customer ID.
	stripeCustomerID, err := h.getOrCreateStripeCustomer(userID, user.Email)
	if err != nil {
		log.Printf("Failed to get/create Stripe customer for user %d: %v", userID, err)
		h.respondError(w, http.StatusInternalServerError, "Failed to set up payment")
		return
	}

	// Create the Checkout Session.
	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(stripeCustomerID),
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(r.Header.Get("Origin") + "/settings?payment=success"),
		CancelURL:  stripe.String(r.Header.Get("Origin") + "/settings?payment=canceled"),
		Metadata: map[string]string{
			"user_id": formatUint(userID),
			"tier":    input.Tier,
		},
	}

	session, err := checkoutsession.New(params)
	if err != nil {
		log.Printf("Failed to create Stripe checkout session: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to create checkout session")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"checkout_url": session.URL,
	})
}

// HandleStripeWebhook processes Stripe webhook events.
// POST /api/v1/payments/webhook
// This endpoint must be PUBLIC (no auth middleware) since Stripe sends webhooks without JWT.
func (h *Handlers) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	const maxBodyBytes = 65536
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	// Verify webhook signature to prevent spoofing.
	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), h.cfg.StripeWebhookSecret)
	if err != nil {
		log.Printf("Stripe webhook signature verification failed: %v", err)
		h.respondError(w, http.StatusBadRequest, "Invalid webhook signature")
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		h.handleCheckoutCompleted(event)
	case "customer.subscription.updated":
		h.handleSubscriptionUpdated(event)
	case "customer.subscription.deleted":
		h.handleSubscriptionDeleted(event)
	case "invoice.payment_failed":
		h.handlePaymentFailed(event)
	default:
		log.Printf("Unhandled Stripe webhook event type: %s", event.Type)
	}

	// Always return 200 to acknowledge receipt.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GetSubscription returns the user's current subscription details.
// GET /api/v1/payments/subscription
func (h *Handlers) GetSubscription(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)

	var sub models.Subscription
	if err := h.db.Where("user_id = ?", userID).First(&sub).Error; err != nil {
		// No subscription found; return default free tier.
		h.respondJSON(w, http.StatusOK, map[string]interface{}{
			"tier":   "free",
			"status": "active",
		})
		return
	}

	h.respondJSON(w, http.StatusOK, sub)
}

// CreateBillingPortal creates a Stripe Billing Portal session for the user.
// GET /api/v1/payments/portal
func (h *Handlers) CreateBillingPortal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)

	var sub models.Subscription
	if err := h.db.Where("user_id = ?", userID).First(&sub).Error; err != nil || sub.StripeCustomerID == "" {
		h.respondError(w, http.StatusBadRequest, "No active subscription found. Subscribe first.")
		return
	}

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(sub.StripeCustomerID),
		ReturnURL: stripe.String(r.Header.Get("Origin") + "/settings"),
	}

	session, err := billingportalsession.New(params)
	if err != nil {
		log.Printf("Failed to create billing portal session: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to create billing portal session")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{
		"portal_url": session.URL,
	})
}

// CancelSubscription cancels the user's subscription at the end of the current period.
// POST /api/v1/payments/cancel
func (h *Handlers) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(uint)

	var sub models.Subscription
	if err := h.db.Where("user_id = ?", userID).First(&sub).Error; err != nil || sub.StripeSubscriptionID == "" {
		h.respondError(w, http.StatusBadRequest, "No active subscription found")
		return
	}

	// Cancel at period end so the user retains access until the period expires.
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	updatedSub, err := subscription.Update(sub.StripeSubscriptionID, params)
	if err != nil {
		log.Printf("Failed to cancel Stripe subscription: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to cancel subscription")
		return
	}

	// Update local record.
	sub.Status = "canceled"
	sub.CurrentPeriodEnd = time.Unix(updatedSub.CurrentPeriodEnd, 0)
	h.db.Save(&sub)

	h.respondJSON(w, http.StatusOK, sub)
}

// --- Internal helpers ---

// getOrCreateStripeCustomer retrieves the existing Stripe customer ID from the local
// Subscription record, or creates a new Stripe customer and saves the ID.
func (h *Handlers) getOrCreateStripeCustomer(userID uint, email string) (string, error) {
	var sub models.Subscription
	err := h.db.Where("user_id = ?", userID).First(&sub).Error

	if err == nil && sub.StripeCustomerID != "" {
		return sub.StripeCustomerID, nil
	}

	// Create a new Stripe customer.
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Metadata: map[string]string{
			"user_id": formatUint(userID),
		},
	}

	cust, err := customer.New(params)
	if err != nil {
		return "", err
	}

	// Ensure we have a local Subscription record to store the customer ID.
	if sub.ID == 0 {
		sub = models.Subscription{
			UserID:           userID,
			StripeCustomerID: cust.ID,
			Tier:             "free",
			Status:           "active",
		}
		h.db.Create(&sub)
	} else {
		sub.StripeCustomerID = cust.ID
		h.db.Save(&sub)
	}

	return cust.ID, nil
}

// handleCheckoutCompleted processes the checkout.session.completed event.
func (h *Handlers) handleCheckoutCompleted(event stripe.Event) {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		log.Printf("Failed to unmarshal checkout session: %v", err)
		return
	}

	userIDStr, ok := session.Metadata["user_id"]
	if !ok {
		log.Println("checkout.session.completed: missing user_id in metadata")
		return
	}

	tier, _ := session.Metadata["tier"]
	if tier == "" {
		tier = "pro"
	}

	userID := parseUint(userIDStr)
	if userID == 0 {
		log.Printf("checkout.session.completed: invalid user_id: %s", userIDStr)
		return
	}

	// Retrieve the full subscription from Stripe to get period dates.
	stripeSub, err := subscription.Get(session.Subscription.ID, nil)
	if err != nil {
		log.Printf("Failed to retrieve Stripe subscription %s: %v", session.Subscription.ID, err)
		return
	}

	var sub models.Subscription
	h.db.Where("user_id = ?", userID).FirstOrCreate(&sub, models.Subscription{UserID: userID})

	sub.StripeCustomerID = session.Customer.ID
	sub.StripeSubscriptionID = session.Subscription.ID
	sub.Tier = tier
	sub.Status = "active"
	sub.CurrentPeriodStart = time.Unix(stripeSub.CurrentPeriodStart, 0)
	sub.CurrentPeriodEnd = time.Unix(stripeSub.CurrentPeriodEnd, 0)

	h.db.Save(&sub)
	log.Printf("Subscription created/updated for user %d: tier=%s", userID, tier)
}

// handleSubscriptionUpdated processes the customer.subscription.updated event.
func (h *Handlers) handleSubscriptionUpdated(event stripe.Event) {
	var stripeSub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &stripeSub); err != nil {
		log.Printf("Failed to unmarshal subscription: %v", err)
		return
	}

	var sub models.Subscription
	if err := h.db.Where("stripe_subscription_id = ?", stripeSub.ID).First(&sub).Error; err != nil {
		log.Printf("Subscription not found for stripe_subscription_id %s", stripeSub.ID)
		return
	}

	// Determine tier from the price ID of the first item.
	if len(stripeSub.Items.Data) > 0 {
		priceID := stripeSub.Items.Data[0].Price.ID
		switch priceID {
		case h.cfg.StripePriceProMonthly:
			sub.Tier = "pro"
		case h.cfg.StripePriceEnterpriseMonthly:
			sub.Tier = "enterprise"
		}
	}

	sub.Status = string(stripeSub.Status)
	sub.CurrentPeriodStart = time.Unix(stripeSub.CurrentPeriodStart, 0)
	sub.CurrentPeriodEnd = time.Unix(stripeSub.CurrentPeriodEnd, 0)

	h.db.Save(&sub)
	log.Printf("Subscription updated for user %d: tier=%s status=%s", sub.UserID, sub.Tier, sub.Status)
}

// handleSubscriptionDeleted processes the customer.subscription.deleted event.
func (h *Handlers) handleSubscriptionDeleted(event stripe.Event) {
	var stripeSub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &stripeSub); err != nil {
		log.Printf("Failed to unmarshal subscription: %v", err)
		return
	}

	var sub models.Subscription
	if err := h.db.Where("stripe_subscription_id = ?", stripeSub.ID).First(&sub).Error; err != nil {
		log.Printf("Subscription not found for stripe_subscription_id %s", stripeSub.ID)
		return
	}

	sub.Tier = "free"
	sub.Status = "canceled"

	h.db.Save(&sub)
	log.Printf("Subscription canceled for user %d", sub.UserID)
}

// handlePaymentFailed processes the invoice.payment_failed event.
func (h *Handlers) handlePaymentFailed(event stripe.Event) {
	var invoice struct {
		Subscription string `json:"subscription"`
	}
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		log.Printf("Failed to unmarshal invoice: %v", err)
		return
	}

	if invoice.Subscription == "" {
		return
	}

	var sub models.Subscription
	if err := h.db.Where("stripe_subscription_id = ?", invoice.Subscription).First(&sub).Error; err != nil {
		log.Printf("Subscription not found for stripe_subscription_id %s", invoice.Subscription)
		return
	}

	sub.Status = "past_due"
	h.db.Save(&sub)
	log.Printf("Subscription marked past_due for user %d", sub.UserID)
}

// formatUint converts a uint to its string representation.
func formatUint(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}

// parseUint converts a string to uint, returning 0 on failure.
func parseUint(s string) uint {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint(n)
}
