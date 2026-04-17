<script>
  import { auth } from '$lib/stores/auth';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Card, CardHeader, CardTitle, CardContent } from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Briefcase, Clock, MessageSquare, Award, Search, FileText, User, Settings, ArrowRight, Plus, TrendingUp } from 'lucide-svelte';

  let profile = $state({});
  let applications = $state([]);
  let loading = $state(true);

  onMount(async () => {
    if (!$auth.isAuthenticated) {
      goto('/login');
      return;
    }

    try {
      const token = $auth.token;
      const [profRes, appsRes] = await Promise.all([
        fetch('/api/v1/profile', { headers: { 'Authorization': 'Bearer ' + token } }),
        fetch('/api/v1/applications', { headers: { 'Authorization': 'Bearer ' + token } })
      ]);

      if (profRes.ok) profile = await profRes.json();
      if (appsRes.ok) {
        const appsData = await appsRes.json();
        applications = appsData.applications || appsData || [];
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  });

  let stats = $derived({
    total: applications.length,
    pending: applications.filter(a => a.status === 'applied').length,
    interview: applications.filter(a => a.status === 'interview').length,
    offer: applications.filter(a => a.status === 'offer').length
  });
</script>

<svelte:head>
  <title>Dashboard - JobApplier</title>
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
  {#if loading}
    <!-- Skeleton loading -->
    <div class="animate-fade-in">
      <div class="skeleton h-8 w-64 mb-2"></div>
      <div class="skeleton h-5 w-48 mb-8"></div>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        {#each [1,2,3,4] as _}
          <div class="skeleton h-28 rounded-xl"></div>
        {/each}
      </div>
      <div class="skeleton h-40 rounded-xl mb-8"></div>
      <div class="skeleton h-60 rounded-xl"></div>
    </div>
  {:else}
    <!-- Welcome header -->
    <div class="mb-8 animate-fade-in">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 class="text-2xl sm:text-3xl font-bold text-foreground">
            Welcome back{profile.name ? ', ' + profile.name : ''}!
          </h1>
          <p class="text-muted-foreground mt-1">Here's an overview of your job search activity</p>
        </div>
        <a href="/jobs" class="no-underline">
          <Button class="shadow-sm">
            <Plus size={16} class="mr-2" />
            Find New Jobs
          </Button>
        </a>
      </div>
    </div>

    <!-- KPI Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <Card class="animate-fade-in-up group hover:shadow-lg transition-all duration-300 border-l-4 border-l-primary">
        <CardContent class="p-5">
          <div class="flex items-center justify-between mb-3">
            <div class="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center group-hover:bg-primary/20 transition-colors">
              <Briefcase size={18} class="text-primary" />
            </div>
            <TrendingUp size={14} class="text-green-500" />
          </div>
          <div class="text-3xl font-bold text-foreground">{stats.total}</div>
          <div class="text-sm text-muted-foreground mt-0.5">Total Applications</div>
        </CardContent>
      </Card>

      <Card class="animate-fade-in-up delay-100 group hover:shadow-lg transition-all duration-300 border-l-4 border-l-blue-500">
        <CardContent class="p-5">
          <div class="flex items-center justify-between mb-3">
            <div class="w-10 h-10 rounded-xl bg-blue-500/10 flex items-center justify-center group-hover:bg-blue-500/20 transition-colors">
              <Clock size={18} class="text-blue-500" />
            </div>
          </div>
          <div class="text-3xl font-bold text-foreground">{stats.pending}</div>
          <div class="text-sm text-muted-foreground mt-0.5">Pending Review</div>
        </CardContent>
      </Card>

      <Card class="animate-fade-in-up delay-200 group hover:shadow-lg transition-all duration-300 border-l-4 border-l-amber-500">
        <CardContent class="p-5">
          <div class="flex items-center justify-between mb-3">
            <div class="w-10 h-10 rounded-xl bg-amber-500/10 flex items-center justify-center group-hover:bg-amber-500/20 transition-colors">
              <MessageSquare size={18} class="text-amber-500" />
            </div>
          </div>
          <div class="text-3xl font-bold text-foreground">{stats.interview}</div>
          <div class="text-sm text-muted-foreground mt-0.5">Interviews</div>
        </CardContent>
      </Card>

      <Card class="animate-fade-in-up delay-300 group hover:shadow-lg transition-all duration-300 border-l-4 border-l-green-500">
        <CardContent class="p-5">
          <div class="flex items-center justify-between mb-3">
            <div class="w-10 h-10 rounded-xl bg-green-500/10 flex items-center justify-center group-hover:bg-green-500/20 transition-colors">
              <Award size={18} class="text-green-500" />
            </div>
          </div>
          <div class="text-3xl font-bold text-foreground">{stats.offer}</div>
          <div class="text-sm text-muted-foreground mt-0.5">Offers Received</div>
        </CardContent>
      </Card>
    </div>

    <!-- Quick Actions -->
    <div class="mb-8 animate-fade-in-up delay-400">
      <h2 class="text-lg font-semibold text-foreground mb-4">Quick Actions</h2>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <a href="/jobs" class="no-underline">
          <Card class="group hover:shadow-md hover:border-primary/30 transition-all duration-300 cursor-pointer">
            <CardContent class="p-5 flex items-center gap-4">
              <div class="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center group-hover:bg-primary/20 transition-colors shrink-0">
                <Search size={18} class="text-primary" />
              </div>
              <div class="flex-1 min-w-0">
                <div class="text-sm font-semibold text-foreground">Search Jobs</div>
                <div class="text-xs text-muted-foreground">Find your next opportunity</div>
              </div>
              <ArrowRight size={16} class="text-muted-foreground group-hover:text-primary transition-colors shrink-0" />
            </CardContent>
          </Card>
        </a>

        <a href="/profile" class="no-underline">
          <Card class="group hover:shadow-md hover:border-primary/30 transition-all duration-300 cursor-pointer">
            <CardContent class="p-5 flex items-center gap-4">
              <div class="w-10 h-10 rounded-xl bg-blue-500/10 flex items-center justify-center group-hover:bg-blue-500/20 transition-colors shrink-0">
                <FileText size={18} class="text-blue-500" />
              </div>
              <div class="flex-1 min-w-0">
                <div class="text-sm font-semibold text-foreground">Update Profile</div>
                <div class="text-xs text-muted-foreground">Improve your applications</div>
              </div>
              <ArrowRight size={16} class="text-muted-foreground group-hover:text-primary transition-colors shrink-0" />
            </CardContent>
          </Card>
        </a>

        <a href="/settings" class="no-underline">
          <Card class="group hover:shadow-md hover:border-primary/30 transition-all duration-300 cursor-pointer">
            <CardContent class="p-5 flex items-center gap-4">
              <div class="w-10 h-10 rounded-xl bg-amber-500/10 flex items-center justify-center group-hover:bg-amber-500/20 transition-colors shrink-0">
                <Settings size={18} class="text-amber-500" />
              </div>
              <div class="flex-1 min-w-0">
                <div class="text-sm font-semibold text-foreground">Configure AI</div>
                <div class="text-xs text-muted-foreground">Set up your preferences</div>
              </div>
              <ArrowRight size={16} class="text-muted-foreground group-hover:text-primary transition-colors shrink-0" />
            </CardContent>
          </Card>
        </a>
      </div>
    </div>

    <!-- Recent Applications -->
    <div class="animate-fade-in-up delay-500">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-lg font-semibold text-foreground">Recent Applications</h2>
        {#if applications.length > 0}
          <a href="/applications" class="text-sm text-primary no-underline hover:underline flex items-center gap-1">
            View all
            <ArrowRight size={14} />
          </a>
        {/if}
      </div>

      {#if applications.length === 0}
        <!-- Empty state -->
        <Card class="border-dashed">
          <CardContent class="p-12 text-center">
            <div class="w-16 h-16 rounded-2xl bg-muted flex items-center justify-center mx-auto mb-4">
              <Briefcase size={28} class="text-muted-foreground" />
            </div>
            <h3 class="text-lg font-semibold text-foreground mb-2">No applications yet</h3>
            <p class="text-sm text-muted-foreground mb-6 max-w-sm mx-auto">
              Get started by searching for jobs. Our AI will help you create tailored applications.
            </p>
            <a href="/jobs" class="no-underline">
              <Button>
                <Search size={16} class="mr-2" />
                Start Searching Jobs
              </Button>
            </a>
          </CardContent>
        </Card>
      {:else}
        <Card>
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr class="border-b border-border">
                  <th class="text-left text-xs font-medium text-muted-foreground uppercase tracking-wider p-4">Company</th>
                  <th class="text-left text-xs font-medium text-muted-foreground uppercase tracking-wider p-4 hidden sm:table-cell">Position</th>
                  <th class="text-left text-xs font-medium text-muted-foreground uppercase tracking-wider p-4">Status</th>
                  <th class="text-left text-xs font-medium text-muted-foreground uppercase tracking-wider p-4 hidden md:table-cell">Date</th>
                </tr>
              </thead>
              <tbody>
                {#each applications.slice(0, 5) as app, i}
                  <tr class="border-b border-border/50 last:border-0 hover:bg-muted/30 transition-colors">
                    <td class="p-4">
                      <div class="flex items-center gap-3">
                        <div class="w-9 h-9 rounded-lg bg-primary/10 flex items-center justify-center text-xs font-bold text-primary shrink-0">
                          {(app.job?.company || 'U')[0].toUpperCase()}
                        </div>
                        <div class="min-w-0">
                          <div class="font-medium text-foreground text-sm truncate">{app.job?.company || 'Unknown'}</div>
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
                      <span class="text-sm text-muted-foreground">
                        {(app.applied_at || app.appliedAt) ? new Date(app.applied_at || app.appliedAt).toLocaleDateString() : '-'}
                      </span>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </Card>
      {/if}
    </div>
  {/if}
</div>
