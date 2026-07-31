<script lang="ts">
	import type { Visualizer } from '$lib/types/visualizer.types';
	import { renderTemplate } from '$lib/utils/markdown';
	import { getOrCreateDeviceId } from '$lib/utils/device';

	let { visualizer = $bindable(), username } = $props<{
		visualizer: Visualizer;
		username: string;
	}>();

	let selectedTheme = $state('');
	$effect(() => {
		if (!selectedTheme) {
			selectedTheme = visualizer.defaultTheme;
		}
	});
	let imageState = $state<'loading' | 'loaded' | 'error'>('loading');
	let reloadKey = $state(0);
	let optInLoaded = $state(false);
	let copied = $state(false);
	let copiedYaml = $state(false);
	let voting = $state(false);

	// Derived image URL & Markdown snippet
	let imageUrl = $derived(renderTemplate(visualizer.imageUrlTemplate, username, selectedTheme));
	let markdownCode = $derived(renderTemplate(visualizer.markdownTemplate, username, selectedTheme));

	// Reset image state when theme or username or reloadKey changes
	$effect(() => {
		// track dependencies
		imageUrl;
		reloadKey;
		imageState = 'loading';
	});

	function handleCopy() {
		if (typeof navigator !== 'undefined' && navigator.clipboard) {
			navigator.clipboard.writeText(markdownCode);
			copied = true;
			setTimeout(() => {
				copied = false;
			}, 2000);
		}
	}

	function handleCopyYaml() {
		if (visualizer.actionWorkflowYaml && typeof navigator !== 'undefined' && navigator.clipboard) {
			navigator.clipboard.writeText(visualizer.actionWorkflowYaml);
			copiedYaml = true;
			setTimeout(() => {
				copiedYaml = false;
			}, 2000);
		}
	}

	function handleRetry() {
		imageState = 'loading';
		reloadKey++;
	}

	function formatStars(count?: number): string {
		if (!count || count === 0) return '';
		if (count >= 1000) {
			return (count / 1000).toFixed(1).replace(/\.0$/, '') + 'k';
		}
		return count.toString();
	}

	function handleSelectTheme(themeKey: string) {
		selectedTheme = themeKey;
	}

	async function handleVote() {
		if (voting) return;
		voting = true;
		const deviceId = getOrCreateDeviceId();

		try {
			const res = await fetch(`/api/visualizers/${visualizer.id}/vote`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ deviceId })
			});
			if (res.ok) {
				const data = await res.json();
				visualizer.userVoted = data.voted;
				visualizer.voteCount = data.voteCount;
			}
		} catch (e) {
			console.error('Failed to register vote:', e);
		} finally {
			voting = false;
		}
	}
</script>

<article class="card">
	<!-- Card Header -->
	<div class="card-header">
		<div class="header-main">
			<h3 class="card-title">{visualizer.name}</h3>
			<div class="badge-row">
				<span class="category-badge">{visualizer.category}</span>
				<button
					type="button"
					class="reddit-upvote-btn"
					class:voted={visualizer.userVoted}
					onclick={handleVote}
					disabled={voting}
					title={visualizer.userVoted ? 'Remove your upvote' : 'Upvote this visualizer'}
				>
					<svg class="upvote-arrow" width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
						<path d="M12 4L3 15h6v5h6v-5h6L12 4z" />
					</svg>
					<span class="vote-count">{visualizer.voteCount}</span>
				</button>
			</div>
		</div>
		<div class="header-links">
			<a
				href={visualizer.repositoryUrl}
				target="_blank"
				rel="noreferrer"
				class="repo-link"
				title="View Open-Source GitHub Repository"
			>
				<svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor">
					<path
						d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"
					/>
				</svg>
				<span>Repo</span>
				{#if visualizer.githubStars && visualizer.githubStars > 0}
					<span class="star-count">⭐ {formatStars(visualizer.githubStars)}</span>
				{/if}
			</a>
			<a
				href={visualizer.websiteUrl}
				target="_blank"
				rel="noreferrer"
				class="link-icon"
				title="Visit Original Site"
			>
				<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
					<path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
					<polyline points="15 3 21 3 21 9" />
					<line x1="10" y1="14" x2="21" y2="3" />
				</svg>
			</a>
		</div>
	</div>

	<!-- Description -->
	<p class="card-desc">{visualizer.description}</p>

	<!-- Subtle Tags -->
	{#if visualizer.tags && visualizer.tags.length > 0}
		<div class="card-tags">
			{#each visualizer.tags as tag}
				<span class="card-tag-chip">#{tag}</span>
			{/each}
		</div>
	{/if}

	<!-- Theme Switcher Row -->
	{#if visualizer.themes.length > 1}
		<div class="themes-row">
			<span class="themes-label">Theme:</span>
			<div class="theme-buttons">
				{#each visualizer.themes as t}
					<button
						type="button"
						class="theme-chip"
						class:active={selectedTheme === t.key}
						onclick={() => handleSelectTheme(t.key)}
					>
						{t.label}
					</button>
				{/each}
			</div>
		</div>
	{/if}

	<!-- Preview Container -->
	<div class="preview-box" style="min-height: {visualizer.estimatedHeight || 170}px;">
		{#if visualizer.previewType === 'setup-card'}
			<!-- Configured Setup Card -->
			<div class="setup-notice">
				<div class="setup-header">
					<span class="setup-icon">⚙️</span>
					<h4>External Setup Required</h4>
				</div>
				<p class="setup-text">{visualizer.setupExplanation}</p>
				{#if visualizer.privacyNotice}
					<p class="setup-privacy">{visualizer.privacyNotice}</p>
				{/if}
				<a href={visualizer.websiteUrl} target="_blank" rel="noreferrer" class="setup-btn">
					View Setup Guide & Documentation ↗
				</a>
			</div>
		{:else if visualizer.previewType === 'opt-in-counter' && !optInLoaded}
			<!-- Opt-in Visitor Counter Placeholder -->
			<div class="opt-in-box">
				<span class="counter-icon">👁️</span>
				<p class="counter-notice">{visualizer.privacyNotice}</p>
				<button type="button" class="opt-in-btn" onclick={() => (optInLoaded = true)}>
					Load Live Counter Preview ({username})
				</button>
			</div>
		{:else}
			<!-- Live Image Preview Pipeline -->
			{#if imageState === 'loading'}
				<div class="skeleton-overlay" style="height: {visualizer.estimatedHeight || 170}px;">
					<div class="spinner"></div>
					<span>Loading visualizer preview...</span>
				</div>
			{/if}

			{#if imageState === 'error'}
				<div class="error-box" style="min-height: {visualizer.estimatedHeight || 170}px;">
					<span class="error-icon">⚙️</span>
					{#if visualizer.requiresExternalSetup}
						<p class="error-text">GitHub Action workflow required to generate SVG output for <strong>{username}</strong>.</p>
						{#if visualizer.setupExplanation}
							<p class="setup-privacy">{visualizer.setupExplanation}</p>
						{/if}
						<div class="error-actions">
							<a href={visualizer.websiteUrl} target="_blank" rel="noreferrer" class="retry-btn">View Setup Guide ↗</a>
							<button type="button" class="open-direct-btn" onclick={handleRetry}>Retry</button>
						</div>
					{:else}
						<p class="error-text">Failed to load preview image from third-party provider.</p>
						<div class="error-actions">
							<button type="button" class="retry-btn" onclick={handleRetry}>Retry</button>
							<a href={imageUrl} target="_blank" rel="noreferrer" class="open-direct-btn">Open Image Direct ↗</a>
						</div>
					{/if}
				</div>
			{/if}

			{#key `${imageUrl}-${reloadKey}`}
				<img
					src={imageUrl}
					alt="{visualizer.name} preview for {username}"
					loading="lazy"
					referrerpolicy="no-referrer"
					class="preview-img"
					class:hidden={imageState === 'error'}
					onload={() => (imageState = 'loaded')}
					onerror={() => (imageState = 'error')}
				/>
			{/key}
		{/if}
	</div>

	<!-- Code Preview & Markdown Copy Section -->
	<div class="code-section">
		<div class="code-header">
			<span class="code-label">README Markdown Code:</span>
			<div class="code-actions">
				{#if visualizer.actionWorkflowYaml}
					<button type="button" class="copy-btn copy-yaml-btn" onclick={handleCopyYaml}>
						{#if copiedYaml}
							<span class="copied-indicator">✓ YAML Copied!</span>
						{:else}
							<span>📋 Copy Workflow (.yml)</span>
						{/if}
					</button>
				{/if}
				<button type="button" class="copy-btn" onclick={handleCopy}>
					{#if copied}
						<span class="copied-indicator">✓ Copied!</span>
					{:else}
						<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
							<rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
							<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
						</svg>
						<span>Copy Markdown</span>
					{/if}
				</button>
			</div>
		</div>
		<pre class="code-block"><code>{markdownCode}</code></pre>

		{#if visualizer.actionWorkflowYaml}
			<details class="workflow-details">
				<summary class="workflow-summary">⚙️ View GitHub Action Workflow (.github/workflows/*.yml)</summary>
				<pre class="code-block yaml-block"><code>{visualizer.actionWorkflowYaml}</code></pre>
			</details>
		{/if}
	</div>
</article>

<style>
	.card {
		background-color: var(--bg-card);
		border: 1px solid var(--border-color);
		border-radius: var(--radius-lg);
		padding: 1.25rem;
		display: flex;
		flex-direction: column;
		gap: 0.85rem;
		box-shadow: var(--shadow-sm);
		transition: border-color 0.2s ease, box-shadow 0.2s ease;
	}

	.card:hover {
		border-color: var(--border-hover);
		box-shadow: var(--shadow-md);
	}

	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 0.5rem;
	}

	.header-main {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.card-title {
		font-size: 1.05rem;
		font-weight: 700;
		color: var(--text-main);
		line-height: 1.3;
	}

	.badge-row {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.category-badge {
		align-self: flex-start;
		font-size: 0.725rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.04em;
		background-color: var(--brand-light);
		color: var(--brand-text);
		padding: 0.15rem 0.5rem;
		border-radius: var(--radius-sm);
	}

	.reddit-upvote-btn {
		display: inline-flex;
		align-items: center;
		gap: 0.35rem;
		font-size: 0.8rem;
		font-weight: 700;
		padding: 0.25rem 0.65rem;
		border-radius: var(--radius-full);
		background-color: var(--bg-subtle);
		color: var(--text-muted);
		border: 1px solid var(--border-color);
		cursor: pointer;
		user-select: none;
		transition: all 0.15s cubic-bezier(0.4, 0, 0.2, 1);
		box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
	}

	.reddit-upvote-btn .upvote-arrow {
		transition: transform 0.15s ease, color 0.15s ease;
	}

	.reddit-upvote-btn:hover {
		border-color: #ff4500;
		color: #ff4500;
		background-color: rgba(255, 69, 0, 0.08);
	}

	.reddit-upvote-btn:hover .upvote-arrow {
		transform: translateY(-2px);
	}

	.reddit-upvote-btn.voted {
		background: linear-gradient(135deg, #ff4500 0%, #ff5722 100%);
		color: #ffffff;
		border-color: #ff4500;
		box-shadow: 0 2px 8px rgba(255, 69, 0, 0.35);
	}

	.reddit-upvote-btn.voted .upvote-arrow {
		color: #ffffff;
		transform: translateY(-1px);
	}

	.reddit-upvote-btn:active {
		transform: scale(0.94);
	}

	.header-links {
		display: flex;
		align-items: center;
		gap: 0.4rem;
	}

	.repo-link {
		display: flex;
		align-items: center;
		gap: 0.3rem;
		font-size: 0.775rem;
		font-weight: 700;
		color: var(--text-main);
		background-color: var(--bg-subtle);
		padding: 4px 8px;
		border-radius: var(--radius-md);
		border: 1px solid var(--border-color);
		transition: all 0.15s ease;
	}

	.star-count {
		font-size: 0.725rem;
		font-weight: 600;
		color: #eab308;
		background-color: rgba(234, 179, 8, 0.12);
		padding: 0.05rem 0.35rem;
		border-radius: var(--radius-sm);
		margin-left: 0.15rem;
	}

	.repo-link:hover {
		background-color: var(--brand-light);
		color: var(--brand-text);
		border-color: var(--brand-primary);
	}

	.link-icon {
		color: var(--text-muted);
		padding: 4px;
		border-radius: var(--radius-sm);
		transition: color 0.15s ease, background-color 0.15s ease;
	}

	.link-icon:hover {
		color: var(--text-main);
		background-color: var(--bg-subtle);
	}

	.card-desc {
		font-size: 0.875rem;
		color: var(--text-muted);
		line-height: 1.45;
	}

	.card-tags {
		display: flex;
		flex-wrap: wrap;
		gap: 0.35rem;
		margin-top: -0.25rem;
		margin-bottom: 0.2rem;
	}

	.card-tag-chip {
		font-size: 0.725rem;
		font-weight: 500;
		color: var(--text-subtle);
		background-color: rgba(255, 255, 255, 0.04);
		border: 1px solid var(--border-color);
		padding: 0.1rem 0.45rem;
		border-radius: var(--radius-sm);
		letter-spacing: -0.01em;
		opacity: 0.85;
		transition: opacity 0.15s ease, color 0.15s ease, border-color 0.15s ease;
	}

	.card-tag-chip:hover {
		opacity: 1;
		color: var(--text-main);
		border-color: var(--brand-primary);
	}

	.themes-row {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.8rem;
		flex-wrap: wrap;
	}

	.themes-label {
		font-weight: 600;
		color: var(--text-subtle);
	}

	.theme-buttons {
		display: flex;
		flex-wrap: wrap;
		gap: 0.35rem;
	}

	.theme-chip {
		font-size: 0.75rem;
		font-weight: 600;
		padding: 0.2rem 0.55rem;
		border-radius: var(--radius-sm);
		background-color: var(--bg-subtle);
		color: var(--text-muted);
		border: 1px solid var(--border-color);
		transition: all 0.15s ease;
	}

	.theme-chip:hover, .theme-chip.active {
		background-color: var(--brand-primary);
		color: #ffffff;
		border-color: var(--brand-primary);
	}

	.preview-box {
		position: relative;
		display: flex;
		align-items: center;
		justify-content: center;
		background-color: var(--bg-subtle);
		border-radius: var(--radius-md);
		border: 1px solid var(--border-color);
		overflow: hidden;
		padding: 0.75rem;
	}

	.preview-img {
		max-width: 100%;
		height: auto;
		display: block;
		object-fit: contain;
		margin: 0 auto;
		border-radius: var(--radius-sm);
	}

	.preview-img.hidden {
		display: none;
	}

	.skeleton-overlay {
		position: absolute;
		inset: 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		background-color: var(--bg-subtle);
		color: var(--text-muted);
		font-size: 0.825rem;
		z-index: 2;
	}

	.spinner {
		width: 22px;
		height: 22px;
		border: 2.5px solid var(--border-color);
		border-top-color: var(--brand-primary);
		border-radius: 50%;
		animation: spin 0.8s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.error-box {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		padding: 1rem;
		text-align: center;
		background-color: var(--error-bg);
		color: var(--error-text);
		border-radius: var(--radius-md);
		width: 100%;
	}

	.error-icon {
		font-size: 1.5rem;
	}

	.error-text {
		font-size: 0.85rem;
		font-weight: 500;
	}

	.error-actions {
		display: flex;
		gap: 0.5rem;
		margin-top: 0.25rem;
	}

	.retry-btn, .open-direct-btn {
		font-size: 0.775rem;
		font-weight: 600;
		padding: 0.3rem 0.65rem;
		border-radius: var(--radius-sm);
	}

	.retry-btn {
		background-color: var(--error-text);
		color: #ffffff;
	}

	.open-direct-btn {
		background-color: transparent;
		color: var(--error-text);
		border: 1px solid var(--error-border);
	}

	.setup-notice, .opt-in-box {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		text-align: center;
		gap: 0.5rem;
		padding: 1rem;
		width: 100%;
	}

	.setup-header {
		display: flex;
		align-items: center;
		gap: 0.4rem;
	}

	.setup-header h4 {
		font-size: 0.95rem;
		font-weight: 700;
		color: var(--text-main);
	}

	.setup-text {
		font-size: 0.825rem;
		color: var(--text-muted);
		line-height: 1.4;
		max-width: 500px;
	}

	.setup-privacy, .counter-notice {
		font-size: 0.775rem;
		color: var(--text-subtle);
		font-style: italic;
	}

	.setup-btn, .opt-in-btn {
		margin-top: 0.4rem;
		font-size: 0.825rem;
		font-weight: 700;
		padding: 0.45rem 0.85rem;
		border-radius: var(--radius-md);
		background-color: var(--brand-primary);
		color: #ffffff;
		transition: background-color 0.15s ease;
	}

	.setup-btn:hover, .opt-in-btn:hover {
		background-color: var(--brand-hover);
		color: #ffffff;
	}

	.code-section {
		display: flex;
		flex-direction: column;
		gap: 0.35rem;
		background-color: var(--bg-code);
		border-radius: var(--radius-md);
		padding: 0.75rem;
	}

	.code-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		flex-wrap: wrap;
		gap: 0.5rem;
	}

	.code-actions {
		display: flex;
		align-items: center;
		gap: 0.4rem;
	}

	.code-label {
		font-size: 0.75rem;
		font-weight: 600;
		color: #94a3b8;
		text-transform: uppercase;
		letter-spacing: 0.03em;
	}

	.copy-btn {
		display: flex;
		align-items: center;
		gap: 0.35rem;
		font-size: 0.725rem;
		font-weight: 600;
		color: var(--text-muted);
		background-color: rgba(255, 255, 255, 0.05);
		border: 1px solid var(--border-color);
		padding: 0.25rem 0.6rem;
		border-radius: var(--radius-sm);
		cursor: pointer;
		transition: all 0.15s ease;
	}

	.copy-yaml-btn {
		background-color: rgba(99, 102, 241, 0.12);
		color: #818cf8;
		border-color: rgba(99, 102, 241, 0.3);
	}

	.copy-yaml-btn:hover {
		background-color: rgba(99, 102, 241, 0.22);
		color: #a5b4fc;
	}

	.workflow-details {
		margin-top: 0.4rem;
		border-top: 1px solid rgba(255, 255, 255, 0.06);
		padding-top: 0.4rem;
	}

	.workflow-summary {
		font-size: 0.75rem;
		font-weight: 600;
		color: #94a3b8;
		cursor: pointer;
		user-select: none;
		transition: color 0.15s ease;
	}

	.workflow-summary:hover {
		color: #e2e8f0;
	}

	.yaml-block {
		max-height: 180px;
		overflow-y: auto;
		margin-top: 0.4rem;
		white-space: pre;
	}

	.copy-btn:hover {
		background-color: rgba(255, 255, 255, 0.15);
		color: #ffffff;
	}

	.copied-indicator {
		color: #4ade80;
		font-weight: 700;
	}

	.code-block {
		font-family: var(--font-mono);
		font-size: 0.8rem;
		color: var(--text-code);
		overflow-x: auto;
		white-space: pre-wrap;
		word-break: break-all;
		margin: 0;
	}
</style>
