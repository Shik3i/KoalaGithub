<script lang="ts">
	import KoalaLogo from '$lib/components/KoalaLogo.svelte';
</script>

<svelte:head>
	<title>Guide: How to Create GitHub Profile README & Action Workflows | KoalaGitHub</title>
	<meta name="description" content="Step-by-step guide on creating a special GitHub profile README, running automated daily GitHub Action workflows, and adding animated contribution graphs." />
</svelte:head>

<div class="guide-page container">
	<!-- Hero Section -->
	<div class="guide-hero">
		<KoalaLogo size={56} />
		<h1 class="guide-title">GitHub Profile README & Action Workflows Guide</h1>
		<p class="guide-subtitle">
			Everything you need to know about setting up your special profile README, running daily automated GitHub Action visualizers, and showcasing live contribution graphs.
		</p>
	</div>

	<div class="guide-content">
		<!-- Step 1: Special Repository -->
		<section class="guide-card">
			<div class="card-step-badge">Step 1</div>
			<h2>Creating Your Special Profile Repository</h2>
			<p>
				GitHub has a special feature: when you create a repository with the <strong>exact same name as your GitHub username</strong>, GitHub displays the contents of that repository's <code>README.md</code> at the top of your public profile page!
			</p>
			
			<div class="info-box">
				<span class="box-icon">💡</span>
				<div>
					<strong>Example:</strong> For username <code>Shik3i</code>, the repository is located at
					<a href="https://github.com/Shik3i/Shik3i" target="_blank" rel="noreferrer">github.com/Shik3i/Shik3i ↗</a>.
				</div>
			</div>

			<ol class="step-list">
				<li>Go to <a href="https://github.com/new" target="_blank" rel="noreferrer">github.com/new ↗</a>.</li>
				<li>In the <strong>Repository name</strong> field, enter your exact GitHub username.</li>
				<li>Ensure the repository is set to <strong>Public</strong>.</li>
				<li>Check <strong>"Add a README file"</strong>.</li>
				<li>Click <strong>Create repository</strong>.</li>
			</ol>
		</section>

		<!-- Step 2: Instant Visualizers -->
		<section class="guide-card">
			<div class="card-step-badge">Step 2</div>
			<h2>Adding Instant Visualizers (Streak Stats, Activity Graphs, Bonsai)</h2>
			<p>
				Instant visualizers (like <strong>Kodama Bonsai Tree</strong>, <strong>GitHub Streak Stats</strong>, or <strong>Activity Graphs</strong>) do not require any configuration in your repository. They generate SVG images on the fly!
			</p>

			<ol class="step-list">
				<li>Search your username at <a href="/">KoalaGitHub Hub</a>.</li>
				<li>Select your favorite theme for any visualizer card.</li>
				<li>Click <strong>Copy Markdown</strong>.</li>
				<li>Edit your profile <code>README.md</code> on GitHub and paste the snippet where you want it to appear!</li>
			</ol>
		</section>

		<!-- Step 3: Self-Hosted GitHub Actions -->
		<section class="guide-card">
			<div class="card-step-badge">Step 3</div>
			<h2>Setting Up Daily GitHub Action Workflows (Snake, Pac-Man, 3D Contrib)</h2>
			<p>
				Some visualizers (such as the <strong>Contribution Grid Snake</strong> or <strong>Arcade Pac-Man Graph</strong>) are interactive animations compiled automatically once a day by a <strong>GitHub Action workflow</strong> running in your repository.
			</p>

			<h3>How to Add a Workflow File:</h3>
			<ol class="step-list">
				<li>In your profile repository, create a directory path named <code>.github/workflows/</code>.</li>
				<li>Create a new file inside it, for example <code>.github/workflows/snake.yml</code>.</li>
				<li>Copy the 1-click <strong>📋 Copy Workflow (.yml)</strong> code from KoalaGitHub into that file.</li>
				<li>Commit the file to your repository's <code>main</code> branch.</li>
			</ol>

			<div class="code-example-box">
				<div class="box-title">Example Workflow Configuration (.github/workflows/snake.yml):</div>
				<pre class="yaml-code"><code>name: Generate Snake

on:
  schedule:
    - cron: "0 0 * * *" # Runs automatically every night at 00:00 UTC
  workflow_dispatch: # Allows manual trigger anytime from Actions tab

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: Platane/snk@v3
        with:
          github_user_name: ${'{'} secrets.GITHUB_USERNAME }
          outputs: |
            output/snake.svg
            output/snake-dark.svg?palette=github-dark
      - uses: crazy-max/ghaction-github-pages@v31
        with:
          target_branch: output
          build_dir: output
        env:
          GITHUB_TOKEN: ${'{'} secrets.GITHUB_TOKEN }</code></pre>
			</div>

			<div class="warning-box">
				<span class="box-icon">🔐</span>
				<div>
					<strong>Required Repository Permission:</strong><br />
					Go to your repo's <strong>Settings → Actions → General → Workflow permissions</strong> and select <strong>"Read and write permissions"</strong>. This permits the workflow to save the compiled SVG image to the <code>output</code> branch!
				</div>
			</div>
		</section>

		<!-- Step 4: Real-World Reference Repo -->
		<section class="guide-card reference-card">
			<div class="card-step-badge reference-badge">Live Reference</div>
			<h2>Real-World Example Repository</h2>
			<p>
				Want to see a fully configured working profile README in action? Check out the reference implementation:
			</p>

			<div class="reference-links">
				<a href="https://github.com/Shik3i/Shik3i" target="_blank" rel="noreferrer" class="ref-link-btn">
					<svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
						<path fill-rule="evenodd" clip-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.53 1.032 1.53 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" />
					</svg>
					View Shik3i Profile Repository ↗
				</a>
				<a href="https://github.com/Shik3i/Shik3i/tree/main/.github/workflows" target="_blank" rel="noreferrer" class="ref-link-btn secondary">
					📂 View Live Workflow Files (.github/workflows) ↗
				</a>
			</div>
		</section>
	</div>
</div>

<style>
	.guide-page {
		padding: 2.5rem 1.25rem 4rem;
	}

	.guide-hero {
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		gap: 0.75rem;
		margin-bottom: 3rem;
	}

	.guide-title {
		font-size: 2.25rem;
		font-weight: 800;
		color: var(--text-main);
		letter-spacing: -0.03em;
	}

	.guide-subtitle {
		font-size: 1.1rem;
		color: var(--text-muted);
		max-width: 680px;
		line-height: 1.55;
	}

	.guide-content {
		display: flex;
		flex-direction: column;
		gap: 2rem;
		max-width: 860px;
		margin: 0 auto;
	}

	.guide-card {
		background-color: var(--bg-card);
		border: 1px solid var(--border-color);
		border-radius: var(--radius-lg);
		padding: 1.75rem;
		position: relative;
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.card-step-badge {
		align-self: flex-start;
		font-size: 0.75rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: var(--brand-primary);
		background-color: var(--brand-light);
		padding: 0.2rem 0.6rem;
		border-radius: var(--radius-full);
	}

	.guide-card h2 {
		font-size: 1.35rem;
		font-weight: 700;
		color: var(--text-main);
	}

	.guide-card h3 {
		font-size: 1.1rem;
		font-weight: 600;
		color: var(--text-main);
		margin-top: 0.5rem;
	}

	.guide-card p {
		color: var(--text-muted);
		line-height: 1.6;
		font-size: 0.975rem;
	}

	.step-list {
		padding-left: 1.25rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		color: var(--text-main);
		font-size: 0.95rem;
		line-height: 1.5;
	}

	.info-box, .warning-box {
		display: flex;
		align-items: flex-start;
		gap: 0.75rem;
		padding: 1rem;
		border-radius: var(--radius-md);
		font-size: 0.9rem;
		line-height: 1.5;
	}

	.info-box {
		background-color: var(--koala-accent-bg);
		color: var(--text-main);
		border: 1px solid rgba(14, 165, 233, 0.3);
	}

	.warning-box {
		background-color: var(--warning-bg);
		color: var(--warning-text);
		border: 1px solid rgba(245, 158, 11, 0.3);
	}

	.box-icon {
		font-size: 1.25rem;
		line-height: 1;
	}

	.code-example-box {
		background-color: var(--bg-code);
		border-radius: var(--radius-md);
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.box-title {
		font-size: 0.8rem;
		font-weight: 600;
		color: #94a3b8;
		text-transform: uppercase;
		letter-spacing: 0.03em;
	}

	.yaml-code {
		font-family: var(--font-mono);
		font-size: 0.825rem;
		color: #e2e8f0;
		overflow-x: auto;
		white-space: pre;
		margin: 0;
	}

	.reference-card {
		border-color: rgba(99, 102, 241, 0.3);
		background: linear-gradient(135deg, var(--bg-card) 0%, rgba(99, 102, 241, 0.03) 100%);
	}

	.reference-badge {
		background-color: rgba(99, 102, 241, 0.15);
		color: #818cf8;
	}

	.reference-links {
		display: flex;
		gap: 0.75rem;
		flex-wrap: wrap;
		margin-top: 0.5rem;
	}

	.ref-link-btn {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.6rem 1.1rem;
		border-radius: var(--radius-md);
		font-weight: 600;
		font-size: 0.9rem;
		background-color: var(--brand-primary);
		color: #ffffff;
		transition: background-color 0.15s ease;
	}

	.ref-link-btn:hover {
		background-color: var(--brand-hover);
		color: #ffffff;
	}

	.ref-link-btn.secondary {
		background-color: var(--bg-subtle);
		color: var(--text-main);
		border: 1px solid var(--border-color);
	}

	.ref-link-btn.secondary:hover {
		background-color: var(--border-color);
	}

	@media (max-width: 640px) {
		.guide-title {
			font-size: 1.75rem;
		}

		.guide-card {
			padding: 1.25rem;
		}
	}
</style>
