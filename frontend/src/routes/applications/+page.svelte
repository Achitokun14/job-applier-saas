<script>
  import { auth } from '$lib/stores/auth';
  import { api } from '$lib/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Card, CardContent } from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Trash2, Calendar } from 'lucide-svelte';

  let applications = $state([]);
  let loading = $state(true);
  let error = $state('');

  onMount(async () => {
    if (!$auth.isAuthenticated) {
      goto('/login');
      return;
    }

    try {
      const token = $auth.token;
      applications = await api.getApplications(token);
    } catch (e) {
      error = e.message || 'Failed to load applications';
    } finally {
      loading = false;
    }
  });

  async function deleteApplication(id) {
    try {
      const token = $auth.token;
      await api.deleteApplication(id, token);
      applications = applications.filter(a => a.id !== id);
    } catch (e) {
      error = e.message || 'Failed to delete application';
    }
  }
</script>

<div class="max-w-[800px] mx-auto">
  <h1 class="text-3xl font-bold text-foreground mb-6">My Applications</h1>

  {#if loading}
    <div class="text-center py-12 text-muted-foreground">Loading...</div>
  {:else if error}
    <div class="bg-destructive/10 text-destructive p-3 rounded-md mb-4 text-sm">{error}</div>
  {:else if applications.length === 0}
    <div class="text-center py-12">
      <p class="text-muted-foreground mb-4">No applications yet.</p>
      <a href="/jobs" class="no-underline">
        <Button>Find Jobs</Button>
      </a>
    </div>
  {:else}
    <div class="flex flex-col gap-3">
      {#each applications as app}
        <Card>
          <CardContent class="p-5 flex justify-between items-center gap-4">
            <div class="flex-1 min-w-0">
              <h3 class="text-base font-semibold text-foreground m-0">{app.job?.title || 'Unknown Position'}</h3>
              <p class="text-sm text-muted-foreground m-0 mt-0.5">{app.job?.company || 'Unknown Company'}</p>
              {#if app.appliedAt}
                <div class="flex items-center gap-1.5 mt-1.5">
                  <Calendar size={12} class="text-muted-foreground" />
                  <span class="text-xs text-muted-foreground">Applied: {app.appliedAt}</span>
                </div>
              {/if}
            </div>
            <div class="flex items-center gap-3 shrink-0">
              {#if app.status === 'applied'}
                <Badge class="bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300 border-transparent">Applied</Badge>
              {:else if app.status === 'interview'}
                <Badge class="bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300 border-transparent">Interview</Badge>
              {:else if app.status === 'offer'}
                <Badge class="bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300 border-transparent">Offer</Badge>
              {:else if app.status === 'rejected'}
                <Badge variant="destructive">Rejected</Badge>
              {:else}
                <Badge variant="secondary">{app.status}</Badge>
              {/if}
              <Button variant="destructive" size="icon" onclick={() => deleteApplication(app.id)} class="h-8 w-8">
                <Trash2 size={14} />
              </Button>
            </div>
          </CardContent>
        </Card>
      {/each}
    </div>
  {/if}
</div>
