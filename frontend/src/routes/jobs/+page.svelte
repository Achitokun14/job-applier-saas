<script>
  import { auth } from '$lib/stores/auth';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Input } from '$lib/components/ui/input';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardContent } from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { MapPin, Building2, ExternalLink, CheckCircle } from 'lucide-svelte';

  let query = $state('');
  let jobs = $state([]);
  let loading = $state(false);
  let error = $state('');
  let applied = $state(new Set());

  onMount(() => {
    if (!$auth.isAuthenticated) {
      goto('/login');
    }
  });

  async function search() {
    if (!query.trim()) return;
    loading = true;
    error = '';
    try {
      const token = $auth.token;
      const res = await fetch('/api/v1/jobs?q=' + encodeURIComponent(query), {
        headers: { 'Authorization': 'Bearer ' + token }
      });
      if (res.ok) {
        jobs = await res.json();
      }
    } catch (e) {
      error = 'Search failed';
    } finally {
      loading = false;
    }
  }

  async function apply(jobId) {
    try {
      const token = $auth.token;
      const res = await fetch('/api/v1/jobs/' + jobId + '/apply', {
        method: 'POST',
        headers: { 'Authorization': 'Bearer ' + token }
      });
      if (res.ok) {
        applied.add(jobId);
        applied = applied;
      }
    } catch (e) {
      error = 'Application failed';
    }
  }
</script>

<div class="max-w-[800px] mx-auto">
  <div class="text-center mb-8">
    <h1 class="text-3xl font-bold text-foreground m-0">Find Jobs</h1>
    <p class="text-muted-foreground mt-2 mb-0">Search thousands of job opportunities</p>
  </div>

  <div class="flex gap-2 mb-6">
    <Input
      type="text"
      bind:value={query}
      placeholder="Search (e.g., Software Engineer)"
      onkeydown={e => e.key === 'Enter' && search()}
      class="flex-1"
    />
    <Button onclick={search} disabled={loading}>
      {loading ? 'Searching...' : 'Search'}
    </Button>
  </div>

  {#if error}
    <div class="bg-destructive/10 text-destructive p-3 rounded-md mb-4 text-sm">{error}</div>
  {/if}

  <div class="flex flex-col gap-4">
    {#if jobs.length === 0 && !loading}
      <p class="text-center text-muted-foreground py-12">Enter keywords to find jobs</p>
    {:else}
      {#each jobs as job}
        <Card>
          <CardContent class="p-5 flex justify-between items-start gap-4">
            <div class="flex-1 min-w-0">
              <h3 class="text-base font-semibold text-foreground m-0">{job.title}</h3>
              <div class="flex items-center gap-1.5 mt-1">
                <Building2 size={14} class="text-primary shrink-0" />
                <span class="text-sm text-primary font-medium">{job.company}</span>
              </div>
              {#if job.location}
                <div class="flex items-center gap-1.5 mt-1">
                  <MapPin size={14} class="text-muted-foreground shrink-0" />
                  <span class="text-sm text-muted-foreground">{job.location}</span>
                </div>
              {/if}
              {#if job.salary}
                <Badge variant="secondary" class="mt-2">{job.salary}</Badge>
              {/if}
            </div>
            <div class="flex items-center gap-3 shrink-0">
              {#if applied.has(job.id)}
                <span class="flex items-center gap-1.5 text-green-600 dark:text-green-400 font-semibold text-sm">
                  <CheckCircle size={16} />
                  Applied
                </span>
              {:else}
                <Button onclick={() => apply(job.id)} size="sm">Apply</Button>
              {/if}
              {#if job.url}
                <a
                  href={job.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  class="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors no-underline"
                >
                  <ExternalLink size={14} />
                  View
                </a>
              {/if}
            </div>
          </CardContent>
        </Card>
      {/each}
    {/if}
  </div>
</div>
