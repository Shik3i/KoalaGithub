<script lang="ts">
	import { clearDeviceId, getStoredDeviceId } from '$lib/utils/device';

	let deleting = $state(false);
	let deletionStatus = $state('');
	let deletionError = $state(false);

	async function deleteVotes() {
		const deviceId = getStoredDeviceId();
		if (!deviceId) {
			deletionError = false;
			deletionStatus = 'No vote identifier is stored in this browser.';
			return;
		}
		deleting = true;
		deletionStatus = '';
		try {
			const response = await fetch('/api/votes', {
				method: 'DELETE',
				headers: { 'X-Device-ID': deviceId }
			});
			if (!response.ok) throw new Error(`HTTP ${response.status}`);
			const result: unknown = await response.json();
			const deleted =
				result &&
				typeof result === 'object' &&
				typeof (result as { deleted?: unknown }).deleted === 'number'
					? (result as { deleted: number }).deleted
					: 0;
			clearDeviceId();
			deletionError = false;
			deletionStatus = `${deleted} stored vote${deleted === 1 ? '' : 's'} deleted. The local identifier was removed too.`;
		} catch {
			deletionError = true;
			deletionStatus = 'Votes could not be deleted right now. Please try again later.';
		} finally {
			deleting = false;
		}
	}
</script>

<svelte:head>
	<title>Privacy Policy – KoalaGitHub</title>
	<meta
		name="description"
		content="How KoalaGitHub processes device identifiers, votes, server requests, local storage, and third-party preview requests."
	/>
	<link rel="canonical" href="https://github.koalastuff.net/privacy" />
	<meta property="og:url" content="https://github.koalastuff.net/privacy" />
</svelte:head>

<div class="container page-container">
	<header class="page-header">
		<p class="eyebrow">Last updated: 31 July 2026</p>
		<h1>Privacy Policy</h1>
		<p class="lead">
			How data is processed when you use KoalaGitHub at <code>github.koalastuff.net</code>.
		</p>
	</header>

	<div class="content-card">
		<section>
			<h2>1. Controller and contact</h2>
			<p>
				The operator identified in the
				<a href="https://koalastuff.net/legal" target="_blank" rel="noopener noreferrer"
					>KoalaStuff legal notice ↗</a
				>
				is responsible for this service. Use the contact details published there for privacy questions
				or requests.
			</p>
		</section>

		<section>
			<h2>2. Data processed by the service</h2>
			<p>
				No account, name, email address, tracking cookie, advertising script, or analytics service
				is required.
			</p>
			<ul>
				<li>
					<strong>Server requests:</strong> The web server and infrastructure providers necessarily process
					connection data such as IP address, request time, path, user agent, and response status to deliver
					and protect the service. KoalaGitHub does not write IP addresses into its SQLite vote database.
					Infrastructure security or access logs may be retained separately.
				</li>
				<li>
					<strong>GitHub usernames:</strong> A username entered in the search field is placed in the page
					URL and compatible preview URLs. It is not written to the SQLite database. Infrastructure providers
					and the selected preview provider may receive it as part of the requested URL.
				</li>
				<li>
					<strong>Votes:</strong> A random UUID v4 device identifier and the selected visualizer ID are
					transmitted to the API and stored in SQLite. They are not intentionally connected to a GitHub
					account, name, or email address.
				</li>
				<li>
					<strong>Abuse protection:</strong> The application temporarily counts vote and vote-deletion
					requests per combination of IP address and device identifier in memory. These rate-limit counters
					reset after one hour and are not persisted. Resetting a counter does not expire, remove, or
					duplicate a vote: a vote remains associated with the same persistent browser identifier until
					it is toggled off or deleted.
				</li>
			</ul>
		</section>

		<section>
			<h2>3. Purpose, legal basis, and retention</h2>
			<p>
				Functional storage is used to provide vote toggling, aggregate vote totals, service
				delivery, and abuse prevention. The intended legal basis is legitimate interests under
				Article 6(1)(f) GDPR: operating a useful and secure community directory without user
				accounts or cross-site tracking.
			</p>
			<p>
				Vote records remain until you remove the vote, use the deletion control below, the
				visualizer is removed, or the service is discontinued. Only the transient in-memory request
				counters reset after one hour; this has no effect on stored votes. Hosting logs follow the
				infrastructure provider's operational retention rules.
			</p>
		</section>

		<section>
			<h2>4. Browser storage</h2>
			<ul>
				<li><code>koala-theme</code> stores the selected light, dark, or system theme.</li>
				<li>
					<code>koala_device_id</code> stores the pseudonymous identifier used to load and toggle your
					vote state. It has no automatic expiry and remains unchanged until you delete it or clear this
					site's browser data.
				</li>
			</ul>
			<p>
				Clearing site data removes these browser entries. It does not by itself delete an existing
				vote record from the server because the identifier needed to locate that record would be
				lost. Use the control below first.
			</p>
			<div class="delete-box">
				<div>
					<h3>Delete votes from this browser</h3>
					<p>
						Deletes every server-side vote associated with the currently stored device identifier.
					</p>
				</div>
				<button type="button" class="danger-button" onclick={deleteVotes} disabled={deleting}>
					{deleting ? 'Deleting…' : 'Delete my votes'}
				</button>
			</div>
			{#if deletionStatus}
				<p class:error={deletionError} class="status" role="status">{deletionStatus}</p>
			{/if}
		</section>

		<section>
			<h2>5. Third-party previews and links</h2>
			<p>
				Visualizer previews are requested directly by your browser from their respective providers,
				including GitHub-hosted content and independent services. Those providers receive normal
				HTTP request data, including your IP address, browser information, request time, and the
				GitHub username embedded in the preview URL. KoalaGitHub uses
				<code>referrerpolicy="no-referrer"</code> for preview images and external links where practical,
				but each provider remains responsible for its own processing.
			</p>
			<p>
				Fonts, application scripts, and styles are served locally; no analytics or advertising
				network is included.
			</p>
		</section>

		<section>
			<h2>6. Recipients, transfers, and your rights</h2>
			<p>
				Data may be processed by the service operator, its hosting and network providers, and the
				third-party preview provider you choose to load. Some preview providers may process data
				outside the EU/EEA under their own terms.
			</p>
			<p>
				Subject to applicable law, you may request access, correction, deletion, restriction,
				objection, or data portability. You may also lodge a complaint with a competent data
				protection authority. Because votes are pseudonymous, include the device identifier when a
				request concerns a vote record; otherwise the record cannot be reliably attributed to you.
			</p>
			<p>KoalaGitHub does not use automated decision-making or profiling.</p>
		</section>
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
	.eyebrow {
		color: var(--brand-primary);
		font-size: 0.8rem;
		font-weight: 700;
		letter-spacing: 0.04em;
		text-transform: uppercase;
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
		gap: 2rem;
		box-shadow: var(--shadow-sm);
	}
	section {
		display: grid;
		gap: 0.8rem;
	}
	h2 {
		font-size: 1.25rem;
		font-weight: 700;
		color: var(--text-main);
		border-bottom: 1px solid var(--border-color);
		padding-bottom: 0.4rem;
	}
	p,
	ul {
		color: var(--text-main);
		line-height: 1.65;
	}
	ul {
		padding-left: 1.25rem;
		display: flex;
		flex-direction: column;
		gap: 0.55rem;
	}
	.delete-box {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
		padding: 1rem;
		background: var(--bg-subtle);
		border: 1px solid var(--border-color);
		border-radius: var(--radius-md);
	}
	.delete-box h3 {
		font-size: 1rem;
	}
	.delete-box p {
		margin-top: 0.2rem;
		font-size: 0.85rem;
		color: var(--text-muted);
	}
	.danger-button {
		flex: 0 0 auto;
		border-radius: var(--radius-md);
		background: var(--error-solid);
		color: var(--error-on-solid);
		font-weight: 700;
		padding: 0.6rem 0.9rem;
	}

	.danger-button:hover:not(:disabled) {
		background: var(--error-solid-hover);
	}
	.danger-button:disabled {
		cursor: wait;
		opacity: 0.65;
	}
	.status {
		color: var(--success-text);
		font-weight: 600;
	}
	.status.error {
		color: var(--error-text);
	}
	@media (max-width: 640px) {
		.content-card {
			padding: 1.25rem;
		}
		.delete-box {
			align-items: stretch;
			flex-direction: column;
		}
	}
</style>
