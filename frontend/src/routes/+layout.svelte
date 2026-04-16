<script lang="ts">
  import "../app.css";
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { ModeWatcher, toggleMode } from 'mode-watcher';
  import { Sun, Moon, LayoutDashboard, Search, FileText, User, Settings, Menu, X, LogOut, ChevronDown } from 'lucide-svelte';
  import { Toaster } from 'svelte-sonner';
  import { Button } from '$lib/components/ui/button';

  let { children } = $props();

  let isAuthenticated = $derived($auth.isAuthenticated);
  let userName = $derived($auth.user?.name || '');
  let userInitials = $derived(
    userName ? userName.split(' ').map((n: string) => n[0]).join('').toUpperCase().slice(0, 2) : 'U'
  );
  let currentPath = $derived($page.url.pathname);

  let mobileMenuOpen = $state(false);
  let userMenuOpen = $state(false);

  function handleLogout() {
    auth.logout();
    userMenuOpen = false;
    mobileMenuOpen = false;
    goto('/login');
  }

  function closeMobileMenu() {
    mobileMenuOpen = false;
  }

  function toggleUserMenu() {
    userMenuOpen = !userMenuOpen;
  }

  function closeUserMenu() {
    userMenuOpen = false;
  }

  const navLinks = [
    { href: '/dashboard', label: 'Dashboard', icon: LayoutDashboard },
    { href: '/jobs', label: 'Jobs', icon: Search },
    { href: '/applications', label: 'Applications', icon: FileText },
    { href: '/profile', label: 'Profile', icon: User },
    { href: '/settings', label: 'Settings', icon: Settings },
  ];

  function isActive(path: string): boolean {
    return currentPath === path;
  }
</script>

<ModeWatcher />
<Toaster />

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="min-h-screen flex flex-col bg-background text-foreground">
  <!-- Navbar -->
  <nav class="glass sticky top-0 z-50 border-b border-border/50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo -->
        <a href="/" class="flex items-center gap-2.5 text-xl font-bold text-primary no-underline shrink-0 group">
          <div class="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center group-hover:bg-primary/20 transition-colors">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
              <path d="M20 6L12 2L4 6V18L12 22L20 18V6Z" stroke="currentColor" stroke-width="2.5" fill="none"/>
              <path d="M12 22V12" stroke="currentColor" stroke-width="2"/>
              <path d="M20 6L12 10L4 6" stroke="currentColor" stroke-width="2"/>
            </svg>
          </div>
          <span class="hidden sm:inline">JobApplier</span>
        </a>

        <!-- Center nav links (desktop) -->
        {#if isAuthenticated}
          <div class="hidden md:flex items-center gap-1">
            {#each navLinks as link}
              <a
                href={link.href}
                class="flex items-center gap-2 px-3 py-2 text-sm font-medium rounded-lg no-underline transition-all duration-200
                  {isActive(link.href)
                    ? 'text-primary bg-primary/10'
                    : 'text-muted-foreground hover:text-foreground hover:bg-accent'}"
              >
                <svelte:component this={link.icon} size={16} />
                {link.label}
              </a>
            {/each}
          </div>
        {/if}

        <!-- Right side -->
        <div class="flex items-center gap-2">
          <!-- Dark mode toggle -->
          <button
            onclick={toggleMode}
            class="inline-flex items-center justify-center h-9 w-9 rounded-lg text-muted-foreground hover:bg-accent hover:text-foreground transition-all duration-200"
            aria-label="Toggle dark mode"
          >
            <Sun size={16} class="dark:hidden" />
            <Moon size={16} class="hidden dark:block" />
          </button>

          {#if isAuthenticated}
            <!-- User menu (desktop) -->
            <div class="hidden md:block relative">
              <button
                onclick={toggleUserMenu}
                class="flex items-center gap-2 pl-2 pr-3 py-1.5 rounded-lg hover:bg-accent transition-all duration-200 cursor-pointer"
              >
                <div class="w-8 h-8 rounded-full bg-primary/10 text-primary flex items-center justify-center text-xs font-bold">
                  {userInitials}
                </div>
                <span class="text-sm font-medium text-foreground max-w-[120px] truncate">{userName || 'User'}</span>
                <ChevronDown size={14} class="text-muted-foreground" />
              </button>

              {#if userMenuOpen}
                <div
                  class="absolute right-0 top-full mt-2 w-48 bg-card rounded-xl border border-border shadow-lg py-1 z-50"
                >
                  <div class="px-3 py-2 border-b border-border">
                    <p class="text-sm font-medium text-foreground truncate">{userName || 'User'}</p>
                    <p class="text-xs text-muted-foreground truncate">{$auth.user?.email || ''}</p>
                  </div>
                  <a href="/profile" onclick={closeUserMenu} class="flex items-center gap-2 px-3 py-2 text-sm text-foreground no-underline hover:bg-accent transition-colors">
                    <User size={14} />
                    Profile
                  </a>
                  <a href="/settings" onclick={closeUserMenu} class="flex items-center gap-2 px-3 py-2 text-sm text-foreground no-underline hover:bg-accent transition-colors">
                    <Settings size={14} />
                    Settings
                  </a>
                  <div class="border-t border-border mt-1 pt-1">
                    <button
                      onclick={handleLogout}
                      class="flex items-center gap-2 px-3 py-2 text-sm text-destructive hover:bg-destructive/10 transition-colors w-full text-left cursor-pointer"
                    >
                      <LogOut size={14} />
                      Logout
                    </button>
                  </div>
                </div>
                <!-- Backdrop for user menu -->
                <div class="fixed inset-0 z-[-1]" onclick={closeUserMenu}></div>
              {/if}
            </div>

            <!-- Mobile hamburger -->
            <button
              onclick={() => mobileMenuOpen = !mobileMenuOpen}
              class="md:hidden inline-flex items-center justify-center h-9 w-9 rounded-lg text-muted-foreground hover:bg-accent hover:text-foreground transition-all duration-200"
              aria-label="Toggle menu"
            >
              {#if mobileMenuOpen}
                <X size={18} />
              {:else}
                <Menu size={18} />
              {/if}
            </button>
          {:else}
            <a href="/login" class="hidden sm:inline-flex px-4 py-2 text-sm font-medium text-muted-foreground no-underline hover:text-foreground transition-colors">
              Sign In
            </a>
            <a href="/register" class="px-4 py-2 bg-primary text-primary-foreground no-underline font-semibold text-sm rounded-lg hover:bg-primary/90 transition-all duration-200 shadow-sm hover:shadow-md">
              Get Started
            </a>
            <!-- Mobile hamburger for unauthenticated too -->
            <button
              onclick={() => mobileMenuOpen = !mobileMenuOpen}
              class="sm:hidden inline-flex items-center justify-center h-9 w-9 rounded-lg text-muted-foreground hover:bg-accent hover:text-foreground transition-all duration-200"
              aria-label="Toggle menu"
            >
              {#if mobileMenuOpen}
                <X size={18} />
              {:else}
                <Menu size={18} />
              {/if}
            </button>
          {/if}
        </div>
      </div>
    </div>

    <!-- Mobile menu -->
    {#if mobileMenuOpen}
      <div class="md:hidden border-t border-border/50 bg-background/95 backdrop-blur-lg">
        <div class="px-4 py-3 space-y-1">
          {#if isAuthenticated}
            <!-- User info -->
            <div class="flex items-center gap-3 px-3 py-3 mb-2 border-b border-border/50">
              <div class="w-10 h-10 rounded-full bg-primary/10 text-primary flex items-center justify-center text-sm font-bold">
                {userInitials}
              </div>
              <div>
                <p class="text-sm font-medium text-foreground">{userName || 'User'}</p>
                <p class="text-xs text-muted-foreground">{$auth.user?.email || ''}</p>
              </div>
            </div>

            {#each navLinks as link}
              <a
                href={link.href}
                onclick={closeMobileMenu}
                class="flex items-center gap-3 px-3 py-2.5 text-sm font-medium rounded-lg no-underline transition-all duration-200
                  {isActive(link.href)
                    ? 'text-primary bg-primary/10'
                    : 'text-muted-foreground hover:text-foreground hover:bg-accent'}"
              >
                <svelte:component this={link.icon} size={18} />
                {link.label}
              </a>
            {/each}

            <div class="border-t border-border/50 mt-2 pt-2">
              <button
                onclick={handleLogout}
                class="flex items-center gap-3 px-3 py-2.5 text-sm font-medium text-destructive hover:bg-destructive/10 rounded-lg transition-colors w-full text-left cursor-pointer"
              >
                <LogOut size={18} />
                Logout
              </button>
            </div>
          {:else}
            <a
              href="/login"
              onclick={closeMobileMenu}
              class="flex items-center gap-3 px-3 py-2.5 text-sm font-medium text-muted-foreground no-underline hover:text-foreground hover:bg-accent rounded-lg transition-colors"
            >
              Sign In
            </a>
            <a
              href="/register"
              onclick={closeMobileMenu}
              class="flex items-center justify-center gap-2 px-3 py-2.5 text-sm font-semibold bg-primary text-primary-foreground no-underline rounded-lg hover:bg-primary/90 transition-colors"
            >
              Get Started
            </a>
          {/if}
        </div>
      </div>
    {/if}
  </nav>

  <!-- Main content -->
  <main class="flex-1">
    {@render children()}
  </main>

  <!-- Footer -->
  <footer class="border-t border-border/50 bg-card/50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
        <!-- Brand -->
        <div class="md:col-span-1">
          <a href="/" class="flex items-center gap-2 text-lg font-bold text-primary no-underline mb-3">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M20 6L12 2L4 6V18L12 22L20 18V6Z" stroke="currentColor" stroke-width="2.5" fill="none"/>
              <path d="M12 22V12" stroke="currentColor" stroke-width="2"/>
              <path d="M20 6L12 10L4 6" stroke="currentColor" stroke-width="2"/>
            </svg>
            JobApplier
          </a>
          <p class="text-sm text-muted-foreground leading-relaxed">
            AI-powered job application platform. Automate your job search and land your dream role.
          </p>
        </div>

        <!-- Product -->
        <div>
          <h4 class="text-sm font-semibold text-foreground mb-3">Product</h4>
          <div class="flex flex-col gap-2">
            <a href="/jobs" class="text-sm text-muted-foreground no-underline hover:text-foreground transition-colors">Job Search</a>
            <a href="/dashboard" class="text-sm text-muted-foreground no-underline hover:text-foreground transition-colors">Dashboard</a>
            <a href="/applications" class="text-sm text-muted-foreground no-underline hover:text-foreground transition-colors">Applications</a>
            <span class="text-sm text-muted-foreground">AI Resume Builder</span>
          </div>
        </div>

        <!-- Company -->
        <div>
          <h4 class="text-sm font-semibold text-foreground mb-3">Company</h4>
          <div class="flex flex-col gap-2">
            <span class="text-sm text-muted-foreground">About</span>
            <span class="text-sm text-muted-foreground">Careers</span>
            <span class="text-sm text-muted-foreground">Blog</span>
            <span class="text-sm text-muted-foreground">Contact</span>
          </div>
        </div>

        <!-- Legal -->
        <div>
          <h4 class="text-sm font-semibold text-foreground mb-3">Legal</h4>
          <div class="flex flex-col gap-2">
            <span class="text-sm text-muted-foreground">Privacy Policy</span>
            <span class="text-sm text-muted-foreground">Terms of Service</span>
            <span class="text-sm text-muted-foreground">Cookie Policy</span>
          </div>
        </div>
      </div>

      <div class="mt-10 pt-6 border-t border-border/50 flex flex-col sm:flex-row justify-between items-center gap-4">
        <p class="text-sm text-muted-foreground">&copy; 2026 JobApplier. All rights reserved.</p>
        <div class="flex items-center gap-4">
          <span class="text-sm text-muted-foreground">Built with AI</span>
        </div>
      </div>
    </div>
  </footer>
</div>
