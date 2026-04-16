<script>
  import { goto } from '$app/navigation';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '$lib/components/ui/card';
  import { Input } from '$lib/components/ui/input';

  let name = $state('');
  let email = $state('');
  let password = $state('');
  let loading = $state(false);
  let error = $state('');

  async function register() {
    if (!name || !email || !password) {
      error = 'Please fill in all fields';
      return;
    }

    loading = true;
    error = '';

    try {
      const res = await fetch('/api/v1/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, email, password })
      });

      if (res.ok) {
        goto('/login');
      } else {
        error = 'Registration failed';
      }
    } catch (e) {
      error = 'Registration failed';
    } finally {
      loading = false;
    }
  }
</script>

<div class="flex justify-center items-center min-h-[60vh] px-4">
  <Card class="w-full max-w-[400px]">
    <CardHeader class="text-center">
      <CardTitle class="text-2xl">Create Account</CardTitle>
      <CardDescription>Start your job search automation</CardDescription>
    </CardHeader>
    <CardContent>
      {#if error}
        <div class="bg-destructive/10 text-destructive text-sm p-3 rounded-md mb-4">
          {error}
        </div>
      {/if}

      <form onsubmit={(e) => { e.preventDefault(); register(); }} class="space-y-4">
        <div class="space-y-2">
          <label for="name" class="text-sm font-medium text-foreground">Full Name</label>
          <Input type="text" id="name" bind:value={name} placeholder="John Doe" />
        </div>

        <div class="space-y-2">
          <label for="email" class="text-sm font-medium text-foreground">Email</label>
          <Input type="email" id="email" bind:value={email} placeholder="you@example.com" />
        </div>

        <div class="space-y-2">
          <label for="password" class="text-sm font-medium text-foreground">Password</label>
          <Input type="password" id="password" bind:value={password} placeholder="********" />
        </div>

        <Button type="submit" class="w-full" disabled={loading}>
          {loading ? 'Creating...' : 'Create Account'}
        </Button>
      </form>
    </CardContent>
    <CardFooter class="justify-center">
      <p class="text-sm text-muted-foreground">
        Already have an account? <a href="/login" class="text-primary font-medium hover:underline">Sign in</a>
      </p>
    </CardFooter>
  </Card>
</div>
