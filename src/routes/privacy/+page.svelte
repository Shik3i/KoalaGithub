<svelte:head>
	<title>Privacy Policy – KoalaGitHub</title>
</svelte:head>

<div class="container page-container">
	<header class="page-header">
		<h1>Privacy Policy</h1>
		<p class="lead">Transparent information about data handling on KoalaGitHub (`github.koalastuff.net`).</p>
	</header>

	<div class="content-card">
		<h2>1. Application Data Processing</h2>
		<p>
			KoalaGitHub is designed as a <strong>100% static, client-side web application</strong>.
		</p>
		<ul>
			<li><strong>No Analytics</strong>: We do not use Google Analytics, Plausible, Telemetry, or any tracking services.</li>
			<li><strong>No Cookies or Ad Networks</strong>: We do not set any tracking cookies or execute advertising scripts.</li>
			<li><strong>No Account System or Database</strong>: We do not collect, store, or process user accounts or personal profiles on any server.</li>
		</ul>

		<h2>2. Local Storage</h2>
		<p>
			We use your browser's <code>localStorage</code> solely to save your active theme preference (e.g. Light, Dark, or System theme). This setting remains strictly inside your web browser and is never transmitted to us or any third party.
		</p>

		<h2>3. Third-Party Visualizer Requests & Network Traffic</h2>
		<p>
			<strong>Important Transparency Notice:</strong> When you view or interact with live visualizers, preview images, stats cards, trophies, or graphs on KoalaGitHub, your web browser makes direct HTTP/HTTPS requests to third-party providers (such as Vercel, Heroku, GitHub, komarev.com, shields.io, etc.).
		</p>
		<p>
			During these requests, third-party hosting services and visualizer maintainers receive standard web traffic data directly from your browser, which may include:
		</p>
		<ul>
			<li>Your IP address</li>
			<li>Your browser User-Agent header</li>
			<li>The requested GitHub username parameter</li>
			<li>Timestamp and standard HTTP request headers</li>
		</ul>
		<p>
			To help protect your privacy, KoalaGitHub applies a restrictive <code>referrerpolicy="no-referrer"</code> attribute to preview images and external links where practical.
		</p>

		<h2>4. Visitor Counters</h2>
		<p>
			Certain visualizers (such as visitor profile view counters) increment an external counter database when their preview image URL is fetched. To prevent unwanted counting, KoalaGitHub requires explicit user action before loading live preview images for counter badges.
		</p>

		<h2>5. Fonts & Assets</h2>
		<p>
			KoalaGitHub does not load font files or stylesheets from external CDNs (such as Google Fonts). All typography relies on native system fonts.
		</p>

		<h2>6. Legal Notices</h2>
		<p>
			For general legal notices, terms, and operator details of the KoalaStuff ecosystem, please refer to <a href="https://koalastuff.net/legal" target="_blank" rel="noreferrer">koalastuff.net/legal ↗</a>.
		</p>
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
		font-size: 1.25rem;
		font-weight: 700;
		color: var(--text-main);
		border-bottom: 1px solid var(--border-color);
		padding-bottom: 0.4rem;
		margin-top: 0.5rem;
	}

	p, ul {
		color: var(--text-main);
		line-height: 1.6;
	}

	ul {
		padding-left: 1.25rem;
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
	}
</style>
