/// <reference types="vitest" />
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	// @ts-expect-error vitest config block
	test: {
		include: ['src/**/*.test.ts', 'tests/**/*.test.ts']
	}
});
