<script>
  import { auth } from '$lib/stores/auth';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription } from '$lib/components/ui/card';
</script>

<svelte:head>
  <title>JobApplier - AI-Powered Job Application Platform</title>
  <meta name="description" content="Automate your job search with AI. Generate tailored resumes, craft cover letters, and apply to jobs automatically." />
</svelte:head>

<div class="min-h-[80vh] bg-background">
  <!-- Hero Section -->
  <section class="text-center py-16 md:py-24 px-4">
    <h1 class="text-4xl md:text-6xl font-bold tracking-tight text-foreground mb-4">
      <span class="bg-gradient-to-r from-primary via-purple-500 to-pink-500 bg-clip-text text-transparent">
        Automate Your Job Search
      </span>
    </h1>
    <p class="text-lg md:text-xl text-muted-foreground max-w-2xl mx-auto mb-8">
      Let AI handle the tedious parts of job hunting. Generate tailored resumes, craft cover letters, and apply to jobs automatically.
    </p>
    <div class="flex flex-col sm:flex-row gap-4 justify-center">
      {#if $auth.isAuthenticated}
        <a href="/dashboard">
          <Button size="lg" class="w-full sm:w-auto text-base px-8">Go to Dashboard</Button>
        </a>
        <a href="/jobs">
          <Button variant="outline" size="lg" class="w-full sm:w-auto text-base px-8">Browse Jobs</Button>
        </a>
      {:else}
        <a href="/register">
          <Button size="lg" class="w-full sm:w-auto text-base px-8">Get Started Free</Button>
        </a>
        <a href="/login">
          <Button variant="outline" size="lg" class="w-full sm:w-auto text-base px-8">Sign In</Button>
        </a>
      {/if}
    </div>
  </section>

  <!-- Features Section -->
  <section class="max-w-6xl mx-auto px-4 pb-16 md:pb-24">
    <h2 class="text-2xl md:text-3xl font-bold text-foreground text-center mb-12">
      Everything you need to land your next role
    </h2>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <Card class="text-center">
        <CardHeader>
          <div class="mx-auto w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center mb-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
          <CardTitle class="text-lg">Smart Job Search</CardTitle>
          <CardDescription>Search across multiple platforms with intelligent filters for location, experience level, and work type.</CardDescription>
        </CardHeader>
      </Card>

      <Card class="text-center">
        <CardHeader>
          <div class="mx-auto w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center mb-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <CardTitle class="text-lg">AI Resume & Cover Letters</CardTitle>
          <CardDescription>Generate tailored resumes and cover letters for each application using advanced AI models.</CardDescription>
        </CardHeader>
      </Card>

      <Card class="text-center">
        <CardHeader>
          <div class="mx-auto w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center mb-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <CardTitle class="text-lg">One-Click Apply</CardTitle>
          <CardDescription>Apply to multiple jobs with a single click. Track all your applications in one dashboard.</CardDescription>
        </CardHeader>
      </Card>
    </div>
  </section>

  <!-- How It Works Section -->
  <section class="bg-muted/50 py-16 md:py-24 px-4">
    <div class="max-w-4xl mx-auto text-center">
      <h2 class="text-2xl md:text-3xl font-bold text-foreground mb-12">How It Works</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div>
          <div class="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center mx-auto mb-4 font-bold text-lg">1</div>
          <h3 class="font-semibold text-foreground mb-2">Set Up Your Profile</h3>
          <p class="text-sm text-muted-foreground">Add your experience, skills, and preferences to personalize your applications.</p>
        </div>
        <div>
          <div class="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center mx-auto mb-4 font-bold text-lg">2</div>
          <h3 class="font-semibold text-foreground mb-2">Configure AI Settings</h3>
          <p class="text-sm text-muted-foreground">Choose your preferred AI model and set job search filters.</p>
        </div>
        <div>
          <div class="w-10 h-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center mx-auto mb-4 font-bold text-lg">3</div>
          <h3 class="font-semibold text-foreground mb-2">Apply Automatically</h3>
          <p class="text-sm text-muted-foreground">Let the AI generate custom documents and apply to matching jobs.</p>
        </div>
      </div>
    </div>
  </section>

  <!-- CTA Section -->
  <section class="text-center py-16 md:py-24 px-4">
    <h2 class="text-2xl md:text-3xl font-bold text-foreground mb-4">Ready to get started?</h2>
    <p class="text-muted-foreground mb-8 max-w-lg mx-auto">Join thousands of job seekers who are automating their job search with AI.</p>
    {#if !$auth.isAuthenticated}
      <a href="/register">
        <Button size="lg" class="text-base px-8">Create Your Free Account</Button>
      </a>
    {:else}
      <a href="/dashboard">
        <Button size="lg" class="text-base px-8">Go to Dashboard</Button>
      </a>
    {/if}
  </section>
</div>
