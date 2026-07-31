<script lang="ts">
	import type { Visualizer } from '$lib/types/visualizer.types';
	import { renderTemplate } from '$lib/utils/markdown';
	import { getOrCreateDeviceId } from '$lib/utils/device';

	let { visualizer = $bindable(), username } = $props<{
		visualizer: Visualizer;
		username: string;
	}>();

	let selectedTheme = $state('');
	let previewImage = $state<HTMLImageElement>();
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
	let actionStatus = $state('');
	let actionFailed = $state(false);

	// Derived image URL & Markdown snippet
	let imageUrl = $derived(renderTemplate(visualizer.imageUrlTemplate, username, selectedTheme));
	let markdownCode = $derived(renderTemplate(visualizer.markdownTemplate, username, selectedTheme));
	let isReferenceProfile = $derived(username.toLowerCase() === 'shik3i');

	// Reset image state when theme or username or reloadKey changes
	$effect(() => {
		// track dependencies
		void imageUrl;
		void reloadKey;
		imageState = 'loading';
	});

	$effect(() => {
		const image = previewImage;
		if (!image) return;

		const handleLoad = () => (imageState = 'loaded');
		const handleError = () => (imageState = 'error');
		image.addEventListener('load', handleLoad);
		image.addEventListener('error', handleError);

		if (image.complete) {
			imageState = image.naturalWidth > 0 ? 'loaded' : 'error';
		}

		return () => {
			image.removeEventListener('load', handleLoad);
			image.removeEventListener('error', handleError);
		};
	});

	async function copyText(value: string, kind: 'markdown' | 'workflow') {
		try {
			if (typeof navigator === 'undefined' || !navigator.clipboard) {
				throw new Error('Clipboard unavailable');
			}
			await navigator.clipboard.writeText(value);
			actionFailed = false;
			actionStatus = `${kind === 'workflow' ? 'Workflow' : 'Markdown'} copied to clipboard.`;
			if (kind === 'workflow') copiedYaml = true;
			else copied = true;
			setTimeout(() => {
				copied = false;
				copiedYaml = false;
				actionStatus = '';
			}, 2000);
		} catch {
			actionFailed = true;
			actionStatus = 'Clipboard access failed. Select and copy the code manually.';
		}
	}

	function handleCopy() {
		void copyText(markdownCode, 'markdown');
	}

	function handleCopyYaml() {
		if (visualizer.actionWorkflowYaml) {
			void copyText(visualizer.actionWorkflowYaml, 'workflow');
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
		if (!deviceId) {
			actionFailed = true;
			actionStatus = 'Voting requires functional browser storage and secure random values.';
			voting = false;
			return;
		}

		try {
			const res = await fetch(`/api/visualizers/${visualizer.id}/vote`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ deviceId })
			});
			if (!res.ok) throw new Error(`HTTP ${res.status}`);
			const data: unknown = await res.json();
			if (
				!data ||
				typeof data !== 'object' ||
				typeof (data as { voted?: unknown }).voted !== 'boolean' ||
				typeof (data as { voteCount?: unknown }).voteCount !== 'number'
			) {
				throw new Error('Invalid response');
			}
			visualizer.userVoted = (data as { voted: boolean }).voted;
			visualizer.voteCount = (data as { voteCount: number }).voteCount;
			actionFailed = false;
			actionStatus = visualizer.userVoted ? 'Vote added.' : 'Vote removed.';
		} catch {
			actionFailed = true;
			actionStatus = 'Vote could not be saved. Please try again.';
		} finally {
			voting = false;
		}
	}
</script>

<article class="card">
	<!-- Card Header -->
	<div class="card-header">
		<div class="header-main">
			<h2 class="card-title">{visualizer.name}</h2>
			<div class="badge-row">
				<span class="category-badge">{visualizer.category}</span>
				<button
					type="button"
					class="vote-btn"
					class:voted={visualizer.userVoted}
					onclick={handleVote}
					disabled={voting}
					title={visualizer.userVoted ? 'Remove your upvote' : 'Upvote this visualizer'}
					aria-label={`${visualizer.userVoted ? 'Remove vote from' : 'Vote for'} ${visualizer.name}; ${visualizer.voteCount} votes`}
					aria-pressed={visualizer.userVoted}
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
				rel="noopener noreferrer"
				class="repo-link"
				title="View Open-Source GitHub Repository"
			>
				<svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
					<path
						d="M12 .3a12 12 0 0 0-3.8 23.4c.6.1.8-.3.8-.6v-2.3c-3.3.7-4-1.4-4-1.4-.5-1.4-1.3-1.7-1.3-1.7-1.1-.7.1-.7.1-.7 1.2.1 1.8 1.2 1.8 1.2 1.1 1.8 2.8 1.3 3.5 1 .1-.8.4-1.3.8-1.6-2.7-.3-5.5-1.3-5.5-5.9 0-1.3.5-2.4 1.2-3.2-.1-.3-.5-1.5.1-3.2 0 0 1-.3 3.3 1.2a11.5 11.5 0 0 1 6 0c2.3-1.5 3.3-1.2 3.3-1.2.7 1.7.3 2.9.1 3.2.8.8 1.2 1.9 1.2 3.2 0 4.6-2.8 5.6-5.5 5.9.4.4.8 1.1.8 2.2v3.3c0 .3.2.7.8.6A12 12 0 0 0 12 .3"
					/>
				</svg>
				<span>Repo</span>
				{#if visualizer.githubStars && visualizer.githubStars > 0}
					<span class="star-count">⭐ {formatStars(visualizer.githubStars)}</span>
				{/if}
			</a>
		</div>
	</div>

	<!-- Description -->
	<p class="card-desc">{visualizer.description}</p>

	<!-- Subtle Tags -->
	{#if visualizer.tags && visualizer.tags.length > 0}
		<div class="card-tags">
			{#each visualizer.tags as tag (tag)}
				<span class="card-tag-chip">#{tag}</span>
			{/each}
		</div>
	{/if}

	<!-- Theme Switcher Row -->
	{#if visualizer.themes.length > 1}
		<div class="themes-row">
			<span class="themes-label">Theme:</span>
			<div class="theme-buttons">
				{#each visualizer.themes as t (t.key)}
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
					<h3>External Setup Required</h3>
				</div>
				<p class="setup-text">{visualizer.setupExplanation}</p>
				{#if visualizer.privacyNotice}
					<p class="setup-privacy">{visualizer.privacyNotice}</p>
				{/if}
				<a href={visualizer.websiteUrl} target="_blank" rel="noopener noreferrer" class="setup-btn">
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
						<p class="error-text">
							{#if isReferenceProfile}
								The scheduled Shik3i profile build has not published this preview yet.
							{:else}
								This visualizer requires a workflow in <strong>{username}/{username}</strong>.
							{/if}
						</p>
						{#if visualizer.setupExplanation && !isReferenceProfile}
							<p class="setup-privacy">{visualizer.setupExplanation}</p>
						{/if}
						<div class="error-actions">
							<a
								href={isReferenceProfile
									? 'https://github.com/Shik3i/Shik3i/actions/workflows/visualizers.yml'
									: visualizer.websiteUrl}
								target="_blank"
								rel="noopener noreferrer"
								class="retry-btn"
								>{isReferenceProfile ? 'View Build Status ↗' : 'View Setup Guide ↗'}</a
							>
							<button type="button" class="open-direct-btn" onclick={handleRetry}>Retry</button>
						</div>
					{:else}
						<p class="error-text">Failed to load preview image from third-party provider.</p>
						<div class="error-actions">
							<button type="button" class="retry-btn" onclick={handleRetry}>Retry</button>
							<a href={imageUrl} target="_blank" rel="noopener noreferrer" class="open-direct-btn"
								>Open Image Direct ↗</a
							>
						</div>
					{/if}
				</div>
			{/if}

			{#key `${imageUrl}-${reloadKey}`}
				<img
					bind:this={previewImage}
					src={imageUrl}
					alt="{visualizer.name} preview for {username}"
					width="960"
					height={visualizer.estimatedHeight || 170}
					loading="lazy"
					referrerpolicy="no-referrer"
					class="preview-img"
					class:hidden={imageState === 'error'}
				/>
			{/key}
		{/if}
	</div>

	<!-- Code Preview & Markdown Copy Section -->
	<div class="code-section">
		{#if visualizer.actionWorkflowYaml}
			<div class="workflow-prerequisite">
				This asset does not update dynamically from the username field; generate it in your own
				profile repository with a GitHub Actions workflow by following the
				<a href="/guide#generated-visualizers">setup guide →</a>.
			</div>
			<div class="code-header">
				<span class="code-label">GitHub Actions workflow:</span>
				<div class="code-actions">
					<button type="button" class="copy-btn copy-yaml-btn" onclick={handleCopyYaml}>
						{#if copiedYaml}
							<span class="copied-indicator">Workflow copied</span>
						{:else}
							<span>Copy workflow (.yml)</span>
						{/if}
					</button>
				</div>
			</div>
			<details class="workflow-details">
				<summary class="workflow-summary"
					>View GitHub Action workflow (.github/workflows/*.yml)</summary
				>
				<pre class="code-block yaml-block"><code>{visualizer.actionWorkflowYaml}</code></pre>
			</details>
		{:else}
			<div class="code-header">
				<span class="code-label">README Markdown Code:</span>
				<div class="code-actions">
					<button type="button" class="copy-btn" onclick={handleCopy}>
						{#if copied}
							<span class="copied-indicator">Markdown copied</span>
						{:else}
							<svg
								width="14"
								height="14"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								aria-hidden="true"
							>
								<rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
								<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
							</svg>
							<span>Copy Markdown</span>
						{/if}
					</button>
				</div>
			</div>
			<pre class="code-block"><code>{markdownCode}</code></pre>
		{/if}
		{#if actionStatus}
			<p class:error={actionFailed} class="action-status" role="status">{actionStatus}</p>
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
		transition:
			border-color 0.2s ease,
			box-shadow 0.2s ease;
		min-width: 0;
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
		min-width: 0;
	}

	.card-title {
		font-size: 1.05rem;
		font-weight: 700;
		color: var(--text-main);
		line-height: 1.3;
		overflow-wrap: anywhere;
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

	.vote-btn {
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
		transition:
			background-color 0.15s cubic-bezier(0.4, 0, 0.2, 1),
			border-color 0.15s cubic-bezier(0.4, 0, 0.2, 1),
			color 0.15s cubic-bezier(0.4, 0, 0.2, 1),
			box-shadow 0.15s cubic-bezier(0.4, 0, 0.2, 1),
			transform 0.15s cubic-bezier(0.4, 0, 0.2, 1);
		box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
		min-height: 2rem;
	}

	.vote-btn .upvote-arrow {
		transition:
			transform 0.15s ease,
			color 0.15s ease;
	}

	.vote-btn:hover {
		border-color: var(--brand-primary);
		color: var(--brand-primary);
		background-color: color-mix(in srgb, var(--brand-primary) 9%, transparent);
	}

	.vote-btn:hover .upvote-arrow {
		transform: translateY(-2px);
	}

	.vote-btn.voted {
		background: linear-gradient(135deg, var(--brand-solid) 0%, #0d9488 100%);
		color: var(--brand-on-solid);
		border-color: var(--brand-primary);
		box-shadow: 0 2px 8px color-mix(in srgb, var(--brand-primary) 35%, transparent);
	}

	.vote-btn.voted .upvote-arrow {
		color: #ffffff;
		transform: translateY(-1px);
	}

	.vote-btn:active {
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
		transition:
			background-color 0.15s ease,
			border-color 0.15s ease,
			color 0.15s ease;
		min-height: 2rem;
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
		transition:
			opacity 0.15s ease,
			color 0.15s ease,
			border-color 0.15s ease;
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
		transition:
			background-color 0.15s ease,
			border-color 0.15s ease,
			color 0.15s ease;
		min-height: 2rem;
	}

	.theme-chip:hover,
	.theme-chip.active {
		background-color: var(--brand-solid);
		color: var(--brand-on-solid);
		border-color: var(--brand-solid);
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
		to {
			transform: rotate(360deg);
		}
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

	.retry-btn,
	.open-direct-btn {
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

	.setup-notice,
	.opt-in-box {
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

	.setup-header h3 {
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

	.setup-privacy,
	.counter-notice {
		font-size: 0.775rem;
		color: var(--text-subtle);
		font-style: italic;
	}

	.setup-btn,
	.opt-in-btn {
		margin-top: 0.4rem;
		font-size: 0.825rem;
		font-weight: 700;
		padding: 0.45rem 0.85rem;
		border-radius: var(--radius-md);
		background-color: var(--brand-solid);
		color: var(--brand-on-solid);
		transition: background-color 0.15s ease;
	}

	.setup-btn:hover,
	.opt-in-btn:hover {
		background-color: var(--brand-solid-hover);
		color: var(--brand-on-solid);
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
		transition:
			background-color 0.15s ease,
			border-color 0.15s ease,
			color 0.15s ease;
	}

	.workflow-prerequisite {
		padding: 0.55rem 0.7rem;
		border-left: 3px solid var(--warning-text);
		background: color-mix(in srgb, var(--warning-bg) 75%, transparent);
		color: var(--text-main);
		font-size: 0.8rem;
		line-height: 1.45;
	}

	.workflow-prerequisite a {
		font-weight: 700;
		text-decoration: underline;
		text-underline-offset: 0.15em;
	}

	.copy-yaml-btn {
		background-color: color-mix(in srgb, var(--brand-primary) 12%, transparent);
		color: var(--brand-primary);
		border-color: color-mix(in srgb, var(--brand-primary) 30%, transparent);
	}

	.copy-yaml-btn:hover {
		background-color: color-mix(in srgb, var(--brand-primary) 22%, transparent);
		color: var(--brand-hover);
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

	.action-status {
		font-size: 0.75rem;
		color: #4ade80;
	}

	.action-status.error {
		color: #fca5a5;
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

	@media (max-width: 420px) {
		.card {
			padding: 1rem;
		}

		.card-header {
			flex-wrap: wrap;
		}

		.header-main {
			width: 100%;
		}

		.header-links {
			width: 100%;
		}

		.repo-link {
			width: 100%;
			justify-content: center;
		}

		.error-actions {
			flex-wrap: wrap;
			justify-content: center;
		}
	}
</style>
