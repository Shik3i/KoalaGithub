<script lang="ts">
	import { DEFAULT_USERNAME, isValidGitHubUsername, sanitizeUsername } from '$lib/utils/username';

	let { username = $bindable(DEFAULT_USERNAME), onUpdate = (_val: string) => {} } = $props<{
		username?: string;
		onUpdate?: (val: string) => void;
	}>();

	let inputVal = $state(username);
	let errorMessage = $state<string | null>(null);

	$effect(() => {
		inputVal = username;
	});

	function handleSubmit(e?: Event) {
		if (e) e.preventDefault();
		const trimmed = inputVal.trim();

		if (!trimmed) {
			inputVal = DEFAULT_USERNAME;
			errorMessage = null;
			username = DEFAULT_USERNAME;
			onUpdate(DEFAULT_USERNAME);
			return;
		}

		if (!isValidGitHubUsername(trimmed)) {
			errorMessage = 'Invalid GitHub username format (letters, numbers, single hyphens).';
			return;
		}

		errorMessage = null;
		const clean = sanitizeUsername(trimmed);
		username = clean;
		onUpdate(clean);
	}

	function handleReset() {
		inputVal = DEFAULT_USERNAME;
		errorMessage = null;
		username = DEFAULT_USERNAME;
		onUpdate(DEFAULT_USERNAME);
	}
</script>

<div class="search-container">
	<form class="search-box" onsubmit={handleSubmit}>
		<div class="input-wrapper">
			<span class="github-icon">
				<svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
					<path
						d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"
					/>
				</svg>
			</span>

			<input
				type="text"
				bind:value={inputVal}
				placeholder="Enter GitHub Username (e.g. Shik3i, octocat)"
				class="username-input"
				aria-label="GitHub Username"
				spellcheck="false"
				autocomplete="off"
			/>

			{#if inputVal !== DEFAULT_USERNAME}
				<button type="button" class="reset-btn" onclick={handleReset} title="Reset to default (Shik3i)">
					✕
				</button>
			{/if}
		</div>

		<button type="submit" class="submit-btn">
			<span>Update Profile Cards</span>
			<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
				<path d="M5 12h14M12 5l7 7-7 7" stroke-linecap="round" stroke-linejoin="round" />
			</svg>
		</button>
	</form>

	{#if errorMessage}
		<p class="error-msg">{errorMessage}</p>
	{/if}

	<div class="quick-examples">
		<span class="label">Try examples:</span>
		<button
			type="button"
			class="chip-btn"
			class:active={username === 'Shik3i'}
			onclick={() => {
				inputVal = 'Shik3i';
				handleSubmit();
			}}
		>
			Shik3i
		</button>
		<button
			type="button"
			class="chip-btn"
			class:active={username === 'octocat'}
			onclick={() => {
				inputVal = 'octocat';
				handleSubmit();
			}}
		>
			octocat
		</button>
		<button
			type="button"
			class="chip-btn"
			class:active={username === 'torvalds'}
			onclick={() => {
				inputVal = 'torvalds';
				handleSubmit();
			}}
		>
			torvalds
		</button>
	</div>
</div>

<style>
	.search-container {
		display: flex;
		flex-direction: column;
		gap: 0.6rem;
		width: 100%;
		max-width: 760px;
		margin: 0 auto;
	}

	.search-box {
		display: flex;
		gap: 0.6rem;
		align-items: center;
	}

	.input-wrapper {
		position: relative;
		display: flex;
		align-items: center;
		flex: 1;
		background-color: var(--bg-card);
		border: 1.5px solid var(--border-color);
		border-radius: var(--radius-md);
		box-shadow: var(--shadow-sm);
		transition: border-color 0.2s ease, box-shadow 0.2s ease;
	}

	.input-wrapper:focus-within {
		border-color: var(--brand-primary);
		box-shadow: 0 0 0 3px var(--brand-light);
	}

	.github-icon {
		padding-left: 0.9rem;
		display: flex;
		align-items: center;
		color: var(--text-muted);
	}

	.username-input {
		width: 100%;
		padding: 0.75rem 0.75rem 0.75rem 0.6rem;
		border: none;
		background: transparent;
		color: var(--text-main);
		font-weight: 600;
		font-size: 1.05rem;
		outline: none;
	}

	.reset-btn {
		padding: 0.4rem 0.75rem;
		color: var(--text-muted);
		font-size: 0.85rem;
		transition: color 0.15s ease;
	}

	.reset-btn:hover {
		color: var(--text-main);
	}

	.submit-btn {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		background-color: var(--brand-primary);
		color: #ffffff;
		font-weight: 700;
		font-size: 0.95rem;
		padding: 0.75rem 1.25rem;
		border-radius: var(--radius-md);
		white-space: nowrap;
		transition: background-color 0.15s ease, transform 0.1s ease;
		box-shadow: var(--shadow-sm);
	}

	.submit-btn:hover {
		background-color: var(--brand-hover);
	}

	.submit-btn:active {
		transform: translateY(1px);
	}

	.error-msg {
		color: var(--error-text);
		font-size: 0.85rem;
		font-weight: 500;
		padding-left: 0.25rem;
	}

	.quick-examples {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.825rem;
		color: var(--text-muted);
	}

	.chip-btn {
		background-color: var(--bg-subtle);
		color: var(--text-muted);
		font-weight: 600;
		padding: 0.2rem 0.6rem;
		border-radius: var(--radius-full);
		transition: all 0.15s ease;
		border: 1px solid var(--border-color);
	}

	.chip-btn:hover, .chip-btn.active {
		background-color: var(--brand-light);
		color: var(--brand-text);
		border-color: var(--brand-primary);
	}

	@media (max-width: 640px) {
		.search-box {
			flex-direction: column;
			align-items: stretch;
		}

		.submit-btn {
			justify-content: center;
		}
	}
</style>
