<script lang="ts">
	import Lock from 'phosphor-svelte/lib/Lock';
	import KoalaLogo from './KoalaLogo.svelte';
</script>

<footer class="footer">
	<div class="container footer-content">
		<div class="footer-brand">
			<div class="logo-row">
				<KoalaLogo size={32} />
				<span class="brand-title">KoalaGitHub</span>
			</div>
			<p class="tagline">
				The central hub for GitHub profile visualizers, stats, trophies, graphs & badges.
			</p>
		</div>

		<div class="footer-links">
			<a href="/guide">Setup Guide</a>
			<a href="/about">About KoalaGitHub</a>
			<a href="/privacy">Privacy Policy</a>
			<a href="https://koalastuff.net/legal" target="_blank" rel="noopener noreferrer"
				>Legal Notice ↗</a
			>
		</div>
	</div>

	<div class="container footer-bottom">
		<p class="privacy-badge">
			<Lock size={14} weight="duotone" aria-hidden="true" /> No analytics • No ad tracking • Open source
		</p>
		<p class="copyright">
			<span class="version">Version {__APP_VERSION__}</span>
			<span>© {new Date().getFullYear()} KoalaStuff Ecosystem</span>
		</p>
	</div>
</footer>

<style>
	.footer {
		margin-top: auto;
		background-color: var(--bg-card);
		border-top: 1px solid var(--border-color);
		padding: 2.5rem 0 1.5rem;
	}

	.footer-content {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 2rem;
		padding-bottom: 2rem;
		border-bottom: 1px solid var(--border-color);
	}

	.footer-brand {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		max-width: 420px;
	}

	.logo-row {
		display: flex;
		align-items: center;
		gap: 0.6rem;
	}

	.brand-title {
		font-weight: 800;
		font-size: 1.15rem;
		color: var(--text-main);
	}

	.tagline {
		font-size: 0.85rem;
		color: var(--text-muted);
		line-height: 1.4;
	}

	.footer-links {
		display: flex;
		gap: 1.5rem;
		font-size: 0.9rem;
		font-weight: 600;
	}

	.footer-links a {
		color: var(--text-muted);
		transition: color 0.15s ease;
	}

	.footer-links a:hover {
		color: var(--brand-primary);
	}

	.footer-bottom {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding-top: 1.25rem;
		font-size: 0.8rem;
		color: var(--text-subtle);
	}

	.privacy-badge {
		font-weight: 600;
		color: var(--text-muted);
	}

	.copyright {
		display: flex;
		gap: 0.75rem;
		align-items: center;
	}

	.version {
		font-family: var(--font-mono);
		color: var(--text-muted);
	}

	@media (max-width: 640px) {
		.footer-content {
			flex-direction: column;
			gap: 1.5rem;
		}

		.footer-links {
			flex-direction: column;
			gap: 0.75rem;
		}

		.footer-bottom {
			flex-direction: column;
			gap: 0.5rem;
			align-items: flex-start;
		}

		.copyright {
			flex-wrap: wrap;
		}
	}
</style>
