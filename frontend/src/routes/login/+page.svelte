<script>
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';

  let email = $state('');
  let password = $state('');
  let loading = $state(false);
  let error = $state('');

  async function login() {
    if (!email || !password) {
      error = 'Please fill in all fields';
      return;
    }

    loading = true;
    error = '';

    try {
      const res = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });

      if (res.ok) {
        const data = await res.json();
        auth.login(data.token, data.user);
        goto('/dashboard');
      } else {
        error = 'Invalid email or password';
      }
    } catch (e) {
      error = 'Login failed';
    } finally {
      loading = false;
    }
  }
</script>

<div class="flex justify-center items-center min-h-[60vh] px-4">
  <Card class="w-full max-w-[400px]">
    <CardHeader class="text-center">
      <CardTitle class="text-2xl">Welcome Back</CardTitle>
      <CardDescription>Sign in to continue</CardDescription>
    </CardHeader>
    <CardContent>
      {#if error}
        <div class="bg-destructive/10 text-destructive text-sm p-3 rounded-md mb-4">
          {error}
        </div>
      {/if}

      <form onsubmit={(e) => { e.preventDefault(); login(); }} class="space-y-4">
        <div class="space-y-2">
          <label for="email" class="text-sm font-medium text-foreground">Email</label>
          <Input type="email" id="email" bind:value={email} placeholder="you@example.com" />
        </div>

        <div class="space-y-2">
          <label for="password" class="text-sm font-medium text-foreground">Password</label>
          <Input type="password" id="password" bind:value={password} placeholder="********" />
        </div>

        <Button type="submit" class="w-full" disabled={loading}>
          {loading ? 'Signing in...' : 'Sign In'}
        </Button>
      </form>
    </CardContent>
    <CardFooter class="justify-center">
      <p class="text-sm text-muted-foreground">
        Don't have an account? <a href="/register" class="text-primary font-medium hover:underline">Sign up</a>
      </p>
    </CardFooter>
  </Card>
</div>
