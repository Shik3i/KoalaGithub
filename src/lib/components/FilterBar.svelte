<script lang="ts">
	import type { VisualizerCategory, SortOption } from '$lib/types/visualizer.types';

	let {
		selectedCategory = $bindable<VisualizerCategory | 'all'>('all'),
		searchQuery = $bindable(''),
		selectedTag = $bindable<string | null>(null),
		selectedSort = $bindable<SortOption>('popularity'),
		totalCount = 0,
		availableTags = [] as string[]
	} = $props<{
		selectedCategory?: VisualizerCategory | 'all';
		searchQuery?: string;
		selectedTag?: string | null;
		selectedSort?: SortOption;
		totalCount?: number;
		availableTags?: string[];
	}>();

	const CATEGORIES: { key: VisualizerCategory | 'all'; label: string; icon: string }[] = [
		{ key: 'all', label: 'All', icon: '✨' },
		{ key: 'stats', label: 'Stats Cards', icon: '📊' },
		{ key: 'streak', label: 'Streak', icon: '🔥' },
		{ key: 'activity', label: 'Activity Graph', icon: '📈' },
		{ key: 'configured', label: 'GitHub Actions (Self-Hosted)', icon: '⚙️' }
	];
</script>

<div class="filter-bar-container">
	<!-- Category Chips (Single Row, No Wrap) -->
	<div class="category-chips" role="tablist" aria-label="Visualizer Categories">
		{#each CATEGORIES as cat}
			<button
				type="button"
				role="tab"
				aria-selected={selectedCategory === cat.key}
				class="category-chip"
				class:active={selectedCategory === cat.key}
				onclick={() => (selectedCategory = cat.key)}
			>
				<span class="icon">{cat.icon}</span>
				<span>{cat.label}</span>
			</button>
		{/each}
	</div>

	<!-- Controls Row: Search Query & Sort -->
	<div class="controls-row">
		<div class="search-input-wrapper">
			<svg class="search-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
				<circle cx="11" cy="11" r="8" />
				<line x1="21" y1="21" x2="16.65" y2="16.65" />
			</svg>
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Filter visualizers by keyword or tag..."
				class="filter-search-input"
				aria-label="Filter visualizers"
			/>
			{#if searchQuery}
				<button type="button" class="clear-search" onclick={() => (searchQuery = '')}>✕</button>
			{/if}
		</div>

		<div class="sort-wrapper">
			<label for="sort-select" class="sort-label">Sort:</label>
			<select id="sort-select" bind:value={selectedSort} class="sort-select">
				<option value="popularity">Popularity</option>
				<option value="most-voted">Most Voted ▲</option>
				<option value="alphabetical">Alphabetical (A-Z)</option>
				<option value="recently-added">Recently Added</option>
			</select>
		</div>
	</div>

	<!-- Active Tags / Available Tag Chips -->
	{#if availableTags.length > 0}
		<div class="tags-container">
			<span class="tags-label">Tags:</span>
			{#each availableTags.slice(0, 10) as tag}
				<button
					type="button"
					class="tag-chip"
					class:active={selectedTag === tag}
					onclick={() => (selectedTag = selectedTag === tag ? null : tag)}
				>
					#{tag}
				</button>
			{/each}
			{#if selectedTag}
				<button type="button" class="clear-tag" onclick={() => (selectedTag = null)}>Clear tag filter</button>
			{/if}
		</div>
	{/if}

	<div class="results-meta">
		Showing <strong>{totalCount}</strong> visualizer{totalCount === 1 ? '' : 's'}
	</div>
</div>

<style>
	.filter-bar-container {
		display: flex;
		flex-direction: column;
		gap: 0.85rem;
		margin: 1.5rem 0 1rem;
	}

	.category-chips {
		display: flex;
		flex-wrap: nowrap;
		gap: 0.5rem;
		overflow-x: auto;
		padding-bottom: 0.25rem;
		scrollbar-width: none;
		-ms-overflow-style: none;
	}

	.category-chips::-webkit-scrollbar {
		display: none;
	}

	.category-chip {
		display: flex;
		align-items: center;
		gap: 0.4rem;
		padding: 0.4rem 0.8rem;
		border-radius: var(--radius-full);
		background-color: var(--bg-card);
		border: 1px solid var(--border-color);
		color: var(--text-muted);
		font-size: 0.85rem;
		font-weight: 600;
		white-space: nowrap;
		flex-shrink: 0;
		transition: all 0.15s ease;
		box-shadow: var(--shadow-sm);
	}

	.category-chip:hover {
		border-color: var(--border-hover);
		color: var(--text-main);
	}

	.category-chip.active {
		background-color: var(--brand-primary);
		color: #ffffff;
		border-color: var(--brand-primary);
	}

	.controls-row {
		display: flex;
		gap: 1rem;
		align-items: center;
		justify-content: space-between;
	}

	.search-input-wrapper {
		position: relative;
		display: flex;
		align-items: center;
		flex: 1;
		background-color: var(--bg-card);
		border: 1px solid var(--border-color);
		border-radius: var(--radius-md);
	}

	.search-icon {
		position: absolute;
		left: 0.75rem;
		color: var(--text-subtle);
		pointer-events: none;
	}

	.filter-search-input {
		width: 100%;
		padding: 0.55rem 2rem 0.55rem 2.25rem;
		background: transparent;
		border: none;
		color: var(--text-main);
		font-size: 0.9rem;
		outline: none;
	}

	.clear-search {
		position: absolute;
		right: 0.75rem;
		color: var(--text-muted);
		font-size: 0.8rem;
	}

	.sort-wrapper {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		white-space: nowrap;
	}

	.sort-label {
		font-size: 0.85rem;
		font-weight: 600;
		color: var(--text-muted);
	}

	.sort-select {
		padding: 0.5rem 0.75rem;
		border-radius: var(--radius-md);
		border: 1px solid var(--border-color);
		background-color: var(--bg-card);
		color: var(--text-main);
		font-size: 0.875rem;
		font-weight: 600;
		outline: none;
	}

	.tags-container {
		display: flex;
		flex-wrap: nowrap;
		overflow-x: auto;
		align-items: center;
		gap: 0.4rem;
		font-size: 0.8rem;
		padding-bottom: 0.2rem;
		scrollbar-width: none;
	}

	.tags-container::-webkit-scrollbar {
		display: none;
	}

	.tag-chip {
		background-color: var(--bg-subtle);
		color: var(--text-muted);
		padding: 0.15rem 0.5rem;
		border-radius: var(--radius-sm);
		font-size: 0.775rem;
		font-weight: 500;
		border: 1px solid transparent;
		white-space: nowrap;
		flex-shrink: 0;
		transition: all 0.15s ease;
	}

	.tag-chip:hover, .tag-chip.active {
		border-color: var(--brand-primary);
		color: var(--brand-primary);
		background-color: var(--brand-light);
	}

	.clear-tag {
		color: var(--error-text);
		font-size: 0.775rem;
		font-weight: 600;
		text-decoration: underline;
		margin-left: 0.25rem;
		white-space: nowrap;
	}

	.results-meta {
		font-size: 0.85rem;
		color: var(--text-muted);
		margin-top: -0.25rem;
	}

	@media (max-width: 640px) {
		.controls-row {
			flex-direction: column;
			align-items: stretch;
		}

		.sort-wrapper {
			justify-content: flex-end;
		}
	}
</style>
