<script>
  import { auth } from '$lib/stores/auth';
  import { api } from '$lib/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';
  import { User, Briefcase, GraduationCap, Code, FolderOpen, CheckCircle, AlertCircle, Loader2, Save } from 'lucide-svelte';

  let profile = $state({});
  let name = $state('');
  let personalInfo = $state('');
  let education = $state('');
  let experience = $state('');
  let skills = $state('');
  let projects = $state('');
  let loading = $state(true);
  let saving = $state(false);
  let message = $state('');
  let messageType = $state('success');

  onMount(async () => {
    if (!$auth.isAuthenticated) {
      goto('/login');
      return;
    }

    try {
      const token = $auth.token;
      profile = await api.getProfile(token);
      name = profile.name || '';
      personalInfo = profile.personalInfo || '';
      education = profile.education || '';
      experience = profile.experience || '';
      skills = profile.skills || '';
      projects = profile.projects || '';
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  });

  async function saveProfile() {
    saving = true;
    message = '';
    try {
      const token = $auth.token;
      await api.updateProfile({
        name,
        personalInfo,
        education,
        experience,
        skills,
        projects
      }, token);
      message = 'Profile saved successfully!';
      messageType = 'success';
      setTimeout(() => { message = ''; }, 3000);
    } catch (e) {
      message = 'Failed to save profile. Please try again.';
      messageType = 'error';
    } finally {
      saving = false;
    }
  }
</script>

<svelte:head>
  <title>My Profile - JobApplier</title>
</svelte:head>

<div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
  <!-- Header -->
  <div class="mb-8 animate-fade-in">
    <h1 class="text-2xl sm:text-3xl font-bold text-foreground">My Profile</h1>
    <p class="text-muted-foreground mt-1">Manage your personal and professional details for better applications</p>
  </div>

  {#if loading}
    <div class="space-y-6 animate-fade-in">
      {#each [1,2,3,4] as _}
        <div class="skeleton h-48 rounded-xl"></div>
      {/each}
    </div>
  {:else}
    <form onsubmit={(e) => { e.preventDefault(); saveProfile(); }} class="space-y-6">
      <!-- Personal Info -->
      <Card class="animate-fade-in-up">
        <CardHeader class="pb-4">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-lg bg-primary/10 flex items-center justify-center">
              <User size={16} class="text-primary" />
            </div>
            <div>
              <CardTitle class="text-base">Personal Information</CardTitle>
              <CardDescription class="text-xs">Your name and personal statement</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <label for="name" class="text-sm font-medium text-foreground">Full Name</label>
            <Input type="text" id="name" bind:value={name} placeholder="John Doe" />
          </div>

          <div class="space-y-2">
            <label for="personalInfo" class="text-sm font-medium text-foreground">Personal Statement</label>
            <textarea
              id="personalInfo"
              bind:value={personalInfo}
              rows="4"
              placeholder="A brief summary about yourself, your career goals, and what makes you unique..."
              class="flex w-full rounded-lg border border-input bg-background px-3 py-2.5 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 resize-y transition-colors"
            ></textarea>
          </div>
        </CardContent>
      </Card>

      <!-- Experience -->
      <Card class="animate-fade-in-up delay-100">
        <CardHeader class="pb-4">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-lg bg-blue-500/10 flex items-center justify-center">
              <Briefcase size={16} class="text-blue-500" />
            </div>
            <div>
              <CardTitle class="text-base">Experience</CardTitle>
              <CardDescription class="text-xs">Your work history and professional experience</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div class="space-y-2">
            <label for="experience" class="text-sm font-medium text-foreground">Work Experience</label>
            <textarea
              id="experience"
              bind:value={experience}
              rows="5"
              placeholder="List your work experience, including company names, roles, dates, and key achievements..."
              class="flex w-full rounded-lg border border-input bg-background px-3 py-2.5 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 resize-y transition-colors"
            ></textarea>
          </div>
        </CardContent>
      </Card>

      <!-- Education -->
      <Card class="animate-fade-in-up delay-200">
        <CardHeader class="pb-4">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-lg bg-green-500/10 flex items-center justify-center">
              <GraduationCap size={16} class="text-green-500" />
            </div>
            <div>
              <CardTitle class="text-base">Education</CardTitle>
              <CardDescription class="text-xs">Your academic background and certifications</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div class="space-y-2">
            <label for="education" class="text-sm font-medium text-foreground">Education</label>
            <textarea
              id="education"
              bind:value={education}
              rows="4"
              placeholder="Your degrees, institutions, graduation dates, and relevant coursework..."
              class="flex w-full rounded-lg border border-input bg-background px-3 py-2.5 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 resize-y transition-colors"
            ></textarea>
          </div>
        </CardContent>
      </Card>

      <!-- Skills & Projects -->
      <Card class="animate-fade-in-up delay-300">
        <CardHeader class="pb-4">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-lg bg-amber-500/10 flex items-center justify-center">
              <Code size={16} class="text-amber-500" />
            </div>
            <div>
              <CardTitle class="text-base">Skills & Projects</CardTitle>
              <CardDescription class="text-xs">Your technical skills and notable projects</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <label for="skills" class="text-sm font-medium text-foreground">Skills</label>
            <Input type="text" id="skills" bind:value={skills} placeholder="JavaScript, Python, React, Node.js, AWS..." />
            <p class="text-xs text-muted-foreground">Comma-separated list of your technical and soft skills</p>
          </div>

          <div class="space-y-2">
            <label for="projects" class="text-sm font-medium text-foreground">Notable Projects</label>
            <textarea
              id="projects"
              bind:value={projects}
              rows="4"
              placeholder="Describe your key projects, technologies used, and impact..."
              class="flex w-full rounded-lg border border-input bg-background px-3 py-2.5 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 resize-y transition-colors"
            ></textarea>
          </div>
        </CardContent>
      </Card>

      <!-- Save -->
      <div class="flex flex-col sm:flex-row items-center gap-4 pt-2 animate-fade-in-up delay-400">
        {#if message}
          <div class="flex items-center gap-2 text-sm {messageType === 'error' ? 'text-destructive' : 'text-green-600 dark:text-green-400'}">
            {#if messageType === 'error'}
              <AlertCircle size={16} />
            {:else}
              <CheckCircle size={16} />
            {/if}
            {message}
          </div>
        {/if}
        <Button type="submit" disabled={saving} class="sm:ml-auto">
          {#if saving}
            <Loader2 size={16} class="mr-2 animate-spin" />
            Saving...
          {:else}
            <Save size={16} class="mr-2" />
            Save Profile
          {/if}
        </Button>
      </div>
    </form>
  {/if}
</div>
