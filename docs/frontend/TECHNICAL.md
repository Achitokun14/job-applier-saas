# Frontend Technical Documentation

## Overview

The frontend is a SvelteKit application using Bun as the runtime, providing a modern web dashboard for managing job applications.

## Directory Structure

```
frontend/
├── src/
│   ├── app.html                    # HTML template
│   ├── app.d.ts                    # TypeScript declarations
│   ├── routes/
│   │   ├── +layout.svelte          # Main layout with navigation
│   │   ├── +page.svelte            # Landing page
│   │   ├── login/
│   │   │   └── +page.svelte        # Login page
│   │   ├── register/
│   │   │   └── +page.svelte        # Registration page
│   │   ├── dashboard/
│   │   │   └── +page.svelte        # Dashboard with stats
│   │   ├── jobs/
│   │   │   └── +page.svelte        # Job search
│   │   ├── applications/
│   │   │   └── +page.svelte        # Application tracking
│   │   ├── profile/
│   │   │   └── +page.svelte        # Profile editor
│   │   └── settings/
│   │       └── +page.svelte        # Settings page
│   └── lib/
│       ├── api/
│       │   └── index.ts            # API client
│       └── stores/
│           └── auth.ts             # Auth store
├── static/                         # Static assets
├── package.json                    # Dependencies
├── svelte.config.js                # SvelteKit config
├── vite.config.ts                  # Vite config
└── tsconfig.json                   # TypeScript config
```

## Configuration

### SvelteKit Config (`svelte.config.js`)

```javascript
import adapter from '@sveltejs/adapter-node';

const config = {
  kit: {
    adapter: adapter({
      out: 'build',
      precompress: true
    }),
    alias: {
      '$lib': 'src/lib',
      '$lib/*': 'src/lib/*'
    }
  }
};
```

### Vite Config (`vite.config.ts`)

```typescript
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    port: 5173,
    host: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
});
```

## State Management

### Auth Store (`src/lib/stores/auth.ts`)

Manages authentication state using Svelte stores:

```typescript
interface AuthState {
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
}
```

**Features:**
- Persists to localStorage
- Provides `login()`, `logout()`, `getToken()` methods
- Reactive subscriptions via `$auth`

## API Client (`src/lib/api/index.ts`)

Centralized API client with typed methods:

```typescript
export const api = {
  // Auth
  login(email: string, password: string): Promise<AuthResponse>
  register(email: string, password: string, name: string): Promise<AuthResponse>
  
  // Jobs
  searchJobs(query: string, token: string): Promise<Job[]>
  applyJob(jobId: number, token: string): Promise<void>
  
  // Applications
  getApplications(token: string): Promise<Application[]>
  getApplication(id: number, token: string): Promise<Application>
  deleteApplication(id: number, token: string): Promise<void>
  
  // Profile
  getProfile(token: string): Promise<Profile>
  updateProfile(data: ProfileData, token: string): Promise<void>
  
  // Settings
  getSettings(token: string): Promise<Settings>
  updateSettings(data: SettingsData, token: string): Promise<void>
  
  // Resume
  generateResume(yaml: string, style: string, token: string): Promise<ResumeResponse>
  generateCoverLetter(resume: string, jobDesc: string, token: string): Promise<CoverLetterResponse>
};
```

## Routes

### Landing Page (`/`)

- Hero section with CTA buttons
- Feature highlights
- Conditional rendering based on auth state

### Login (`/login`)

- Email/password form
- Error handling
- Redirect to dashboard on success
- Link to register

### Register (`/register`)

- Name/email/password form
- Password validation (min 6 chars)
- Auto-login after registration

### Dashboard (`/dashboard`)

- Statistics cards (total, pending, interviews, offers)
- Quick action buttons
- Recent applications list
- Requires authentication

### Jobs (`/jobs`)

- Search input with query
- Job cards with apply buttons
- Applied status tracking
- Link to original posting

### Applications (`/applications`)

- Table view of all applications
- Status badges with colors
- Delete functionality
- Date formatting

### Profile (`/profile`)

- Personal information form
- Education textarea
- Experience textarea
- Skills textarea
- Projects textarea
- Save functionality

### Settings (`/settings`)

- LLM provider selection
- Model configuration
- API key input (password field)
- Job search preferences
- Experience level dropdown
- Job types dropdown
- Location settings

## Components

### Layout (`+layout.svelte`)

Global layout with:
- Navigation bar
- Logo
- Conditional nav links (auth state)
- Logout button
- Main content area

### Styling

Uses scoped CSS with:
- CSS variables for theming
- Responsive design
- Purple accent color (#7D56F4)
- Clean, modern aesthetic

## Authentication Flow

1. User submits login/register form
2. API call to backend
3. JWT token received
4. Token stored in localStorage via auth store
5. Subsequent requests include Authorization header
6. Protected routes check `$auth.isAuthenticated`

## Error Handling

- Network errors caught in try/catch
- API errors displayed in UI
- Form validation before submission
- Loading states during requests

## Environment Variables

| Variable | Description |
|----------|-------------|
| `PUBLIC_API_URL` | Backend API URL |

## Dependencies

```json
{
  "@sveltejs/adapter-node": "^2.0.0",
  "@sveltejs/kit": "^2.0.0",
  "svelte": "^4.2.7",
  "vite": "^5.0.3",
  "typescript": "^5.0.0"
}
```

## Development

```bash
# Install dependencies
bun install

# Start dev server
bun run dev

# Type checking
bun run check

# Linting
bun run lint

# Build for production
bun run build
```

## Production Build

```bash
bun run build
# Output in build/ directory
# Serve with: node build
```

## Browser Support

- Modern browsers (Chrome, Firefox, Safari, Edge)
- ES2020+ features
- CSS Grid and Flexbox
