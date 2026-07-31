<svelte:head>
	<title>About KoalaGitHub – Central GitHub Profile Visualizer Hub</title>
	<link rel="canonical" href="https://github.koalastuff.net/about" />
	<meta property="og:url" content="https://github.koalastuff.net/about" />
</svelte:head>

<div class="container page-container">
	<header class="page-header">
		<h1>About KoalaGitHub</h1>
		<p class="lead">
			Centralized, open-source directory for GitHub profile stats, graphs, trophies, and badges.
		</p>
	</header>

	<div class="content-card">
		<h2>What is KoalaGitHub?</h2>
		<p>
			<strong>KoalaGitHub</strong> (<code>github.koalastuff.net</code>) solves a common problem
			developers face when designing their GitHub profile README: instead of hunting through dozens
			of separate websites and repositories to test different stats cards, streak graphs, activity
			charts, and trophies, KoalaGitHub brings them together on a single, responsive comparison
			dashboard.
		</p>

		<h2>Core principles</h2>
		<ul class="features-list">
			<li>
				⚡ <strong>Fast UI</strong>: SvelteKit prerenders the interface; a small Go API provides
				live stars and pseudonymous votes.
			</li>
			<li>
				🔒 <strong>Privacy Focused</strong>: No analytics, ad trackers, accounts, or tracking
				cookies.
			</li>
			<li>
				🔗 <strong>Shareable Custom URLs</strong>: Share your customized view using
				<code>?user=yourusername</code>.
			</li>
			<li>
				🛠️ <strong>Guided setup</strong>: Copy Markdown for instant visualizers or the required
				workflow for repository-generated assets.
			</li>
		</ul>

		<h2>How to contribute a new visualizer</h2>
		<p>
			Visualizer metadata is maintained in the central JSON registry <code
				>src/lib/data/visualizers.json</code
			>.
		</p>
		<ol class="steps-list">
			<li>
				Fork the repository on GitHub: <a
					href="https://github.com/Shik3i/KoalaGithub"
					target="_blank"
					rel="noopener noreferrer">github.com/Shik3i/KoalaGithub</a
				>
			</li>
			<li>Add your visualizer definition to <code>src/lib/data/visualizers.json</code>.</li>
			<li>Run <code>npm run test</code> to run data integrity checks.</li>
			<li>Submit a pull request.</li>
		</ol>

		<h2>How sorting works</h2>
		<ul class="features-list">
			<li><strong>Most Voted</strong> uses only votes recorded by KoalaGitHub.</li>
			<li>
				<strong>Most Stars</strong> uses only the current star count of each source repository.
			</li>
			<li><strong>Alphabetical</strong> sorts by visualizer name.</li>
			<li><strong>Recently Added</strong> sorts by the registry's addition date.</li>
		</ul>
		<p>Votes and repository stars are never combined into a weighted score.</p>

		<div class="eco-note">
			<h3>Part of the KoalaStuff ecosystem</h3>
			<p>
				KoalaGitHub is part of the open-source KoalaStuff project suite. For legal and domain
				notices, please visit <a
					href="https://koalastuff.net/legal"
					target="_blank"
					rel="noopener noreferrer">koalastuff.net/legal ↗</a
				>.
			</p>
		</div>
	</div>
</div>

<style>
	.page-container {
		max-width: 840px;
		margin: 2.5rem auto 4rem;
	}

	.page-header {
		margin-bottom: 2rem;
		text-align: center;
	}

	.page-header h1 {
		font-size: 2.25rem;
		font-weight: 800;
		color: var(--text-main);
	}

	.lead {
		font-size: 1.1rem;
		color: var(--text-muted);
		margin-top: 0.4rem;
	}

	.content-card {
		background-color: var(--bg-card);
		border: 1px solid var(--border-color);
		border-radius: var(--radius-lg);
		padding: 2rem;
		display: flex;
		flex-direction: column;
		gap: 1.25rem;
		box-shadow: var(--shadow-sm);
	}

	h2 {
		font-size: 1.35rem;
		font-weight: 700;
		color: var(--text-main);
		border-bottom: 1px solid var(--border-color);
		padding-bottom: 0.5rem;
	}

	p {
		color: var(--text-main);
		line-height: 1.6;
	}

	.features-list,
	.steps-list {
		display: flex;
		flex-direction: column;
		gap: 0.6rem;
		padding-left: 1.25rem;
		color: var(--text-main);
	}

	.features-list li,
	.steps-list li {
		line-height: 1.5;
	}

	.eco-note {
		margin-top: 1rem;
		background-color: var(--bg-subtle);
		padding: 1.25rem;
		border-radius: var(--radius-md);
		border: 1px solid var(--border-color);
	}

	.eco-note h3 {
		font-size: 1.05rem;
		font-weight: 700;
		margin-bottom: 0.35rem;
	}
</style>
