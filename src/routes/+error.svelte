<script lang="ts">
	import { page } from '$app/state';
</script>

<svelte:head>
	<title>{page.status === 404 ? 'Page not found' : 'Something went wrong'} – KoalaGitHub</title>
	<meta name="robots" content="noindex" />
</svelte:head>

<section class="error-page">
	<p class="status">{page.status}</p>
	<h1>{page.status === 404 ? 'This page wandered off.' : 'Something went wrong.'}</h1>
	<p>{page.error?.message || 'The requested page could not be loaded.'}</p>
	<a href="/">Back to the visualizer hub</a>
</section>

<style>
	.error-page {
		width: min(100% - 2rem, 680px);
		margin: 5rem auto;
		padding: 2rem;
		text-align: center;
		background: var(--bg-card);
		border: 1px solid var(--border-color);
		border-radius: var(--radius-lg);
	}
	.status {
		color: var(--brand-primary);
		font: 800 4rem/1 var(--font-mono);
	}
	h1 {
		margin-top: 0.75rem;
		font-size: 1.75rem;
	}
	p:not(.status) {
		margin: 0.75rem 0 1.5rem;
		color: var(--text-muted);
	}
	a {
		display: inline-flex;
		padding: 0.6rem 0.9rem;
		border-radius: var(--radius-md);
		background: var(--brand-solid);
		color: var(--brand-on-solid);
		color: white;
		font-weight: 700;
	}
</style>
