<script lang="ts">
	import { onMount } from 'svelte';
	import UserSearch from '$lib/components/UserSearch.svelte';
	import FilterBar from '$lib/components/FilterBar.svelte';
	import VisualizerGrid from '$lib/components/VisualizerGrid.svelte';
	import { VISUALIZERS as FALLBACK_VISUALIZERS } from '$lib/data/visualizers';
	import { DEFAULT_USERNAME, sanitizeUsername, updateUsernameUrl } from '$lib/utils/username';
	import { getOrCreateDeviceId } from '$lib/utils/device';
	import type {
		Visualizer,
		VisualizerCategory,
		SortDirection,
		SortOption
	} from '$lib/types/visualizer.types';

	let activeUsername = $state(DEFAULT_USERNAME);
	let selectedCategory = $state<VisualizerCategory | 'all'>('all');
	let searchQuery = $state('');
	let selectedTag = $state<string | null>(null);
	let selectedSort = $state<SortOption>('most-voted');
	let sortDirection = $state<SortDirection>('descending');
	let visualizerList = $state<Visualizer[]>(FALLBACK_VISUALIZERS);
	let usingFallback = $state(false);

	onMount(async () => {
		const urlParams = new URLSearchParams(window.location.search);
		const userParam = urlParams.get('user');
		if (userParam) {
			activeUsername = sanitizeUsername(userParam);
		}

		// Fetch live visualizers + vote data from Go backend API
		const deviceId = getOrCreateDeviceId();
		try {
			const headers: HeadersInit = deviceId ? { 'X-Device-ID': deviceId } : {};
			const res = await fetch('/api/visualizers', { headers });
			if (res.ok) {
				const data: unknown = await res.json();
				if (Array.isArray(data) && data.every(isVisualizer)) {
					visualizerList = data;
					return;
				}
			}
			usingFallback = true;
		} catch {
			usingFallback = true;
		}
	});

	function isVisualizer(value: unknown): value is Visualizer {
		if (!value || typeof value !== 'object') return false;
		const item = value as Record<string, unknown>;
		const themes = item.themes;
		const categories = new Set([
			'stats',
			'languages',
			'activity',
			'streak',
			'trophies',
			'headers',
			'configured',
			'badges',
			'quotes'
		]);
		const previewTypes = new Set(['image', 'iframe', 'setup-card', 'opt-in-counter']);
		const imageFormats = new Set(['svg', 'png', 'gif', 'jpg', 'auto']);
		return (
			typeof item.id === 'string' &&
			typeof item.name === 'string' &&
			typeof item.description === 'string' &&
			typeof item.category === 'string' &&
			categories.has(item.category) &&
			Array.isArray(item.tags) &&
			item.tags.every((tag) => typeof tag === 'string') &&
			typeof item.previewType === 'string' &&
			previewTypes.has(item.previewType) &&
			typeof item.imageFormat === 'string' &&
			imageFormats.has(item.imageFormat) &&
			typeof item.imageUrlTemplate === 'string' &&
			typeof item.markdownTemplate === 'string' &&
			typeof item.websiteUrl === 'string' &&
			typeof item.repositoryUrl === 'string' &&
			Array.isArray(themes) &&
			themes.every(
				(theme) =>
					theme &&
					typeof theme === 'object' &&
					typeof (theme as Record<string, unknown>).key === 'string' &&
					typeof (theme as Record<string, unknown>).label === 'string'
			) &&
			typeof item.defaultTheme === 'string' &&
			typeof item.requiresUsername === 'boolean' &&
			typeof item.requiresExternalSetup === 'boolean' &&
			typeof item.addedAt === 'string' &&
			/^\d{4}-\d{2}-\d{2}$/.test(item.addedAt) &&
			typeof item.enabled === 'boolean' &&
			typeof item.voteCount === 'number' &&
			Number.isFinite(item.voteCount) &&
			item.voteCount >= 0 &&
			typeof item.userVoted === 'boolean' &&
			(item.githubStars === undefined ||
				(typeof item.githubStars === 'number' &&
					Number.isFinite(item.githubStars) &&
					item.githubStars >= 0))
		);
	}

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
				let comparison = 0;
				if (selectedSort === 'most-voted') {
					comparison = a.voteCount - b.voteCount;
				} else if (selectedSort === 'most-stars') {
					comparison = (a.githubStars || 0) - (b.githubStars || 0);
				} else if (selectedSort === 'alphabetical') {
					comparison = a.name.localeCompare(b.name);
				} else if (selectedSort === 'recently-added') {
					comparison = new Date(a.addedAt).getTime() - new Date(b.addedAt).getTime();
				}
				const directed = sortDirection === 'ascending' ? comparison : -comparison;
				return directed || a.name.localeCompare(b.name);
			});
	});
</script>

<svelte:head>
	<title>KoalaGitHub – Central GitHub Profile Visualizer Directory</title>
	<link rel="canonical" href="https://github.koalastuff.net/" />
	<meta property="og:url" content="https://github.koalastuff.net/" />
</svelte:head>

<section class="hero-section">
	<div class="container hero-container">
		<img
			class="hero-logo"
			src="/assets/brand/koalagithub-logo-256.png"
			srcset="
				/assets/brand/koalagithub-logo-128.png 128w,
				/assets/brand/koalagithub-logo-256.png 256w,
				/assets/brand/koalagithub-logo-512.png 512w
			"
			sizes="(max-width: 640px) 132px, 180px"
			width="180"
			height="180"
			alt=""
			fetchpriority="high"
			decoding="async"
		/>
		<h1 class="hero-title">Compare GitHub Profile Visualizers in One Place</h1>
		<p class="hero-subtitle">
			Enter any GitHub username below to instantly preview, customize themes, upvote favorites, and
			copy ready-to-use README Markdown code.
		</p>

		<UserSearch bind:username={activeUsername} onUpdate={handleUsernameUpdate} />
	</div>
</section>

<section class="hub-section">
	<div class="container">
		{#if usingFallback}
			<div class="service-notice" role="status">
				Live votes and star counts are temporarily unavailable. Visualizer previews and Markdown
				generation still work.
			</div>
		{/if}

		<FilterBar
			bind:selectedCategory
			bind:searchQuery
			bind:selectedTag
			bind:selectedSort
			bind:sortDirection
			totalCount={filteredVisualizers.length}
			{availableTags}
		/>

		{#if selectedCategory === 'configured'}
			<div class="category-info-banner">
				<span class="banner-icon">⚙️</span>
				<div class="banner-text">
					<h2>Repository-generated visualizers</h2>
					<p>
						These visualizers require an automated <strong>GitHub Actions workflow</strong> in your
						profile repository. Depending on the project, the generated image is committed to the
						<code>output</code>
						branch or the default branch. A preview appears only after that file exists for
						<strong>{activeUsername}</strong>. Read the
						<a href="/guide" class="banner-guide-link">Complete Setup Guide 📖</a> or copy the 1-click
						workflow shown on each card.
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

	.hero-logo {
		width: clamp(8.25rem, 18vw, 11.25rem);
		height: auto;
		filter: drop-shadow(0 10px 18px rgba(15, 23, 42, 0.16));
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

	.service-notice {
		margin: 0 0 1rem;
		padding: 0.75rem 1rem;
		border: 1px solid var(--warning-text);
		border-radius: var(--radius-md);
		background: var(--warning-bg);
		color: var(--warning-text);
		font-size: 0.875rem;
	}

	.category-info-banner {
		display: flex;
		align-items: flex-start;
		gap: 1rem;
		background: linear-gradient(
			135deg,
			color-mix(in srgb, var(--brand-primary) 8%, transparent) 0%,
			color-mix(in srgb, var(--koala-accent) 4%, transparent) 100%
		);
		border: 1px solid color-mix(in srgb, var(--brand-primary) 25%, transparent);
		border-radius: var(--radius-md);
		padding: 1.1rem 1.3rem;
		margin-bottom: 1.5rem;
	}

	.category-info-banner .banner-icon {
		font-size: 1.5rem;
		line-height: 1;
		margin-top: 0.15rem;
	}

	.category-info-banner .banner-text h2 {
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
