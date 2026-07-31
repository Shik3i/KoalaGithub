<script lang="ts">
	import VisualizerCard from './VisualizerCard.svelte';
	import type { Visualizer } from '$lib/types/visualizer.types';

	let { visualizers, username } = $props<{
		visualizers: Visualizer[];
		username: string;
	}>();
</script>

<div class="grid-container">
	{#if visualizers.length === 0}
		<div class="empty-state">
			<span class="empty-icon">🔍</span>
			<h3>No visualizers found</h3>
			<p>Try adjusting your search, category, tag, or workflow filter.</p>
		</div>
	{:else}
		<div class="grid">
			{#each visualizers as item (item.id)}
				<VisualizerCard visualizer={item} {username} />
			{/each}
		</div>
	{/if}
</div>

<style>
	.grid-container {
		width: 100%;
		margin: 1rem 0 3rem;
	}

	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(420px, 1fr));
		gap: 1.5rem;
		min-width: 0;
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 4rem 1.5rem;
		text-align: center;
		background-color: var(--bg-card);
		border: 1px solid var(--border-color);
		border-radius: var(--radius-lg);
		gap: 0.5rem;
	}

	.empty-icon {
		font-size: 2.5rem;
	}

	.empty-state h3 {
		font-size: 1.25rem;
		font-weight: 700;
		color: var(--text-main);
	}

	.empty-state p {
		font-size: 0.9rem;
		color: var(--text-muted);
		max-width: 440px;
	}

	@media (max-width: 640px) {
		.grid {
			grid-template-columns: 1fr;
		}
	}
</style>
