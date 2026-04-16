<script>
  import { auth } from '$lib/stores/auth';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Zap, Search, FileText, ArrowRight, CheckCircle, Star, Users, Send, Shield, Sparkles, Target } from 'lucide-svelte';
</script>

<svelte:head>
  <title>JobApplier - AI-Powered Job Application Platform</title>
  <meta name="description" content="Automate your job search with AI. Generate tailored resumes, craft cover letters, and apply to jobs automatically." />
</svelte:head>

<div class="bg-background">
  <!-- Hero Section -->
  <section class="relative overflow-hidden">
    <!-- Background gradient effects -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-[40%] -right-[20%] w-[600px] h-[600px] rounded-full bg-purple-500/10 blur-3xl"></div>
      <div class="absolute -bottom-[40%] -left-[20%] w-[600px] h-[600px] rounded-full bg-pink-500/8 blur-3xl"></div>
    </div>

    <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-20 pb-16 md:pt-32 md:pb-24">
      <div class="text-center max-w-4xl mx-auto">
        <!-- Badge -->
        <div class="animate-fade-in mb-6">
          <Badge variant="outline" class="px-4 py-1.5 text-sm font-medium border-primary/30 text-primary">
            <Sparkles size={14} class="mr-1.5" />
            AI-Powered Job Applications
          </Badge>
        </div>

        <!-- Headline -->
        <h1 class="animate-fade-in-up text-4xl sm:text-5xl md:text-6xl lg:text-7xl font-extrabold tracking-tight leading-[1.1] mb-6">
          <span class="text-foreground">Land Your Dream Job</span>
          <br />
          <span class="gradient-text animate-gradient">10x Faster with AI</span>
        </h1>

        <!-- Subtitle -->
        <p class="animate-fade-in-up delay-100 text-lg md:text-xl text-muted-foreground max-w-2xl mx-auto mb-10 leading-relaxed">
          Stop spending hours on each application. Let AI generate tailored resumes, craft compelling cover letters, and apply to jobs automatically.
        </p>

        <!-- CTAs -->
        <div class="animate-fade-in-up delay-200 flex flex-col sm:flex-row gap-4 justify-center">
          {#if $auth.isAuthenticated}
            <a href="/dashboard" class="no-underline">
              <Button size="lg" class="w-full sm:w-auto text-base px-8 h-12 shadow-lg shadow-primary/25 hover:shadow-xl hover:shadow-primary/30 transition-all duration-300">
                Go to Dashboard
                <ArrowRight size={18} class="ml-2" />
              </Button>
            </a>
            <a href="/jobs" class="no-underline">
              <Button variant="outline" size="lg" class="w-full sm:w-auto text-base px-8 h-12">
                Browse Jobs
              </Button>
            </a>
          {:else}
            <a href="/register" class="no-underline">
              <Button size="lg" class="w-full sm:w-auto text-base px-8 h-12 shadow-lg shadow-primary/25 hover:shadow-xl hover:shadow-primary/30 transition-all duration-300">
                Get Started Free
                <ArrowRight size={18} class="ml-2" />
              </Button>
            </a>
            <a href="/login" class="no-underline">
              <Button variant="outline" size="lg" class="w-full sm:w-auto text-base px-8 h-12">
                Sign In
              </Button>
            </a>
          {/if}
        </div>

        <!-- Social proof -->
        <div class="animate-fade-in-up delay-300 mt-10 flex items-center justify-center gap-2 text-sm text-muted-foreground">
          <div class="flex -space-x-2">
            {#each [0,1,2,3,4] as i}
              <div class="w-7 h-7 rounded-full bg-primary/20 border-2 border-background flex items-center justify-center text-[10px] font-bold text-primary">
                {['A','J','M','S','K'][i]}
              </div>
            {/each}
          </div>
          <span class="ml-2">Trusted by <strong class="text-foreground">2,000+</strong> job seekers</span>
        </div>
      </div>
    </div>
  </section>

  <!-- Stats Bar -->
  <section class="border-y border-border/50 bg-muted/30">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-6 md:gap-8">
        <div class="text-center">
          <div class="text-2xl md:text-3xl font-bold text-foreground">10,000+</div>
          <div class="text-sm text-muted-foreground mt-1">Applications Sent</div>
        </div>
        <div class="text-center">
          <div class="text-2xl md:text-3xl font-bold text-foreground">500+</div>
          <div class="text-sm text-muted-foreground mt-1">Companies</div>
        </div>
        <div class="text-center">
          <div class="text-2xl md:text-3xl font-bold text-foreground">95%</div>
          <div class="text-sm text-muted-foreground mt-1">Satisfaction Rate</div>
        </div>
        <div class="text-center">
          <div class="text-2xl md:text-3xl font-bold text-foreground">3x</div>
          <div class="text-sm text-muted-foreground mt-1">Faster Hiring</div>
        </div>
      </div>
    </div>
  </section>

  <!-- Features Section -->
  <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 md:py-28">
    <div class="text-center mb-16">
      <Badge variant="outline" class="px-3 py-1 text-xs font-medium border-primary/30 text-primary mb-4">Features</Badge>
      <h2 class="text-3xl md:text-4xl font-bold text-foreground mb-4">
        Everything you need to land your next role
      </h2>
      <p class="text-lg text-muted-foreground max-w-2xl mx-auto">
        From intelligent job matching to AI-crafted applications, we handle the heavy lifting so you can focus on what matters.
      </p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 lg:gap-8">
      <Card class="group hover:shadow-lg hover:border-primary/20 transition-all duration-300 hover:-translate-y-1">
        <CardHeader class="pb-4">
          <div class="w-12 h-12 rounded-xl bg-blue-500/10 flex items-center justify-center mb-4 group-hover:bg-blue-500/20 transition-colors">
            <Search size={22} class="text-blue-500" />
          </div>
          <CardTitle class="text-lg">Smart Job Search</CardTitle>
          <CardDescription class="text-sm leading-relaxed">
            Search across multiple platforms with intelligent filters. Find the perfect match for your skills and preferences.
          </CardDescription>
        </CardHeader>
      </Card>

      <Card class="group hover:shadow-lg hover:border-primary/20 transition-all duration-300 hover:-translate-y-1">
        <CardHeader class="pb-4">
          <div class="w-12 h-12 rounded-xl bg-purple-500/10 flex items-center justify-center mb-4 group-hover:bg-purple-500/20 transition-colors">
            <FileText size={22} class="text-purple-500" />
          </div>
          <CardTitle class="text-lg">AI Resume & Cover Letters</CardTitle>
          <CardDescription class="text-sm leading-relaxed">
            Generate tailored resumes and compelling cover letters for each application using advanced AI models.
          </CardDescription>
        </CardHeader>
      </Card>

      <Card class="group hover:shadow-lg hover:border-primary/20 transition-all duration-300 hover:-translate-y-1">
        <CardHeader class="pb-4">
          <div class="w-12 h-12 rounded-xl bg-amber-500/10 flex items-center justify-center mb-4 group-hover:bg-amber-500/20 transition-colors">
            <Zap size={22} class="text-amber-500" />
          </div>
          <CardTitle class="text-lg">One-Click Apply</CardTitle>
          <CardDescription class="text-sm leading-relaxed">
            Apply to multiple jobs with a single click. Track all your applications in one unified dashboard.
          </CardDescription>
        </CardHeader>
      </Card>

      <Card class="group hover:shadow-lg hover:border-primary/20 transition-all duration-300 hover:-translate-y-1">
        <CardHeader class="pb-4">
          <div class="w-12 h-12 rounded-xl bg-green-500/10 flex items-center justify-center mb-4 group-hover:bg-green-500/20 transition-colors">
            <Target size={22} class="text-green-500" />
          </div>
          <CardTitle class="text-lg">Intelligent Matching</CardTitle>
          <CardDescription class="text-sm leading-relaxed">
            Our AI analyzes your profile and matches you with jobs where you are most likely to succeed.
          </CardDescription>
        </CardHeader>
      </Card>

      <Card class="group hover:shadow-lg hover:border-primary/20 transition-all duration-300 hover:-translate-y-1">
        <CardHeader class="pb-4">
          <div class="w-12 h-12 rounded-xl bg-pink-500/10 flex items-center justify-center mb-4 group-hover:bg-pink-500/20 transition-colors">
            <Send size={22} class="text-pink-500" />
          </div>
          <CardTitle class="text-lg">Application Tracking</CardTitle>
          <CardDescription class="text-sm leading-relaxed">
            Monitor every application from submission to offer. Never lose track of where you stand.
          </CardDescription>
        </CardHeader>
      </Card>

      <Card class="group hover:shadow-lg hover:border-primary/20 transition-all duration-300 hover:-translate-y-1">
        <CardHeader class="pb-4">
          <div class="w-12 h-12 rounded-xl bg-red-500/10 flex items-center justify-center mb-4 group-hover:bg-red-500/20 transition-colors">
            <Shield size={22} class="text-red-500" />
          </div>
          <CardTitle class="text-lg">Secure & Private</CardTitle>
          <CardDescription class="text-sm leading-relaxed">
            Your data is encrypted with AES-256-GCM. We never share your information with third parties.
          </CardDescription>
        </CardHeader>
      </Card>
    </div>
  </section>

  <!-- How It Works -->
  <section class="bg-muted/30 border-y border-border/50 py-20 md:py-28">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="text-center mb-16">
        <Badge variant="outline" class="px-3 py-1 text-xs font-medium border-primary/30 text-primary mb-4">How It Works</Badge>
        <h2 class="text-3xl md:text-4xl font-bold text-foreground mb-4">
          Three steps to your next role
        </h2>
        <p class="text-lg text-muted-foreground max-w-2xl mx-auto">
          Getting started takes less than 5 minutes. Set up once and let AI do the rest.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-5xl mx-auto">
        <div class="relative text-center">
          <div class="w-14 h-14 rounded-2xl bg-primary text-primary-foreground flex items-center justify-center mx-auto mb-6 font-bold text-xl shadow-lg shadow-primary/25">
            1
          </div>
          <h3 class="text-lg font-semibold text-foreground mb-3">Set Up Your Profile</h3>
          <p class="text-sm text-muted-foreground leading-relaxed">
            Add your experience, skills, education, and preferences to personalize every application.
          </p>
        </div>

        <div class="relative text-center">
          <div class="w-14 h-14 rounded-2xl bg-primary text-primary-foreground flex items-center justify-center mx-auto mb-6 font-bold text-xl shadow-lg shadow-primary/25">
            2
          </div>
          <h3 class="text-lg font-semibold text-foreground mb-3">Configure AI Settings</h3>
          <p class="text-sm text-muted-foreground leading-relaxed">
            Choose your preferred AI model, set job search filters, and customize your preferences.
          </p>
        </div>

        <div class="relative text-center">
          <div class="w-14 h-14 rounded-2xl bg-primary text-primary-foreground flex items-center justify-center mx-auto mb-6 font-bold text-xl shadow-lg shadow-primary/25">
            3
          </div>
          <h3 class="text-lg font-semibold text-foreground mb-3">Apply Automatically</h3>
          <p class="text-sm text-muted-foreground leading-relaxed">
            Search for jobs, and let AI generate custom resumes and cover letters for each application.
          </p>
        </div>
      </div>
    </div>
  </section>

  <!-- Testimonials -->
  <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 md:py-28">
    <div class="text-center mb-16">
      <Badge variant="outline" class="px-3 py-1 text-xs font-medium border-primary/30 text-primary mb-4">Testimonials</Badge>
      <h2 class="text-3xl md:text-4xl font-bold text-foreground mb-4">
        Loved by job seekers everywhere
      </h2>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 lg:gap-8">
      <Card class="hover:shadow-lg transition-all duration-300">
        <CardContent class="p-6">
          <div class="flex gap-1 mb-4">
            {#each [1,2,3,4,5] as _}
              <Star size={16} class="fill-amber-400 text-amber-400" />
            {/each}
          </div>
          <p class="text-sm text-muted-foreground leading-relaxed mb-4">
            "I applied to 50 jobs in one weekend. The AI-generated cover letters were so personalized that I got callbacks from 12 companies. Landed my dream role in 3 weeks!"
          </p>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-full bg-blue-500/10 flex items-center justify-center text-sm font-bold text-blue-500">
              SM
            </div>
            <div>
              <div class="text-sm font-semibold text-foreground">Sarah M.</div>
              <div class="text-xs text-muted-foreground">Software Engineer at Google</div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card class="hover:shadow-lg transition-all duration-300">
        <CardContent class="p-6">
          <div class="flex gap-1 mb-4">
            {#each [1,2,3,4,5] as _}
              <Star size={16} class="fill-amber-400 text-amber-400" />
            {/each}
          </div>
          <p class="text-sm text-muted-foreground leading-relaxed mb-4">
            "The resume tailoring feature is incredible. Each resume was customized for the specific role. My response rate went from 5% to 35%. Absolutely game-changing."
          </p>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-full bg-green-500/10 flex items-center justify-center text-sm font-bold text-green-500">
              AK
            </div>
            <div>
              <div class="text-sm font-semibold text-foreground">Ahmed K.</div>
              <div class="text-xs text-muted-foreground">Product Manager at Meta</div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card class="hover:shadow-lg transition-all duration-300">
        <CardContent class="p-6">
          <div class="flex gap-1 mb-4">
            {#each [1,2,3,4,5] as _}
              <Star size={16} class="fill-amber-400 text-amber-400" />
            {/each}
          </div>
          <p class="text-sm text-muted-foreground leading-relaxed mb-4">
            "As a career changer, I struggled with cover letters. JobApplier's AI understood my transferable skills and presented them perfectly. Got 3 offers in a month!"
          </p>
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-full bg-purple-500/10 flex items-center justify-center text-sm font-bold text-purple-500">
              LR
            </div>
            <div>
              <div class="text-sm font-semibold text-foreground">Lisa R.</div>
              <div class="text-xs text-muted-foreground">UX Designer at Stripe</div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </section>

  <!-- Pricing Preview -->
  <section class="bg-muted/30 border-y border-border/50 py-20 md:py-28">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="text-center mb-16">
        <Badge variant="outline" class="px-3 py-1 text-xs font-medium border-primary/30 text-primary mb-4">Pricing</Badge>
        <h2 class="text-3xl md:text-4xl font-bold text-foreground mb-4">
          Simple, transparent pricing
        </h2>
        <p class="text-lg text-muted-foreground max-w-2xl mx-auto">
          Start free and scale as you grow. No hidden fees, cancel anytime.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 lg:gap-8 max-w-5xl mx-auto">
        <!-- Free -->
        <Card class="hover:shadow-lg transition-all duration-300">
          <CardContent class="p-8">
            <h3 class="text-lg font-semibold text-foreground mb-2">Free</h3>
            <div class="flex items-baseline gap-1 mb-1">
              <span class="text-4xl font-bold text-foreground">$0</span>
              <span class="text-muted-foreground text-sm">/month</span>
            </div>
            <p class="text-sm text-muted-foreground mb-6">Perfect for getting started</p>
            <a href="/register" class="no-underline">
              <Button variant="outline" class="w-full mb-6">Get Started</Button>
            </a>
            <div class="space-y-3">
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <CheckCircle size={16} class="text-green-500 shrink-0" />
                5 applications / month
              </div>
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <CheckCircle size={16} class="text-green-500 shrink-0" />
                Basic AI resume generation
              </div>
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <CheckCircle size={16} class="text-green-500 shrink-0" />
                Job search & tracking
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Pro -->
        <Card class="relative border-primary shadow-lg shadow-primary/10 hover:shadow-xl transition-all duration-300">
          <div class="absolute -top-3 left-1/2 -translate-x-1/2">
            <Badge class="px-3 py-1 bg-primary text-primary-foreground border-transparent">Most Popular</Badge>
          </div>
          <CardContent class="p-8">
            <h3 class="text-lg font-semibold text-foreground mb-2">Pro</h3>
            <div class="flex items-baseline gap-1 mb-1">
              <span class="text-4xl font-bold text-foreground">$19</span>
              <span class="text-muted-foreground text-sm">/month</span>
            </div>
            <p class="text-sm text-muted-foreground mb-6">For active job seekers</p>
            <a href="/register" class="no-underline">
              <Button class="w-full mb-6">Start Pro Trial</Button>
            </a>
            <div class="space-y-3">
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <CheckCircle size={16} class="text-green-500 shrink-0" />
                Unlimited applications
              </div>
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <CheckCircle size={16} class="text-green-500 shrink-0" />
                Advanced AI (GPT-4, Claude)
              </div>
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <CheckCircle size={16} class="text-green-500 shrink-0" />
                Custom cover letters
              </div>
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <CheckCircle size={16} class="text-green-500 shrink-0" />
                Priority support
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Enterprise -->
        <Card class="hover:shadow-lg transition-all duration-300">
          <CardContent class="p-8">
            <h3 class="text-lg font-semibold text-foreground mb-2">Enterprise</h3>
            <div class="flex items-baseline gap-1 mb-1">
              <span class="text-4xl font-bold text-foreground">$49</span>
              <span class="text-muted-foreground text-sm">/month</span>
            </div>
            <p class="text-sm text-muted-foreground mb-6">For teams and agencies</p>
            <a href="/register" class="no-underline">
              <Button variant="outline" class="w-full mb-6">Contact Sales</Button>
            </a>
            <div class="space-y-3">
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <CheckCircle size={16} class="text-green-500 shrink-0" />
                Everything in Pro
              </div>
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <CheckCircle size={16} class="text-green-500 shrink-0" />
                Team collaboration
              </div>
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <CheckCircle size={16} class="text-green-500 shrink-0" />
                API access
              </div>
              <div class="flex items-center gap-2 text-sm text-muted-foreground">
                <CheckCircle size={16} class="text-green-500 shrink-0" />
                Dedicated account manager
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </section>

  <!-- Final CTA -->
  <section class="relative overflow-hidden py-20 md:py-28">
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[400px] rounded-full bg-primary/5 blur-3xl"></div>
    </div>
    <div class="relative max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
      <h2 class="text-3xl md:text-4xl font-bold text-foreground mb-4">
        Ready to transform your job search?
      </h2>
      <p class="text-lg text-muted-foreground mb-8 max-w-xl mx-auto leading-relaxed">
        Join thousands of job seekers who are landing interviews faster with AI-powered applications.
      </p>
      {#if !$auth.isAuthenticated}
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <a href="/register" class="no-underline">
            <Button size="lg" class="w-full sm:w-auto text-base px-8 h-12 shadow-lg shadow-primary/25 hover:shadow-xl hover:shadow-primary/30 transition-all duration-300">
              Create Your Free Account
              <ArrowRight size={18} class="ml-2" />
            </Button>
          </a>
        </div>
        <p class="text-sm text-muted-foreground mt-4">No credit card required. Free plan available.</p>
      {:else}
        <a href="/dashboard" class="no-underline">
          <Button size="lg" class="text-base px-8 h-12 shadow-lg shadow-primary/25">
            Go to Dashboard
            <ArrowRight size={18} class="ml-2" />
          </Button>
        </a>
      {/if}
    </div>
  </section>
</div>
