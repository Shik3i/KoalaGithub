<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import '../app.css';

	let { children } = $props();
</script>

<div class="app-layout">
	<a class="skip-link" href="#main-content">Skip to main content</a>
	<Header />
	<main id="main-content" class="main-content">
		{@render children()}
	</main>
	<Footer />
</div>

<style>
	.app-layout {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
	}

	.main-content {
		flex: 1;
	}

	.skip-link {
		position: fixed;
		top: 0.75rem;
		left: 0.75rem;
		z-index: 1000;
		padding: 0.65rem 0.9rem;
		border-radius: var(--radius-md);
		background: var(--brand-solid);
		color: var(--brand-on-solid);
		font-weight: 700;
		transform: translateY(-150%);
		transition: transform 0.15s ease;
	}

	.skip-link:focus {
		transform: translateY(0);
	}
</style>
