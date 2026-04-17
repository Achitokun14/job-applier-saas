package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"job-applier-backend/internal/models"
)

// RefreshTokenHandler handles token refresh requests.
// POST /api/v1/auth/refresh
func (h *Handlers) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.RefreshToken == "" {
		h.respondError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	// Find the refresh token in DB
	var storedToken models.RefreshToken
	if err := h.db.Where("token = ?", input.RefreshToken).First(&storedToken).Error; err != nil {
		h.respondError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	// Check if revoked
	if storedToken.Revoked {
		h.respondError(w, http.StatusUnauthorized, "Refresh token has been revoked")
		return
	}

	// Check if expired
	if time.Now().After(storedToken.ExpiresAt) {
		h.respondError(w, http.StatusUnauthorized, "Refresh token has expired")
		return
	}

	// Revoke the old refresh token (rotation)
	if err := h.db.Model(&storedToken).Update("revoked", true).Error; err != nil {
		log.Printf("ERROR: Failed to revoke refresh token %d: %v", storedToken.ID, err)
		h.respondError(w, http.StatusInternalServerError, "Failed to revoke old refresh token")
		return
	}

	// Generate new access token
	accessToken, err := h.generateAccessToken(storedToken.UserID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	// Generate new refresh token
	newRefreshToken, err := h.generateRefreshToken(storedToken.UserID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"token":         accessToken,
		"refresh_token": newRefreshToken,
		"expires_in":    900,
	})
}

// Logout handles user logout by blacklisting the current JWT and revoking all refresh tokens.
// POST /api/v1/auth/logout
func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uint)
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Extract the token to get the jti claim
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.cfg.JWTSecret), nil
	})
	if err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if jti, ok := claims["jti"].(string); ok && jti != "" {
				// Calculate remaining expiry for blacklist TTL
				if exp, ok := claims["exp"].(float64); ok {
					remaining := time.Until(time.Unix(int64(exp), 0))
					if remaining > 0 {
						h.blacklist.Blacklist(jti, remaining)
					}
				}
			}
		}
	}

	// Revoke all refresh tokens for this user
	h.db.Model(&models.RefreshToken{}).Where("user_id = ? AND revoked = ?", userID, false).Update("revoked", true)

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// RequestPasswordReset handles forgot-password requests.
// POST /api/v1/auth/forgot-password
func (h *Handlers) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Always return 200 to not reveal if email exists
	defer h.respondJSON(w, http.StatusOK, map[string]string{"message": "If the email exists, a password reset link has been sent"})

	if input.Email == "" {
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return
	}

	// Generate a random reset token
	rawToken, err := h.generateRandomToken()
	if err != nil {
		log.Printf("Failed to generate password reset token: %v", err)
		return
	}

	// Store the hash of the token
	tokenHash := hashToken(rawToken)

	resetToken := models.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := h.db.Create(&resetToken).Error; err != nil {
		log.Printf("Failed to store password reset token: %v", err)
		return
	}

	// Log the reset URL (actual email sending is a separate task)
	resetURL := fmt.Sprintf("/reset-password?token=%s", rawToken)
	log.Printf("Password reset requested for user %d. Reset URL: %s", user.ID, resetURL)
}

// ResetPassword handles the actual password reset.
// POST /api/v1/auth/reset-password
func (h *Handlers) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.Token == "" || input.NewPassword == "" {
		h.respondError(w, http.StatusBadRequest, "Token and new password are required")
		return
	}

	// Find the reset token by hash
	tokenHash := hashToken(input.Token)
	var resetToken models.PasswordResetToken
	if err := h.db.Where("token_hash = ? AND used = ?", tokenHash, false).First(&resetToken).Error; err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid or expired reset token")
		return
	}

	// Check if expired
	if time.Now().After(resetToken.ExpiresAt) {
		h.respondError(w, http.StatusBadRequest, "Reset token has expired")
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Update the user's password
	if err := h.db.Model(&models.User{}).Where("id = ?", resetToken.UserID).Update("password", string(hashedPassword)).Error; err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	// Mark the reset token as used
	h.db.Model(&resetToken).Update("used", true)

	// Revoke all refresh tokens for this user
	h.db.Model(&models.RefreshToken{}).Where("user_id = ? AND revoked = ?", resetToken.UserID, false).Update("revoked", true)

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Password reset successfully"})
}
