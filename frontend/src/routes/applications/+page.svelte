<script>
  import { auth } from '$lib/stores/auth';
  import { api } from '$lib/api';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Card, CardContent } from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Trash2, Calendar, FileText, Search, AlertCircle, Loader2, Filter } from 'lucide-svelte';
  import { toast } from 'svelte-sonner';

  let applications = $state([]);
  let loading = $state(true);
  let error = $state('');
  let deletingId = $state(null);
  let statusFilter = $state('all');

  onMount(async () => {
    if (!$auth.isAuthenticated) {
      goto('/login');
      return;
    }

    await loadApplications();
  });

  async function loadApplications() {
    loading = true;
    error = '';
    try {
      const token = $auth.token;
      const data = await api.getApplications(token);
      applications = data.applications || data || [];
    } catch (e) {
      error = e.message || 'Failed to load applications';
    } finally {
      loading = false;
    }
  }

  let filteredApplications = $derived(
    statusFilter === 'all'
      ? applications
      : applications.filter(a => a.status === statusFilter)
  );

  async function deleteApplication(id) {
    if (!window.confirm('Are you sure you want to delete this application? This action cannot be undone.')) {
      return;
    }
    deletingId = id;
    try {
      const token = $auth.token;
      await api.deleteApplication(id, token);
      applications = applications.filter(a => a.id !== id);
      toast.success('Application deleted successfully');
    } catch (e) {
      toast.error(e.message || 'Failed to delete application');
    } finally {
      deletingId = null;
    }
  }

  function formatDate(dateStr) {
    if (!dateStr) return '-';
    try {
      return new Date(dateStr).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      });
    } catch {
      return '-';
    }
  }
</script>

<svelte:head>
  <title>My Applications - JobApplier</title>
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
  <!-- Header -->
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8 animate-fade-in">
    <div>
      <h1 class="text-2xl sm:text-3xl font-bold text-foreground">My Applications</h1>
      <p class="text-muted-foreground mt-1">Track and manage all your job applications</p>
    </div>
    {#if !loading && applications.length > 0}
      <div class="flex items-center gap-2">
        <Badge variant="secondary" class="text-xs px-3 py-1">
          {filteredApplications.length} of {applications.length} total
        </Badge>
      </div>
    {/if}
  </div>

  <!-- Status Filter -->
  {#if !loading && applications.length > 0}
    <div class="flex items-center gap-3 mb-6 animate-fade-in">
      <Filter size={14} class="text-muted-foreground" />
      <span class="text-sm text-muted-foreground">Status:</span>
      <select
        bind:value={statusFilter}
        class="h-9 rounded-lg border border-input bg-background px-3 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
      >
        <option value="all">All statuses</option>
        <option value="applied">Applied</option>
        <option value="interview">Interview</option>
        <option value="offer">Offer</option>
        <option value="rejected">Rejected</option>
      </select>
    </div>
  {/if}

  {#if error}
    <div class="flex items-start gap-3 bg-destructive/10 text-destructive text-sm p-3.5 rounded-lg mb-6 border border-destructive/20 animate-fade-in">
      <AlertCircle size={18} class="shrink-0 mt-0.5" />
      <span>{error}</span>
    </div>
  {/if}

  {#if loading}
    <!-- Skeleton -->
    <Card class="animate-fade-in">
      <div class="p-6 space-y-4">
        {#each [1,2,3,4,5] as _}
          <div class="flex items-center gap-4">
            <div class="skeleton h-10 w-10 rounded-lg"></div>
            <div class="flex-1 space-y-2">
              <div class="skeleton h-4 w-48"></div>
              <div class="skeleton h-3 w-32"></div>
            </div>
            <div class="skeleton h-6 w-16 rounded-full"></div>
            <div class="skeleton h-8 w-8 rounded-lg"></div>
          </div>
        {/each}
      </div>
    </Card>
  {:else if applications.length === 0}
    <!-- Empty state -->
    <Card class="border-dashed animate-fade-in-up">
      <CardContent class="p-16 text-center">
        <div class="w-16 h-16 rounded-2xl bg-muted flex items-center justify-center mx-auto mb-4">
          <FileText size={28} class="text-muted-foreground" />
        </div>
        <h3 class="text-lg font-semibold text-foreground mb-2">No applications yet</h3>
        <p class="text-sm text-muted-foreground mb-6 max-w-sm mx-auto">
          Start by searching for jobs and submitting your first application.
        </p>
        <a href="/jobs" class="no-underline">
          <Button>
            <Search size={16} class="mr-2" />
            Find Jobs
          </Button>
        </a>
      </CardContent>
    </Card>
  {:else if filteredApplications.length === 0}
    <!-- Empty filtered state -->
    <Card class="border-dashed animate-fade-in-up">
      <CardContent class="p-12 text-center">
        <div class="w-14 h-14 rounded-2xl bg-muted flex items-center justify-center mx-auto mb-4">
          <Filter size={24} class="text-muted-foreground" />
        </div>
        <h3 class="text-lg font-semibold text-foreground mb-2">No matching applications</h3>
        <p class="text-sm text-muted-foreground mb-4 max-w-sm mx-auto">
          No applications found with status "{statusFilter}". Try a different filter.
        </p>
        <Button variant="outline" onclick={() => statusFilter = 'all'}>
          Show All
        </Button>
      </CardContent>
    </Card>
  {:else}
    <!-- Table view -->
    <Card class="animate-fade-in-up overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-border bg-muted/30">
              <th class="text-left text-xs font-medium text-muted-foreground uppercase tracking-wider p-4">Company</th>
              <th class="text-left text-xs font-medium text-muted-foreground uppercase tracking-wider p-4 hidden sm:table-cell">Position</th>
              <th class="text-left text-xs font-medium text-muted-foreground uppercase tracking-wider p-4">Status</th>
              <th class="text-left text-xs font-medium text-muted-foreground uppercase tracking-wider p-4 hidden md:table-cell">Applied</th>
              <th class="text-right text-xs font-medium text-muted-foreground uppercase tracking-wider p-4">Actions</th>
            </tr>
          </thead>
          <tbody>
            {#each filteredApplications as app}
              <tr class="border-b border-border/50 last:border-0 hover:bg-muted/20 transition-colors group">
                <td class="p-4">
                  <div class="flex items-center gap-3">
                    <div class="w-9 h-9 rounded-lg bg-primary/10 flex items-center justify-center text-xs font-bold text-primary shrink-0">
                      {(app.job?.company || 'U')[0].toUpperCase()}
                    </div>
                    <div class="min-w-0">
                      <div class="font-medium text-foreground text-sm truncate">{app.job?.company || 'Unknown Company'}</div>
                      <div class="text-xs text-muted-foreground sm:hidden truncate">{app.job?.title || 'Unknown'}</div>
                    </div>
                  </div>
                </td>
                <td class="p-4 hidden sm:table-cell">
                  <span class="text-sm text-foreground">{app.job?.title || 'Unknown Position'}</span>
                </td>
                <td class="p-4">
                  {#if app.status === 'applied'}
                    <Badge class="bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300 border-transparent text-xs">Applied</Badge>
                  {:else if app.status === 'interview'}
                    <Badge class="bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300 border-transparent text-xs">Interview</Badge>
                  {:else if app.status === 'offer'}
                    <Badge class="bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-300 border-transparent text-xs">Offer</Badge>
                  {:else if app.status === 'rejected'}
                    <Badge variant="destructive" class="text-xs">Rejected</Badge>
                  {:else}
                    <Badge variant="secondary" class="text-xs">{app.status}</Badge>
                  {/if}
                </td>
                <td class="p-4 hidden md:table-cell">
                  <div class="flex items-center gap-1.5">
                    <Calendar size={12} class="text-muted-foreground" />
                    <span class="text-sm text-muted-foreground">
                      {formatDate(app.applied_at || app.appliedAt)}
                    </span>
                  </div>
                </td>
                <td class="p-4 text-right">
                  <Button
                    variant="ghost"
                    size="icon"
                    onclick={() => deleteApplication(app.id)}
                    disabled={deletingId === app.id}
                    class="h-8 w-8 text-muted-foreground hover:text-destructive hover:bg-destructive/10"
                  >
                    {#if deletingId === app.id}
                      <Loader2 size={14} class="animate-spin" />
                    {:else}
                      <Trash2 size={14} />
                    {/if}
                  </Button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </Card>
  {/if}
</div>
