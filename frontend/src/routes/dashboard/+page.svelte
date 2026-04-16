<script>
  import { auth } from '$lib/stores/auth';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Card, CardHeader, CardTitle, CardContent } from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Briefcase, Clock, MessageSquare, Award, Search, FileText, User, Settings } from 'lucide-svelte';

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
      if (appsRes.ok) applications = await appsRes.json();
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

<div class="max-w-[900px] mx-auto">
  {#if loading}
    <div class="text-center py-12 text-muted-foreground">Loading...</div>
  {:else}
    <div class="flex justify-between items-center mb-8">
      <div>
        <h1 class="text-3xl font-bold text-foreground m-0">Welcome{profile.name ? ', ' + profile.name : ''}!</h1>
        <p class="text-muted-foreground mt-1 mb-0">Here's your job search overview</p>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <Card class="bg-primary text-primary-foreground border-primary">
        <CardContent class="p-6">
          <div class="flex items-center justify-between mb-2">
            <Briefcase size={20} class="opacity-80" />
          </div>
          <div class="text-3xl font-bold">{stats.total}</div>
          <div class="text-sm opacity-80">Total Applications</div>
        </CardContent>
      </Card>
      <Card>
        <CardContent class="p-6">
          <div class="flex items-center justify-between mb-2">
            <Clock size={20} class="text-blue-500" />
          </div>
          <div class="text-3xl font-bold text-foreground">{stats.pending}</div>
          <div class="text-sm text-muted-foreground">Pending</div>
        </CardContent>
      </Card>
      <Card>
        <CardContent class="p-6">
          <div class="flex items-center justify-between mb-2">
            <MessageSquare size={20} class="text-yellow-500" />
          </div>
          <div class="text-3xl font-bold text-foreground">{stats.interview}</div>
          <div class="text-sm text-muted-foreground">Interviews</div>
        </CardContent>
      </Card>
      <Card>
        <CardContent class="p-6">
          <div class="flex items-center justify-between mb-2">
            <Award size={20} class="text-green-500" />
          </div>
          <div class="text-3xl font-bold text-foreground">{stats.offer}</div>
          <div class="text-sm text-muted-foreground">Offers</div>
        </CardContent>
      </Card>
    </div>

    <div class="mb-8">
      <h2 class="text-lg font-semibold text-foreground mb-4">Quick Actions</h2>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
        <a href="/jobs" class="no-underline">
          <Button variant="outline" class="w-full h-auto py-4 flex flex-col gap-2">
            <Search size={18} />
            <span>Search Jobs</span>
          </Button>
        </a>
        <a href="/applications" class="no-underline">
          <Button variant="outline" class="w-full h-auto py-4 flex flex-col gap-2">
            <FileText size={18} />
            <span>View All</span>
          </Button>
        </a>
        <a href="/profile" class="no-underline">
          <Button variant="outline" class="w-full h-auto py-4 flex flex-col gap-2">
            <User size={18} />
            <span>Edit Profile</span>
          </Button>
        </a>
        <a href="/settings" class="no-underline">
          <Button variant="outline" class="w-full h-auto py-4 flex flex-col gap-2">
            <Settings size={18} />
            <span>Settings</span>
          </Button>
        </a>
      </div>
    </div>

    <div>
      <h2 class="text-lg font-semibold text-foreground mb-4">Recent Applications</h2>
      {#if applications.length === 0}
        <p class="text-muted-foreground">No applications yet. <a href="/jobs" class="text-primary hover:underline">Find jobs</a></p>
      {:else}
        <div class="flex flex-col gap-2">
          {#each applications.slice(0, 5) as app}
            <Card>
              <CardContent class="p-4 flex justify-between items-center">
                <div>
                  <div class="font-semibold text-foreground">{app.job?.title || 'Unknown'}</div>
                  <div class="text-sm text-muted-foreground">{app.job?.company || ''}</div>
                </div>
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
              </CardContent>
            </Card>
          {/each}
        </div>
      {/if}
    </div>
  {/if}
</div>
