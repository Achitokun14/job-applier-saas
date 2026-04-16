<script lang="ts">
  import "../app.css";
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { ModeWatcher, toggleMode } from 'mode-watcher';
  import { Sun, Moon, LayoutDashboard, Search, FileText, User, Settings } from 'lucide-svelte';
  import { Toaster } from 'svelte-sonner';
  import { Button } from '$lib/components/ui/button';

  let { children } = $props();

  let isAuthenticated = $derived($auth.isAuthenticated);

  function handleLogout() {
    auth.logout();
    goto('/login');
  }
</script>

<ModeWatcher />
<Toaster />

<div class="min-h-screen flex flex-col bg-background text-foreground">
  <nav class="flex justify-between items-center px-8 py-3 bg-background/90 backdrop-blur-xl border-b border-border sticky top-0 z-50">
    <div class="flex items-center">
      <a href="/" class="flex items-center gap-2 text-xl font-bold text-primary no-underline">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
          <path d="M20 6L12 2L4 6V18L12 22L20 18V6Z" stroke="currentColor" stroke-width="2" fill="none"/>
          <path d="M12 22V12" stroke="currentColor" stroke-width="2"/>
          <path d="M20 6L12 10L4 6" stroke="currentColor" stroke-width="2"/>
        </svg>
        JobApplier
      </a>
    </div>

    <div class="flex gap-1 items-center">
      {#if isAuthenticated}
        <a href="/dashboard" class="flex items-center gap-1.5 px-3 py-2 text-muted-foreground no-underline text-sm font-medium rounded-md hover:text-primary hover:bg-accent transition-colors">
          <LayoutDashboard size={18} />
          Dashboard
        </a>
        <a href="/jobs" class="flex items-center gap-1.5 px-3 py-2 text-muted-foreground no-underline text-sm font-medium rounded-md hover:text-primary hover:bg-accent transition-colors">
          <Search size={18} />
          Jobs
        </a>
        <a href="/applications" class="flex items-center gap-1.5 px-3 py-2 text-muted-foreground no-underline text-sm font-medium rounded-md hover:text-primary hover:bg-accent transition-colors">
          <FileText size={18} />
          Applications
        </a>
        <a href="/profile" class="flex items-center gap-1.5 px-3 py-2 text-muted-foreground no-underline text-sm font-medium rounded-md hover:text-primary hover:bg-accent transition-colors">
          <User size={18} />
          Profile
        </a>
        <a href="/settings" class="flex items-center gap-1.5 px-3 py-2 text-muted-foreground no-underline text-sm font-medium rounded-md hover:text-primary hover:bg-accent transition-colors">
          <Settings size={18} />
          Settings
        </a>
        <button
          onclick={toggleMode}
          class="inline-flex items-center justify-center h-9 w-9 rounded-md text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors ml-1"
          aria-label="Toggle dark mode"
        >
          <Sun size={18} class="dark:hidden" />
          <Moon size={18} class="hidden dark:block" />
        </button>
        <Button variant="outline" size="sm" onclick={handleLogout} class="ml-2">
          Logout
        </Button>
      {:else}
        <button
          onclick={toggleMode}
          class="inline-flex items-center justify-center h-9 w-9 rounded-md text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors"
          aria-label="Toggle dark mode"
        >
          <Sun size={18} class="dark:hidden" />
          <Moon size={18} class="hidden dark:block" />
        </button>
        <a href="/login" class="px-4 py-2 text-muted-foreground no-underline font-medium text-sm hover:text-primary transition-colors">Login</a>
        <a href="/register" class="px-4 py-2 bg-primary text-primary-foreground no-underline font-semibold text-sm rounded-md hover:bg-primary/90 transition-colors">Get Started</a>
      {/if}
    </div>
  </nav>

  <main class="flex-1 p-8 max-w-[1280px] mx-auto w-full">
    {@render children()}
  </main>
</div>
