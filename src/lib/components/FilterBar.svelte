<script lang="ts">
	import type {
		VisualizerCategory,
		SortDirection,
		SortOption,
		WorkflowFilter
	} from '$lib/types/visualizer.types';

	let {
		selectedCategory = $bindable<VisualizerCategory | 'all'>('all'),
		searchQuery = $bindable(''),
		selectedTag = $bindable<string | null>(null),
		workflowFilter = $bindable<WorkflowFilter>('all'),
		selectedSort = $bindable<SortOption>('most-voted'),
		sortDirection = $bindable<SortDirection>('descending'),
		totalCount = 0,
		availableTags = [] as string[]
	} = $props<{
		selectedCategory?: VisualizerCategory | 'all';
		searchQuery?: string;
		selectedTag?: string | null;
		workflowFilter?: WorkflowFilter;
		selectedSort?: SortOption;
		sortDirection?: SortDirection;
		totalCount?: number;
		availableTags?: string[];
	}>();

	const CATEGORIES: { key: VisualizerCategory | 'all'; label: string; icon: string }[] = [
		{ key: 'all', label: 'All', icon: '✨' },
		{ key: 'stats', label: 'Stats Cards', icon: '📊' },
		{ key: 'streak', label: 'Streak', icon: '🔥' },
		{ key: 'activity', label: 'Activity Graph', icon: '📈' },
		{ key: 'configured', label: 'Setup required', icon: '⚙️' }
	];

	const SORT_DESCRIPTIONS: Record<SortOption, [string, string]> = {
		'most-voted': ['Fewest community votes first.', 'Most community votes first.'],
		'most-stars': ['Fewest repository stars first.', 'Most repository stars first.'],
		alphabetical: ['Visualizer names from A to Z.', 'Visualizer names from Z to A.'],
		'recently-added': ['Oldest registry entries first.', 'Newest registry entries first.']
	};

	const WORKFLOW_FILTERS: { key: WorkflowFilter; label: string }[] = [
		{ key: 'all', label: 'All' },
		{ key: 'workflow', label: 'Workflow' },
		{ key: 'instant', label: 'No workflow' }
	];
</script>

<div class="filter-bar-container">
	<!-- Category Chips (Single Row, No Wrap) -->
	<div class="category-chips" role="group" aria-label="Visualizer categories">
		{#each CATEGORIES as cat (cat.key)}
			<button
				type="button"
				aria-pressed={selectedCategory === cat.key}
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
			<svg
				class="search-icon"
				width="16"
				height="16"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
			>
				<circle cx="11" cy="11" r="8" />
				<line x1="21" y1="21" x2="16.65" y2="16.65" />
			</svg>
			<input
				type="search"
				bind:value={searchQuery}
				placeholder="Filter visualizers by keyword or tag..."
				class="filter-search-input"
				aria-label="Filter visualizers"
				autocomplete="off"
				spellcheck="false"
			/>
			{#if searchQuery}
				<button
					type="button"
					class="clear-search"
					onclick={() => (searchQuery = '')}
					aria-label="Clear visualizer search">✕</button
				>
			{/if}
		</div>

		<div class="sort-wrapper">
			<label for="sort-select" class="sort-label">Sort:</label>
			<select id="sort-select" bind:value={selectedSort} class="sort-select">
				<option value="most-voted">Votes</option>
				<option value="most-stars">Repository Stars</option>
				<option value="alphabetical">Name</option>
				<option value="recently-added">Recently Added</option>
			</select>
			<button
				type="button"
				class="direction-toggle"
				onclick={() => (sortDirection = sortDirection === 'ascending' ? 'descending' : 'ascending')}
				aria-label={`Sort ${sortDirection === 'ascending' ? 'ascending' : 'descending'}`}
				title={`Sort ${sortDirection === 'ascending' ? 'ascending' : 'descending'}; click to reverse`}
			>
				<svg
					width="18"
					height="18"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					aria-hidden="true"
				>
					{#if sortDirection === 'ascending'}
						<path d="M12 19V5" />
						<path d="m6 11 6-6 6 6" />
					{:else}
						<path d="M12 5v14" />
						<path d="m18 13-6 6-6-6" />
					{/if}
				</svg>
			</button>
		</div>
	</div>
	<fieldset class="workflow-filter">
		<legend>Setup:</legend>
		<div class="workflow-options">
			{#each WORKFLOW_FILTERS as option (option.key)}
				<button
					type="button"
					class:active={workflowFilter === option.key}
					aria-pressed={workflowFilter === option.key}
					onclick={() => (workflowFilter = option.key)}
				>
					{option.label}
				</button>
			{/each}
		</div>
	</fieldset>
	<p class="sort-help" role="status">
		{SORT_DESCRIPTIONS[selectedSort as SortOption][sortDirection === 'ascending' ? 0 : 1]}
	</p>

	<!-- Active Tags / Available Tag Chips -->
	{#if availableTags.length > 0}
		<div class="tags-container">
			<span class="tags-label">Tags:</span>
			{#each availableTags.slice(0, 10) as tag (tag)}
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
				<button type="button" class="clear-tag" onclick={() => (selectedTag = null)}
					>Clear tag filter</button
				>
			{/if}
		</div>
	{/if}

	<div class="results-meta" role="status" aria-live="polite">
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
		min-height: 2.5rem;
		transition:
			background-color 0.15s ease,
			border-color 0.15s ease,
			color 0.15s ease;
		box-shadow: var(--shadow-sm);
	}

	.category-chip:hover {
		border-color: var(--border-hover);
		color: var(--text-main);
	}

	.category-chip.active {
		background-color: var(--brand-solid);
		color: var(--brand-on-solid);
		border-color: var(--brand-solid);
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
		min-width: 2.5rem;
		min-height: 2.5rem;
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
		min-height: 2.5rem;
	}

	.direction-toggle {
		display: inline-grid;
		width: 2.5rem;
		height: 2.5rem;
		place-items: center;
		flex: 0 0 auto;
		border: 1px solid var(--border-color);
		border-radius: var(--radius-md);
		background: var(--bg-card);
		color: var(--text-main);
		transition:
			background-color 0.15s ease,
			border-color 0.15s ease,
			color 0.15s ease;
	}

	.direction-toggle:hover {
		background: var(--bg-card-hover);
		border-color: var(--border-hover);
		color: var(--brand-primary);
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

	.sort-help {
		margin-top: -0.35rem;
		color: var(--text-subtle);
		font-size: 0.75rem;
	}

	.workflow-filter {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		min-width: 0;
		border: 0;
		padding: 0;
	}

	.workflow-filter legend {
		float: left;
		font-size: 0.85rem;
		font-weight: 600;
		color: var(--text-muted);
	}

	.workflow-options {
		display: inline-flex;
		min-width: 0;
		padding: 0.2rem;
		border: 1px solid var(--border-color);
		border-radius: var(--radius-md);
		background: var(--bg-card);
		box-shadow: var(--shadow-sm);
	}

	.workflow-options button {
		padding: 0.35rem 0.7rem;
		border-radius: calc(var(--radius-md) - 0.2rem);
		color: var(--text-muted);
		font-size: 0.8rem;
		font-weight: 650;
		white-space: nowrap;
		min-height: 2.5rem;
	}

	.workflow-options button:hover {
		color: var(--text-main);
		background: var(--bg-card-hover);
	}

	.workflow-options button.active {
		color: var(--brand-on-solid);
		background: var(--brand-solid);
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
		min-height: 2rem;
		transition:
			background-color 0.15s ease,
			border-color 0.15s ease,
			color 0.15s ease;
	}

	.tag-chip:hover,
	.tag-chip.active {
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

		.workflow-filter {
			align-items: flex-start;
			flex-direction: column;
		}

		.workflow-options {
			width: 100%;
		}

		.workflow-options button {
			flex: 1;
		}
	}
</style>
