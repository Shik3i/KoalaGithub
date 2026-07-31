<script lang="ts">
	import { onMount } from 'svelte';
	import UserSearch from '$lib/components/UserSearch.svelte';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import VisualizerGrid from '$lib/components/VisualizerGrid.svelte';
	import { VISUALIZERS as FALLBACK_VISUALIZERS } from '$lib/data/visualizers';
	import { DEFAULT_USERNAME, sanitizeUsername, updateUsernameUrl } from '$lib/utils/username';
	import { getOrCreateDeviceId } from '$lib/utils/device';
	import type { Visualizer, VisualizerCategory, SortOption } from '$lib/types/visualizer.types';

	let activeUsername = $state(DEFAULT_USERNAME);
	let selectedCategory = $state<VisualizerCategory | 'all'>('all');
	let searchQuery = $state('');
	let selectedTag = $state<string | null>(null);
	let selectedSort = $state<SortOption>('popularity');
	let visualizerList = $state<Visualizer[]>(FALLBACK_VISUALIZERS);

	onMount(async () => {
		const urlParams = new URLSearchParams(window.location.search);
		const userParam = urlParams.get('user');
		if (userParam) {
			activeUsername = sanitizeUsername(userParam);
		}

		// Fetch live visualizers + vote data from Go backend API
		const deviceId = getOrCreateDeviceId();
		try {
			const res = await fetch(`/api/visualizers?deviceId=${encodeURIComponent(deviceId)}`);
			if (res.ok) {
				const data = await res.json();
				if (Array.isArray(data) && data.length > 0) {
					visualizerList = data;
				}
			}
		} catch (e) {
			console.log('Serving offline static visualizers list');
		}
	});

	function handleUsernameUpdate(newUsername: string) {
		activeUsername = sanitizeUsername(newUsername);
		updateUsernameUrl(activeUsername);
	}

	// Compute list of unique tags across enabled visualizers
	let availableTags = $derived.by(() => {
		const tagSet = new Set<string>();
		for (const v of visualizerList) {
			if (v.enabled) {
				for (const tag of v.tags) {
					tagSet.add(tag);
				}
			}
		}
		return Array.from(tagSet).sort();
	});

	// Filter & sort visualizers
	let filteredVisualizers = $derived.by(() => {
		return visualizerList
			.filter((v) => {
				if (!v.enabled) return false;

				// Category filter
				if (selectedCategory !== 'all' && v.category !== selectedCategory) {
					return false;
				}

				// Tag filter
				if (selectedTag && !v.tags.includes(selectedTag)) {
					return false;
				}

				// Search query filter
				if (searchQuery.trim()) {
					const q = searchQuery.toLowerCase().trim();
					const matchName = v.name.toLowerCase().includes(q);
					const matchDesc = v.description.toLowerCase().includes(q);
					const matchTag = v.tags.some((t) => t.toLowerCase().includes(q));
					if (!matchName && !matchDesc && !matchTag) {
						return false;
					}
				}

				return true;
			})
			.sort((a, b) => {
				if (selectedSort === 'most-voted') {
					return b.voteCount - a.voteCount;
				} else if (selectedSort === 'most-stars') {
					return (b.githubStars || 0) - (a.githubStars || 0);
				} else if (selectedSort === 'alphabetical') {
					return a.name.localeCompare(b.name);
				} else if (selectedSort === 'recently-added') {
					return new Date(b.addedAt).getTime() - new Date(a.addedAt).getTime();
				} else {
					// popularity
					return a.popularityRank - b.popularityRank;
				}
			});
	});
</script>

<svelte:head>
	<title>KoalaGitHub – Central GitHub Profile Visualizer Directory</title>
</svelte:head>

<section class="hero-section">
	<div class="container hero-container">
		<h1 class="hero-title">
			Compare All GitHub Profile Visualizers in One Place
		</h1>
		<p class="hero-subtitle">
			Enter any GitHub username below to instantly preview, customize themes, upvote favorites, and copy ready-to-use README Markdown code.
		</p>

		<UserSearch bind:username={activeUsername} onUpdate={handleUsernameUpdate} />
	</div>
</section>

<section class="hub-section">
	<div class="container">
		<FilterBar
			bind:selectedCategory
			bind:searchQuery
			bind:selectedTag
			bind:selectedSort
			totalCount={filteredVisualizers.length}
			{availableTags}
		/>

		{#if selectedCategory === 'configured'}
			<div class="category-info-banner">
				<span class="banner-icon">⚙️</span>
				<div class="banner-text">
					<h3>Self-Hosted GitHub Action Visualizers</h3>
					<p>
						These visualizers require running an automated <strong>GitHub Action workflow</strong> in your own profile repository.
						They build once a day directly into your repo's output branch. Below, live output is demonstrated for <strong>{activeUsername}</strong> as an example.
						Need help? Read our <a href="/guide" class="banner-guide-link">Complete Setup Guide 📖</a> or copy the 1-click workflow file on each card!
					</p>
				</div>
			</div>
		{/if}

		<VisualizerGrid visualizers={filteredVisualizers} username={activeUsername} />
	</div>
</section>

<style>
	.hero-section {
		background: linear-gradient(180deg, var(--bg-card) 0%, var(--bg-app) 100%);
		border-bottom: 1px solid var(--border-color);
		padding: 3rem 0 2.5rem;
		text-align: center;
	}

	.hero-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.75rem;
	}

	.hero-title {
		font-size: 2.25rem;
		font-weight: 800;
		letter-spacing: -0.03em;
		color: var(--text-main);
		line-height: 1.25;
		max-width: 800px;
	}

	.hero-subtitle {
		font-size: 1.05rem;
		color: var(--text-muted);
		max-width: 680px;
		margin-bottom: 0.75rem;
		line-height: 1.5;
	}

	.hub-section {
		padding: 1rem 0 3rem;
	}

	.category-info-banner {
		display: flex;
		align-items: flex-start;
		gap: 1rem;
		background: linear-gradient(135deg, rgba(99, 102, 241, 0.08) 0%, rgba(139, 92, 246, 0.04) 100%);
		border: 1px solid rgba(99, 102, 241, 0.25);
		border-radius: var(--radius-md);
		padding: 1.1rem 1.3rem;
		margin-bottom: 1.5rem;
	}

	.category-info-banner .banner-icon {
		font-size: 1.5rem;
		line-height: 1;
		margin-top: 0.15rem;
	}

	.category-info-banner .banner-text h3 {
		font-size: 1.05rem;
		font-weight: 700;
		color: var(--text-main);
		margin-bottom: 0.35rem;
	}

	.category-info-banner .banner-text p {
		font-size: 0.9rem;
		color: var(--text-muted);
		line-height: 1.5;
	}

	@media (max-width: 640px) {
		.hero-title {
			font-size: 1.65rem;
		}

		.hero-subtitle {
			font-size: 0.95rem;
		}
	}
</style>
