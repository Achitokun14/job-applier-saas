<script>
  import { auth } from '$lib/stores/auth';
  import { api } from '$lib/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Input } from '$lib/components/ui/input';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardContent } from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { MapPin, Building2, ExternalLink, CheckCircle, Search, Filter, Briefcase, DollarSign, Loader2, Square, CheckSquare, X } from 'lucide-svelte';
  import { toast } from 'svelte-sonner';

  let query = $state('');
  let location = $state('');
  let jobType = $state('all');
  let remote = $state('all');
  let jobs = $state([]);
  let loading = $state(false);
  let error = $state('');
  let applied = $state(new Set());
  let applyingTo = $state(null);
  let searched = $state(false);
  let selectedJobs = $state(new Set());
  let bulkApplying = $state(false);

  onMount(() => {
    if (!$auth.isAuthenticated) {
      goto('/login');
    }
  });

  async function search() {
    if (!query.trim()) return;
    loading = true;
    error = '';
    searched = true;
    selectedJobs = new Set();
    try {
      const token = $auth.token;
      const data = await api.searchJobs({
        query: query.trim(),
        location: location.trim() || undefined,
        remote: remote,
        jobType: jobType
      }, token);
      jobs = data.jobs || [];
    } catch (e) {
      error = 'Search failed. Please try again.';
    } finally {
      loading = false;
    }
  }

  async function applyToJob(jobId) {
    applyingTo = jobId;
    try {
      const token = $auth.token;
      await api.applyJob(jobId, token);
      applied = new Set([...applied, jobId]);
      selectedJobs.delete(jobId);
      selectedJobs = new Set(selectedJobs);
      toast.success('Application submitted successfully!');
    } catch (e) {
      toast.error('Application failed. Please try again.');
    } finally {
      applyingTo = null;
    }
  }

  function toggleSelect(jobId) {
    if (selectedJobs.has(jobId)) {
      selectedJobs.delete(jobId);
    } else {
      selectedJobs.add(jobId);
    }
    selectedJobs = new Set(selectedJobs);
  }

  function selectAll() {
    const selectableJobs = jobs.filter(j => !applied.has(j.id));
    selectedJobs = new Set(selectableJobs.map(j => j.id));
  }

  function deselectAll() {
    selectedJobs = new Set();
  }

  async function bulkApply() {
    if (selectedJobs.size === 0) return;
    bulkApplying = true;
    try {
      const token = $auth.token;
      const data = await api.bulkApply(Array.from(selectedJobs), token);
      const newApplied = new Set(applied);
      for (const id of selectedJobs) {
        newApplied.add(id);
      }
      applied = newApplied;
      toast.success(`Applied to ${data.applied} job${data.applied !== 1 ? 's' : ''}${data.skipped ? `, ${data.skipped} skipped` : ''}`);
      if (data.errors && data.errors.length > 0) {
        toast.error(`${data.errors.length} error(s) during bulk apply`);
      }
      selectedJobs = new Set();
    } catch (e) {
      toast.error('Bulk apply failed. Please try again.');
    } finally {
      bulkApplying = false;
    }
  }
</script>

<svelte:head>
  <title>Find Jobs - JobApplier</title>
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 {selectedJobs.size > 0 ? 'pb-24' : ''}">
  <!-- Header -->
  <div class="mb-8 animate-fade-in">
    <h1 class="text-2xl sm:text-3xl font-bold text-foreground">Find Jobs</h1>
    <p class="text-muted-foreground mt-1">Search thousands of job opportunities across companies</p>
  </div>

  <!-- Search bar -->
  <Card class="mb-6 animate-fade-in-up">
    <CardContent class="p-4 sm:p-6">
      <form onsubmit={(e) => { e.preventDefault(); search(); }} class="space-y-4">
        <div class="flex flex-col sm:flex-row gap-3">
          <div class="relative flex-1">
            <Search size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
            <Input
              type="text"
              bind:value={query}
              placeholder="Job title, keywords, or company..."
              class="pl-10"
            />
          </div>
          <div class="relative flex-1 sm:max-w-[240px]">
            <MapPin size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
            <Input
              type="text"
              bind:value={location}
              placeholder="Location..."
              class="pl-10"
            />
          </div>
          <Button type="submit" disabled={loading} class="sm:w-auto shrink-0">
            {#if loading}
              <Loader2 size={16} class="mr-2 animate-spin" />
              Searching...
            {:else}
              <Search size={16} class="mr-2" />
              Search
            {/if}
          </Button>
        </div>

        <!-- Filters row -->
        <div class="flex flex-wrap gap-3">
          <div class="flex items-center gap-2">
            <Filter size={14} class="text-muted-foreground" />
            <span class="text-sm text-muted-foreground">Filters:</span>
          </div>
          <select
            bind:value={remote}
            class="h-9 rounded-lg border border-input bg-background px-3 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
          >
            <option value="all">All work types</option>
            <option value="remote">Remote</option>
            <option value="hybrid">Hybrid</option>
            <option value="onsite">On-site</option>
          </select>
          <select
            bind:value={jobType}
            class="h-9 rounded-lg border border-input bg-background px-3 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
          >
            <option value="all">All job types</option>
            <option value="full_time">Full-time</option>
            <option value="part_time">Part-time</option>
            <option value="contract">Contract</option>
            <option value="internship">Internship</option>
          </select>
        </div>
      </form>
    </CardContent>
  </Card>

  {#if error}
    <div class="bg-destructive/10 text-destructive text-sm p-3.5 rounded-lg mb-4 border border-destructive/20 animate-fade-in">
      {error}
    </div>
  {/if}

  <!-- Results -->
  <div class="space-y-3">
    {#if loading}
      <!-- Skeleton loading -->
      {#each [1,2,3,4,5] as _}
        <Card>
          <CardContent class="p-5">
            <div class="flex justify-between items-start gap-4">
              <div class="flex-1 space-y-3">
                <div class="skeleton h-5 w-48"></div>
                <div class="skeleton h-4 w-32"></div>
                <div class="skeleton h-4 w-40"></div>
              </div>
              <div class="skeleton h-9 w-20 rounded-lg"></div>
            </div>
          </CardContent>
        </Card>
      {/each}
    {:else if !searched}
      <!-- Initial state -->
      <Card class="border-dashed">
        <CardContent class="p-16 text-center">
          <div class="w-16 h-16 rounded-2xl bg-muted flex items-center justify-center mx-auto mb-4">
            <Search size={28} class="text-muted-foreground" />
          </div>
          <h3 class="text-lg font-semibold text-foreground mb-2">Search for jobs</h3>
          <p class="text-sm text-muted-foreground max-w-sm mx-auto">
            Enter keywords like job title, skills, or company name to find matching opportunities.
          </p>
        </CardContent>
      </Card>
    {:else if jobs.length === 0}
      <!-- No results -->
      <Card class="border-dashed">
        <CardContent class="p-16 text-center">
          <div class="w-16 h-16 rounded-2xl bg-muted flex items-center justify-center mx-auto mb-4">
            <Briefcase size={28} class="text-muted-foreground" />
          </div>
          <h3 class="text-lg font-semibold text-foreground mb-2">No jobs found</h3>
          <p class="text-sm text-muted-foreground max-w-sm mx-auto">
            Try adjusting your search terms or filters to find more opportunities.
          </p>
        </CardContent>
      </Card>
    {:else}
      <div class="flex items-center justify-between mb-2">
        <p class="text-sm text-muted-foreground">{jobs.length} job{jobs.length !== 1 ? 's' : ''} found</p>
        <div class="flex items-center gap-2">
          {#if selectedJobs.size === jobs.filter(j => !applied.has(j.id)).length && jobs.filter(j => !applied.has(j.id)).length > 0}
            <Button variant="outline" size="sm" onclick={deselectAll}>
              <X size={14} class="mr-1.5" />
              Deselect All
            </Button>
          {:else}
            <Button variant="outline" size="sm" onclick={selectAll}>
              <CheckSquare size={14} class="mr-1.5" />
              Select All
            </Button>
          {/if}
        </div>
      </div>
      {#each jobs as job, i}
        <Card class="group hover:shadow-md hover:border-primary/20 transition-all duration-300 animate-fade-in {selectedJobs.has(job.id) ? 'border-primary/40 bg-primary/5' : ''}" style="animation-delay: {i * 50}ms">
          <CardContent class="p-5">
            <div class="flex flex-col sm:flex-row sm:items-start gap-4">
              <div class="flex items-start gap-4 flex-1 min-w-0">
                <!-- Checkbox -->
                {#if !applied.has(job.id)}
                  <button
                    onclick={() => toggleSelect(job.id)}
                    class="mt-0.5 shrink-0 text-muted-foreground hover:text-primary transition-colors cursor-pointer"
                    aria-label={selectedJobs.has(job.id) ? 'Deselect job' : 'Select job'}
                  >
                    {#if selectedJobs.has(job.id)}
                      <CheckSquare size={20} class="text-primary" />
                    {:else}
                      <Square size={20} />
                    {/if}
                  </button>
                {:else}
                  <div class="w-5 shrink-0"></div>
                {/if}
                <!-- Company avatar -->
                <div class="w-11 h-11 rounded-xl bg-primary/10 flex items-center justify-center text-sm font-bold text-primary shrink-0">
                  {(job.company || 'C')[0].toUpperCase()}
                </div>
                <div class="flex-1 min-w-0">
                  <h3 class="text-base font-semibold text-foreground m-0 group-hover:text-primary transition-colors">{job.title}</h3>
                  <div class="flex items-center gap-1.5 mt-1">
                    <Building2 size={14} class="text-primary shrink-0" />
                    <span class="text-sm text-primary font-medium">{job.company}</span>
                  </div>
                  <div class="flex flex-wrap items-center gap-x-4 gap-y-1 mt-2">
                    {#if job.location}
                      <div class="flex items-center gap-1.5">
                        <MapPin size={13} class="text-muted-foreground shrink-0" />
                        <span class="text-sm text-muted-foreground">{job.location}</span>
                      </div>
                    {/if}
                    {#if job.salary}
                      <div class="flex items-center gap-1.5">
                        <DollarSign size={13} class="text-muted-foreground shrink-0" />
                        <span class="text-sm text-muted-foreground">{job.salary}</span>
                      </div>
                    {/if}
                  </div>
                  {#if job.tags && job.tags.length > 0}
                    <div class="flex flex-wrap gap-1.5 mt-2.5">
                      {#each job.tags.slice(0, 4) as tag}
                        <Badge variant="secondary" class="text-xs">{tag}</Badge>
                      {/each}
                    </div>
                  {/if}
                </div>
              </div>
              <div class="flex items-center gap-3 sm:shrink-0 sm:self-center">
                {#if applied.has(job.id)}
                  <span class="flex items-center gap-1.5 text-green-600 dark:text-green-400 font-semibold text-sm">
                    <CheckCircle size={16} />
                    Applied
                  </span>
                {:else}
                  <Button onclick={() => applyToJob(job.id)} size="sm" disabled={applyingTo === job.id}>
                    {#if applyingTo === job.id}
                      <Loader2 size={14} class="mr-1.5 animate-spin" />
                      Applying...
                    {:else}
                      Apply
                    {/if}
                  </Button>
                {/if}
                {#if job.url}
                  <a
                    href={job.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    class="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors no-underline"
                  >
                    <ExternalLink size={14} />
                  </a>
                {/if}
              </div>
            </div>
          </CardContent>
        </Card>
      {/each}
    {/if}
  </div>
</div>

<!-- Sticky bulk apply bar -->
{#if selectedJobs.size > 0}
  <div class="fixed bottom-0 left-0 right-0 z-40 border-t border-border bg-background/95 backdrop-blur-lg shadow-lg">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <Badge variant="secondary" class="text-sm px-3 py-1">
          {selectedJobs.size} selected
        </Badge>
        <button onclick={deselectAll} class="text-sm text-muted-foreground hover:text-foreground transition-colors cursor-pointer">
          Clear selection
        </button>
      </div>
      <Button onclick={bulkApply} disabled={bulkApplying}>
        {#if bulkApplying}
          <Loader2 size={16} class="mr-2 animate-spin" />
          Applying...
        {:else}
          <CheckCircle size={16} class="mr-2" />
          Apply to {selectedJobs.size} selected
        {/if}
      </Button>
    </div>
  </div>
{/if}
