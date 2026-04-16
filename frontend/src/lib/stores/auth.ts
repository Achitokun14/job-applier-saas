import { writable } from 'svelte/store';
import { browser } from '$app/environment';

interface User {
  id: number;
  email: string;
  name: string;
}

interface AuthState {
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
}

function createAuthStore() {
  let initial: AuthState = {
    token: null,
    user: null,
    isAuthenticated: false
  };

  if (browser) {
    const saved = localStorage.getItem('auth');
    if (saved) {
      try {
        initial = JSON.parse(saved);
      } catch {
        localStorage.removeItem('auth');
      }
    }
  }

  const { subscribe, set, update } = writable<AuthState>(initial);

  return {
    subscribe,
    login: (token: string, user: User) => {
      const state: AuthState = { token, user, isAuthenticated: true };
      set(state);
      if (browser) {
        localStorage.setItem('auth', JSON.stringify(state));
      }
    },
    logout: () => {
      const state: AuthState = { token: null, user: null, isAuthenticated: false };
      set(state);
      if (browser) {
        localStorage.removeItem('auth');
      }
    },
    getToken: () => {
      let token: string | null = null;
      update(state => {
        token = state.token;
        return state;
      });
      return token;
    }
  };
}

export const auth = createAuthStore();
