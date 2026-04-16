<script>
  import { auth } from '$lib/stores/auth';
  import { api } from '$lib/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';
  import { Bot, MapPin, CheckCircle, AlertCircle, Loader2, Save, Eye, EyeOff, Key } from 'lucide-svelte';

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
  let messageType = $state('success');
  let showApiKey = $state(false);

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
      message = 'Settings saved successfully!';
      messageType = 'success';
      setTimeout(() => { message = ''; }, 3000);
    } catch (e) {
      message = 'Failed to save settings. Please try again.';
      messageType = 'error';
    } finally {
      saving = false;
    }
  }

  let availableModels = $derived(() => {
    if (llmProvider === 'openai') return [
      { value: 'gpt-4o-mini', label: 'GPT-4o Mini' },
      { value: 'gpt-4o', label: 'GPT-4o' },
      { value: 'gpt-4-turbo', label: 'GPT-4 Turbo' },
    ];
    if (llmProvider === 'anthropic') return [
      { value: 'claude-3-haiku', label: 'Claude 3 Haiku' },
      { value: 'claude-3-sonnet', label: 'Claude 3.5 Sonnet' },
      { value: 'claude-3-opus', label: 'Claude 3 Opus' },
    ];
    return [
      { value: 'gemini-pro', label: 'Gemini Pro' },
      { value: 'gemini-ultra', label: 'Gemini Ultra' },
    ];
  });
</script>

<svelte:head>
  <title>Settings - JobApplier</title>
</svelte:head>

<div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
  <!-- Header -->
  <div class="mb-8 animate-fade-in">
    <h1 class="text-2xl sm:text-3xl font-bold text-foreground">Settings</h1>
    <p class="text-muted-foreground mt-1">Configure your AI provider and job search preferences</p>
  </div>

  {#if loading}
    <div class="space-y-6 animate-fade-in">
      {#each [1,2] as _}
        <div class="skeleton h-64 rounded-xl"></div>
      {/each}
    </div>
  {:else}
    <form onsubmit={(e) => { e.preventDefault(); saveSettings(); }} class="space-y-6">
      <!-- AI Configuration -->
      <Card class="animate-fade-in-up">
        <CardHeader class="pb-4">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-lg bg-primary/10 flex items-center justify-center">
              <Bot size={16} class="text-primary" />
            </div>
            <div>
              <CardTitle class="text-base">AI Configuration</CardTitle>
              <CardDescription class="text-xs">Choose your AI provider and model for generating applications</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-2">
              <label for="llmProvider" class="text-sm font-medium text-foreground">Provider</label>
              <select
                id="llmProvider"
                bind:value={llmProvider}
                class="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 transition-colors"
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
                class="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 transition-colors"
              >
                {#each availableModels() as model}
                  <option value={model.value}>{model.label}</option>
                {/each}
              </select>
            </div>
          </div>

          <div class="space-y-2">
            <label for="llmApiKey" class="text-sm font-medium text-foreground flex items-center gap-2">
              <Key size={14} class="text-muted-foreground" />
              API Key
            </label>
            <div class="relative">
              <Input
                type={showApiKey ? 'text' : 'password'}
                id="llmApiKey"
                bind:value={llmApiKey}
                placeholder="sk-..."
                class="pr-10 font-mono text-xs"
              />
              <button
                type="button"
                onclick={() => showApiKey = !showApiKey}
                class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
              >
                {#if showApiKey}
                  <EyeOff size={16} />
                {:else}
                  <Eye size={16} />
                {/if}
              </button>
            </div>
            <p class="text-xs text-muted-foreground">Your API key is encrypted with AES-256-GCM before storage</p>
          </div>
        </CardContent>
      </Card>

      <!-- Job Search Preferences -->
      <Card class="animate-fade-in-up delay-100">
        <CardHeader class="pb-4">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-lg bg-blue-500/10 flex items-center justify-center">
              <MapPin size={16} class="text-blue-500" />
            </div>
            <div>
              <CardTitle class="text-base">Job Search Preferences</CardTitle>
              <CardDescription class="text-xs">Set your preferred work type, experience level, and locations</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent class="space-y-5">
          <!-- Work Type -->
          <div class="space-y-3">
            <label class="text-sm font-medium text-foreground">Work Type</label>
            <div class="flex flex-wrap gap-3">
              <label class="flex items-center gap-2.5 px-4 py-2.5 rounded-lg border border-input hover:border-primary/30 transition-colors cursor-pointer {jobSearchRemote ? 'bg-primary/5 border-primary/30' : ''}">
                <input
                  type="checkbox"
                  bind:checked={jobSearchRemote}
                  class="h-4 w-4 rounded border-input text-primary focus:ring-ring"
                />
                <span class="text-sm text-foreground">Remote</span>
              </label>
              <label class="flex items-center gap-2.5 px-4 py-2.5 rounded-lg border border-input hover:border-primary/30 transition-colors cursor-pointer {jobSearchHybrid ? 'bg-primary/5 border-primary/30' : ''}">
                <input
                  type="checkbox"
                  bind:checked={jobSearchHybrid}
                  class="h-4 w-4 rounded border-input text-primary focus:ring-ring"
                />
                <span class="text-sm text-foreground">Hybrid</span>
              </label>
              <label class="flex items-center gap-2.5 px-4 py-2.5 rounded-lg border border-input hover:border-primary/30 transition-colors cursor-pointer {jobSearchOnsite ? 'bg-primary/5 border-primary/30' : ''}">
                <input
                  type="checkbox"
                  bind:checked={jobSearchOnsite}
                  class="h-4 w-4 rounded border-input text-primary focus:ring-ring"
                />
                <span class="text-sm text-foreground">On-site</span>
              </label>
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="space-y-2">
              <label for="experienceLevel" class="text-sm font-medium text-foreground">Experience Level</label>
              <select
                id="experienceLevel"
                bind:value={experienceLevel}
                class="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 transition-colors"
              >
                <option value="entry">Entry Level</option>
                <option value="mid">Mid-Level</option>
                <option value="mid_senior">Mid-Senior</option>
                <option value="senior">Senior</option>
                <option value="lead">Lead / Principal</option>
              </select>
            </div>

            <div class="space-y-2">
              <label for="jobTypes" class="text-sm font-medium text-foreground">Job Type</label>
              <select
                id="jobTypes"
                bind:value={jobTypes}
                class="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 transition-colors"
              >
                <option value="full_time">Full-time</option>
                <option value="part_time">Part-time</option>
                <option value="contract">Contract</option>
                <option value="internship">Internship</option>
              </select>
            </div>
          </div>

          <div class="space-y-2">
            <label for="locations" class="text-sm font-medium text-foreground">Preferred Locations</label>
            <Input type="text" id="locations" bind:value={locations} placeholder="e.g., San Francisco, New York, Remote" />
            <p class="text-xs text-muted-foreground">Comma-separated list of preferred locations</p>
          </div>
        </CardContent>
      </Card>

      <!-- Save -->
      <div class="flex flex-col sm:flex-row items-center gap-4 pt-2 animate-fade-in-up delay-200">
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
            Save Settings
          {/if}
        </Button>
      </div>
    </form>
  {/if}
</div>
