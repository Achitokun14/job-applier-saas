<script>
  import { auth } from '$lib/stores/auth';
  import { api } from '$lib/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';

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
      message = 'Profile saved!';
    } catch (e) {
      message = 'Failed to save profile';
    } finally {
      saving = false;
    }
  }
</script>

<div class="max-w-2xl mx-auto">
  <div class="mb-6">
    <h1 class="text-3xl font-bold text-foreground">My Profile</h1>
    <p class="text-muted-foreground mt-1">Manage your personal information and professional details</p>
  </div>

  {#if loading}
    <div class="text-center py-12 text-muted-foreground">Loading...</div>
  {:else}
    <form onsubmit={(e) => { e.preventDefault(); saveProfile(); }} class="space-y-6">
      <!-- Personal Info -->
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Personal Information</CardTitle>
          <CardDescription>Your name and personal statement</CardDescription>
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
              placeholder="Brief about yourself..."
              class="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 resize-y"
            ></textarea>
          </div>
        </CardContent>
      </Card>

      <!-- Experience -->
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Experience</CardTitle>
          <CardDescription>Your work history and professional experience</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-2">
            <label for="experience" class="text-sm font-medium text-foreground">Work Experience</label>
            <textarea
              id="experience"
              bind:value={experience}
              rows="4"
              placeholder="Your work experience..."
              class="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 resize-y"
            ></textarea>
          </div>
        </CardContent>
      </Card>

      <!-- Education -->
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Education</CardTitle>
          <CardDescription>Your academic background</CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-2">
            <label for="education" class="text-sm font-medium text-foreground">Education</label>
            <textarea
              id="education"
              bind:value={education}
              rows="3"
              placeholder="Your education..."
              class="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 resize-y"
            ></textarea>
          </div>
        </CardContent>
      </Card>

      <!-- Skills & Projects -->
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Skills & Projects</CardTitle>
          <CardDescription>Your technical skills and notable projects</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <label for="skills" class="text-sm font-medium text-foreground">Skills</label>
            <Input type="text" id="skills" bind:value={skills} placeholder="JavaScript, Python, React..." />
          </div>

          <div class="space-y-2">
            <label for="projects" class="text-sm font-medium text-foreground">Projects</label>
            <textarea
              id="projects"
              bind:value={projects}
              rows="3"
              placeholder="Notable projects..."
              class="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 resize-y"
            ></textarea>
          </div>
        </CardContent>
      </Card>

      <!-- Actions -->
      <div class="flex items-center gap-4 pt-2">
        {#if message}
          <span class="text-sm {message.includes('Failed') ? 'text-destructive' : 'text-green-600 dark:text-green-400'}">{message}</span>
        {/if}
        <Button type="submit" disabled={saving} class="ml-auto">
          {saving ? 'Saving...' : 'Save Profile'}
        </Button>
      </div>
    </form>
  {/if}
</div>
