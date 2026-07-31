<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import KoalaLogo from './KoalaLogo.svelte';
	import type { AppTheme } from '$lib/types/visualizer.types';
	import Sun from 'phosphor-svelte/lib/Sun';
	import Monitor from 'phosphor-svelte/lib/Monitor';
	import Moon from 'phosphor-svelte/lib/Moon';
	import GithubLogo from 'phosphor-svelte/lib/GithubLogo';

	let currentTheme = $state<AppTheme>('system');
	let compact = $state(false);

	onMount(() => {
		const updateHeaderSize = () => {
			const nextCompact = window.scrollY > 0;
			if (compact !== nextCompact) compact = nextCompact;
		};
		updateHeaderSize();
		window.addEventListener('scroll', updateHeaderSize, { passive: true });

		let stored: AppTheme | null = null;
		try {
			stored = localStorage.getItem('koala-theme') as AppTheme | null;
		} catch {
			// Storage can be disabled without preventing the site from rendering.
		}
		if (stored && ['light', 'dark', 'system'].includes(stored)) {
			currentTheme = stored;
		}
		applyTheme(currentTheme);

		return () => {
			window.removeEventListener('scroll', updateHeaderSize);
		};
	});

	function setTheme(theme: AppTheme) {
		currentTheme = theme;
		try {
			localStorage.setItem('koala-theme', theme);
		} catch {
			// The selected theme still applies for this page view.
		}
		applyTheme(theme);
	}

	function applyTheme(theme: AppTheme) {
		if (typeof document === 'undefined') return;
		if (theme === 'system') {
			document.documentElement.removeAttribute('data-theme');
			document.body.removeAttribute('data-theme');
		} else {
			document.documentElement.setAttribute('data-theme', theme);
			document.body.setAttribute('data-theme', theme);
		}
	}

	function isActive(path: string): boolean {
		return path === '/' ? page.url.pathname === '/' : page.url.pathname.startsWith(path);
	}
</script>

<header class="header" class:compact>
	<div class="container header-content">
		<a href="/" class="brand">
			<KoalaLogo size={compact ? 34 : 42} class="brand-logo" />
			<div class="brand-text">
				<span class="title">Koala<span class="highlight">GitHub</span></span>
				<span class="subtitle">Discover GitHub profile visualizers</span>
			</div>
		</a>

		<nav class="nav" aria-label="Primary navigation">
			<a
				href="/"
				class="nav-link"
				class:active={isActive('/')}
				aria-current={isActive('/') ? 'page' : undefined}>Hub</a
			>
			<a
				href="/guide"
				class="nav-link"
				class:active={isActive('/guide')}
				aria-current={isActive('/guide') ? 'page' : undefined}>Guide</a
			>
			<a
				href="/about"
				class="nav-link"
				class:active={isActive('/about')}
				aria-current={isActive('/about') ? 'page' : undefined}>About</a
			>
			<a
				href="/privacy"
				class="nav-link"
				class:active={isActive('/privacy')}
				aria-current={isActive('/privacy') ? 'page' : undefined}>Privacy</a
			>
		</nav>

		<div class="actions">
			<div class="theme-toggle" role="group" aria-label="Theme selector">
				<button
					type="button"
					class="theme-btn"
					class:active={currentTheme === 'light'}
					onclick={() => setTheme('light')}
					title="Light Theme"
					aria-label="Use light theme"
					aria-pressed={currentTheme === 'light'}
				>
					<Sun size={16} weight="regular" aria-hidden="true" />
				</button>
				<button
					type="button"
					class="theme-btn"
					class:active={currentTheme === 'system'}
					onclick={() => setTheme('system')}
					title="System Theme"
					aria-label="Use system theme"
					aria-pressed={currentTheme === 'system'}
				>
					<Monitor size={16} weight="regular" aria-hidden="true" />
				</button>
				<button
					type="button"
					class="theme-btn"
					class:active={currentTheme === 'dark'}
					onclick={() => setTheme('dark')}
					title="Dark Theme"
					aria-label="Use dark theme"
					aria-pressed={currentTheme === 'dark'}
				>
					<Moon size={16} weight="regular" aria-hidden="true" />
				</button>
			</div>

			<a
				href="https://github.com/Shik3i/KoalaGithub"
				target="_blank"
				rel="noopener noreferrer"
				class="github-link"
				title="GitHub Repository"
				aria-label="Open the KoalaGitHub source repository"
			>
				<GithubLogo size={20} weight="fill" aria-hidden="true" />
			</a>
		</div>
	</div>
</header>

<style>
	.header {
		background-color: color-mix(in srgb, var(--bg-card) 94%, transparent);
		border-bottom: 1px solid var(--border-color);
		position: sticky;
		top: 0;
		z-index: 100;
		backdrop-filter: blur(14px);
		transition:
			box-shadow 0.2s ease,
			background-color 0.2s ease;
	}

	.header.compact {
		background-color: color-mix(in srgb, var(--bg-card) 97%, transparent);
		box-shadow: var(--shadow-sm);
	}

	.header-content {
		display: flex;
		align-items: center;
		justify-content: space-between;
		height: 74px;
		gap: 1rem;
		transition: height 0.2s ease;
	}

	.header.compact .header-content {
		height: 54px;
	}

	.brand {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		color: var(--text-main);
		transition: gap 0.2s ease;
	}

	.header.compact .brand {
		gap: 0.55rem;
	}

	:global(.brand-logo) {
		flex: 0 0 auto;
		transition:
			width 0.2s ease,
			height 0.2s ease;
	}

	.brand-text {
		display: flex;
		flex-direction: column;
	}

	.title {
		font-weight: 800;
		font-size: 1.25rem;
		letter-spacing: -0.02em;
		line-height: 1.2;
	}

	.highlight {
		color: var(--brand-primary);
	}

	.subtitle {
		font-size: 0.725rem;
		color: var(--text-muted);
		font-weight: 500;
		line-height: 1.25;
		max-height: 1rem;
		opacity: 1;
		overflow: hidden;
		transform: translateY(0);
		transition:
			max-height 0.18s ease,
			opacity 0.14s ease,
			transform 0.18s ease;
	}

	.header.compact .subtitle {
		max-height: 0;
		opacity: 0;
		transform: translateY(-0.25rem);
	}

	.nav {
		display: flex;
		gap: 1.5rem;
	}

	.nav-link {
		font-weight: 600;
		font-size: 0.925rem;
		color: var(--text-muted);
		padding: 0.4rem 0.6rem;
		border-radius: var(--radius-sm);
		transition:
			background-color 0.15s ease,
			color 0.15s ease;
	}

	.nav-link:hover,
	.nav-link.active {
		color: var(--text-main);
		background-color: var(--bg-subtle);
	}

	.actions {
		display: flex;
		align-items: center;
		gap: 0.85rem;
	}

	.theme-toggle {
		display: flex;
		background-color: var(--bg-subtle);
		padding: 3px;
		border-radius: var(--radius-full);
		border: 1px solid var(--border-color);
	}

	.theme-btn {
		display: inline-grid;
		place-items: center;
		width: 2.25rem;
		height: 2.25rem;
		border-radius: var(--radius-full);
		font-size: 0.85rem;
		transition: background-color 0.2s ease;
		line-height: 1;
		color: var(--text-muted);
	}

	.theme-btn.active {
		background-color: var(--bg-card);
		box-shadow: var(--shadow-sm);
		color: var(--brand-primary);
	}

	.github-link {
		color: var(--text-muted);
		display: flex;
		align-items: center;
		width: 2.5rem;
		height: 2.5rem;
		justify-content: center;
		border-radius: var(--radius-md);
		transition:
			color 0.15s ease,
			background-color 0.15s ease;
	}

	.github-link:hover {
		color: var(--text-main);
		background-color: var(--bg-subtle);
	}

	@media (max-width: 640px) {
		.header-content {
			height: auto;
			min-height: 70px;
			padding: 0.65rem 0;
			flex-wrap: wrap;
			row-gap: 0.45rem;
		}

		.header.compact .header-content {
			height: auto;
			min-height: 54px;
			padding: 0.35rem 0;
			row-gap: 0.25rem;
		}

		.nav {
			order: 3;
			width: 100%;
			justify-content: center;
			gap: 0.5rem;
		}

		.header.compact .nav-link {
			padding-block: 0.2rem;
			font-size: 0.825rem;
		}

		.actions {
			margin-left: auto;
		}
	}
</style>
