<script>
  import { auth } from '$lib/stores/auth';
  import { api } from '$lib/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';

  let settings = $state({});
  let llmProvider = $state('openai');
  let llmModel = $state('gpt-4o-mini');
  let llmApiKey = $state('');
  let jobSearchRemote = $state(true);
  let jobSearchHybrid = $state(true);
  let jobSearchOnsite = $state(false);
  let experienceLevel = $state('mid_senior');
  let jobTypes = $state('full_time');
  let positions = $state('');
  let locations = $state('');
  let distance = $state(50);
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
      settings = await api.getSettings(token);
      if (settings.llmProvider) llmProvider = settings.llmProvider;
      if (settings.llmModel) llmModel = settings.llmModel;
      if (settings.jobSearchRemote !== undefined) jobSearchRemote = settings.jobSearchRemote;
      if (settings.jobSearchHybrid !== undefined) jobSearchHybrid = settings.jobSearchHybrid;
      if (settings.experienceLevel) experienceLevel = settings.experienceLevel;
      if (settings.jobTypes) jobTypes = settings.jobTypes;
      if (settings.locations) locations = settings.locations.join(', ');
      if (settings.distance) distance = settings.distance;
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  });

  async function saveSettings() {
    saving = true;
    message = '';
    try {
      const token = $auth.token;
      await api.updateSettings({
        llmProvider,
        llmModel,
        llmApiKey,
        jobSearchRemote,
        jobSearchHybrid,
        jobSearchOnsite,
        experienceLevel,
        jobTypes,
        locations: locations.split(',').map(l => l.trim()).filter(l => l),
        distance
      }, token);
      message = 'Settings saved!';
    } catch (e) {
      message = 'Failed to save settings';
    } finally {
      saving = false;
    }
  }
</script>

<div class="max-w-2xl mx-auto">
  <div class="mb-6">
    <h1 class="text-3xl font-bold text-foreground">Settings</h1>
    <p class="text-muted-foreground mt-1">Configure your AI provider and job search preferences</p>
  </div>

  {#if loading}
    <div class="text-center py-12 text-muted-foreground">Loading...</div>
  {:else}
    <form onsubmit={(e) => { e.preventDefault(); saveSettings(); }} class="space-y-6">
      <!-- AI Configuration -->
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">AI Configuration</CardTitle>
          <CardDescription>Choose your AI provider and model for generating applications</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <label for="llmProvider" class="text-sm font-medium text-foreground">LLM Provider</label>
            <select
              id="llmProvider"
              bind:value={llmProvider}
              class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            >
              <option value="openai">OpenAI</option>
              <option value="anthropic">Anthropic</option>
              <option value="google">Google Gemini</option>
            </select>
          </div>

          <div class="space-y-2">
            <label for="llmModel" class="text-sm font-medium text-foreground">Model</label>
            <select
              id="llmModel"
              bind:value={llmModel}
              class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            >
              <option value="gpt-4o-mini">GPT-4o Mini</option>
              <option value="gpt-4o">GPT-4o</option>
              <option value="claude-3-haiku">Claude 3 Haiku</option>
              <option value="gemini-pro">Gemini Pro</option>
            </select>
          </div>

          <div class="space-y-2">
            <label for="llmApiKey" class="text-sm font-medium text-foreground">API Key</label>
            <Input type="password" id="llmApiKey" bind:value={llmApiKey} placeholder="Enter API key" />
          </div>
        </CardContent>
      </Card>

      <!-- Job Search Preferences -->
      <Card>
        <CardHeader>
          <CardTitle class="text-lg">Job Search Preferences</CardTitle>
          <CardDescription>Set your preferred work type, experience level, and locations</CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <label class="text-sm font-medium text-foreground">Work Type</label>
            <div class="flex flex-wrap gap-4">
              <label class="flex items-center gap-2 text-sm text-foreground cursor-pointer">
                <input
                  type="checkbox"
                  bind:checked={jobSearchRemote}
                  class="h-4 w-4 rounded border-input text-primary focus:ring-ring"
                />
                Remote
              </label>
              <label class="flex items-center gap-2 text-sm text-foreground cursor-pointer">
                <input
                  type="checkbox"
                  bind:checked={jobSearchHybrid}
                  class="h-4 w-4 rounded border-input text-primary focus:ring-ring"
                />
                Hybrid
              </label>
              <label class="flex items-center gap-2 text-sm text-foreground cursor-pointer">
                <input
                  type="checkbox"
                  bind:checked={jobSearchOnsite}
                  class="h-4 w-4 rounded border-input text-primary focus:ring-ring"
                />
                On-site
              </label>
            </div>
          </div>

          <div class="space-y-2">
            <label for="experienceLevel" class="text-sm font-medium text-foreground">Experience Level</label>
            <select
              id="experienceLevel"
              bind:value={experienceLevel}
              class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            >
              <option value="entry">Entry</option>
              <option value="mid">Mid-Level</option>
              <option value="mid_senior">Mid-Senior</option>
              <option value="senior">Senior</option>
              <option value="lead">Lead</option>
            </select>
          </div>

          <div class="space-y-2">
            <label for="locations" class="text-sm font-medium text-foreground">Preferred Locations</label>
            <Input type="text" id="locations" bind:value={locations} placeholder="e.g., Casablanca, Rabat, Remote" />
          </div>
        </CardContent>
      </Card>

      <!-- Actions -->
      <div class="flex items-center gap-4 pt-2">
        {#if message}
          <span class="text-sm {message.includes('Failed') ? 'text-destructive' : 'text-green-600 dark:text-green-400'}">{message}</span>
        {/if}
        <Button type="submit" disabled={saving} class="ml-auto">
          {saving ? 'Saving...' : 'Save Settings'}
        </Button>
      </div>
    </form>
  {/if}
</div>
