/// <reference types="@sveltejs/kit" />

declare namespace App {
  interface Locals {
    user?: {
      id: number;
      email: string;
      name: string;
    };
  }
}

declare module '$env/static/public' {
  export const PUBLIC_API_URL: string;
}
