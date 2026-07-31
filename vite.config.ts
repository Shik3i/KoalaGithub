/// <reference types="vitest" />
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

const appVersion =
	process.env.PUBLIC_APP_VERSION?.trim() || process.env.npm_package_version || 'dev';

export default defineConfig({
	plugins: [sveltekit()],
	define: {
		__APP_VERSION__: JSON.stringify(appVersion)
	},
	// @ts-expect-error vitest config block
	test: {
		include: ['src/**/*.test.ts', 'tests/**/*.test.ts']
	}
});
