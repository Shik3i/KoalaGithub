<svelte:head>
	<title>Privacy Policy – KoalaGitHub</title>
</svelte:head>

<div class="container page-container">
	<header class="page-header">
		<h1>Privacy Policy</h1>
		<p class="lead">Transparent information about data handling, anonymous upvoting, and privacy on KoalaGitHub (<code>github.koalastuff.net</code>).</p>
	</header>

	<div class="content-card">
		<h2>1. Overview & Backend Architecture</h2>
		<p>
			KoalaGitHub is an open-source web application designed with a strict <strong>privacy-first approach</strong>. It consists of a static SvelteKit frontend served by an embedded Go backend server.
		</p>
		<ul>
			<li><strong>No User Accounts or Registration</strong>: You do not need to register, log in, or provide any personal information (such as name or email) to use KoalaGitHub.</li>
			<li><strong>No Tracking Cookies or Ad Networks</strong>: We do not set any tracking cookies, execute advertising scripts, or use third-party analytics (no Google Analytics, no Plausible, no Telemetry).</li>
			<li><strong>No IP Logging for Votes</strong>: IP addresses are not stored or associated with user actions.</li>
		</ul>

		<h2>2. Anonymous Device-Based Upvoting System</h2>
		<p>
			To allow visitors to upvote their favorite visualizers without requiring an account or login, KoalaGitHub uses a pseudonymous **Device ID**:
		</p>
		<ul>
			<li><strong>Device Identifier (UUID v4)</strong>: When you visit KoalaGitHub, a random UUID v4 string (e.g. <code>koala_device_id</code>) is generated locally in your browser and saved in <code>localStorage</code>.</li>
			<li><strong>Vote Registration</strong>: When you click an upvote button, your pseudonymous <code>device_id</code> and the target <code>visualizer_id</code> are transmitted to <code>POST /api/visualizers/&#123;id&#125;/vote</code> and stored in our local SQLite database.</li>
			<li><strong>Purpose</strong>: This identifier is used exclusively to keep track of which visualizers you have upvoted on your current device so you can toggle your votes on or off. It is completely pseudonymous and is never linked to your identity, IP address, or GitHub account.</li>
		</ul>

		<h2>3. Local Storage Usage</h2>
		<p>
			We use your browser's <code>localStorage</code> for two specific functional purposes only:
		</p>
		<ol>
			<li><code>koala-theme</code>: Stores your selected interface theme (Light, Dark, or System mode).</li>
			<li><code>koala_device_id</code>: Stores your pseudonymous UUID v4 device identifier for visualizer upvoting.</li>
		</ol>
		<p>
			Both entries remain inside your browser and can be cleared at any time by clearing your browser site data or local storage.
		</p>

		<h2>4. Third-Party Visualizer Requests & Image Previews</h2>
		<p>
			<strong>Important Transparency Notice:</strong> When you view live visualizer preview images, stats cards, graphs, or bonsai trees on KoalaGitHub, your browser makes direct HTTPS requests to third-party hosting services (such as Vercel, GitHub raw content, demolab, etc.).
		</p>
		<p>
			During these HTTP requests, third-party hosting services and visualizer maintainers receive standard web traffic data directly from your browser, which may include:
		</p>
		<ul>
			<li>Your IP address</li>
			<li>Your browser User-Agent header</li>
			<li>The requested GitHub username parameter</li>
			<li>Standard HTTP request timestamp</li>
		</ul>
		<p>
			To protect your privacy, KoalaGitHub applies a restrictive <code>referrerpolicy="no-referrer"</code> attribute to all preview images and external links where practical.
		</p>

		<h2>5. Fonts & Asset Hosting</h2>
		<p>
			KoalaGitHub does not load font files, scripts, or stylesheets from external CDNs (such as Google Fonts). All typography relies on native system font stacks.
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

	p, ul, ol {
		color: var(--text-main);
		line-height: 1.6;
	}

	ul, ol {
		padding-left: 1.25rem;
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
	}
</style>
