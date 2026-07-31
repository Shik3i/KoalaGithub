<script lang="ts">
	import { onMount } from 'svelte';
	import KoalaLogo from './KoalaLogo.svelte';
	import type { AppTheme } from '$lib/types/visualizer.types';

	let currentTheme = $state<AppTheme>('system');

	onMount(() => {
		const stored = localStorage.getItem('koala-theme') as AppTheme | null;
		if (stored && ['light', 'dark', 'system'].includes(stored)) {
			currentTheme = stored;
		}
		applyTheme(currentTheme);
	});

	function setTheme(theme: AppTheme) {
		currentTheme = theme;
		localStorage.setItem('koala-theme', theme);
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
</script>

<header class="header">
	<div class="container header-content">
		<a href="/" class="brand">
			<KoalaLogo size={42} />
			<div class="brand-text">
				<span class="title">Koala<span class="highlight">GitHub</span></span>
				<span class="subtitle">github.koalastuff.net</span>
			</div>
		</a>

		<nav class="nav">
			<a href="/" class="nav-link">Hub</a>
			<a href="/guide" class="nav-link">Guide 📖</a>
			<a href="/about" class="nav-link">About</a>
			<a href="/privacy" class="nav-link">Privacy</a>
		</nav>

		<div class="actions">
			<div class="theme-toggle" role="group" aria-label="Theme selector">
				<button
					type="button"
					class="theme-btn"
					class:active={currentTheme === 'light'}
					onclick={() => setTheme('light')}
					title="Light Theme"
				>
					☀️
				</button>
				<button
					type="button"
					class="theme-btn"
					class:active={currentTheme === 'system'}
					onclick={() => setTheme('system')}
					title="System Theme"
				>
					💻
				</button>
				<button
					type="button"
					class="theme-btn"
					class:active={currentTheme === 'dark'}
					onclick={() => setTheme('dark')}
					title="Dark Theme"
				>
					🌙
				</button>
			</div>

			<a
				href="https://github.com/Shik3i/KoalaGithub"
				target="_blank"
				rel="noreferrer"
				class="github-link"
				title="GitHub Repository"
			>
				<svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
					<path
						fill-rule="evenodd"
						clip-rule="evenodd"
						d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.53 1.032 1.53 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
					/>
				</svg>
			</a>
		</div>
	</div>
</header>

<style>
	.header {
		background-color: var(--bg-card);
		border-bottom: 1px solid var(--border-color);
		position: sticky;
		top: 0;
		z-index: 100;
		backdrop-filter: blur(8px);
	}

	.header-content {
		display: flex;
		align-items: center;
		justify-content: space-between;
		height: 70px;
		gap: 1rem;
	}

	.brand {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		color: var(--text-main);
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
		transition: all 0.15s ease;
	}

	.nav-link:hover {
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
		padding: 4px 8px;
		border-radius: var(--radius-full);
		font-size: 0.85rem;
		transition: background-color 0.2s ease;
		line-height: 1;
	}

	.theme-btn.active {
		background-color: var(--bg-card);
		box-shadow: var(--shadow-sm);
	}

	.github-link {
		color: var(--text-muted);
		display: flex;
		align-items: center;
		padding: 6px;
		border-radius: var(--radius-md);
		transition: color 0.15s ease, background-color 0.15s ease;
	}

	.github-link:hover {
		color: var(--text-main);
		background-color: var(--bg-subtle);
	}

	@media (max-width: 640px) {
		.header-content {
			height: auto;
			padding: 0.75rem 0;
			flex-wrap: wrap;
		}

		.nav {
			gap: 0.5rem;
		}
	}
</style>
