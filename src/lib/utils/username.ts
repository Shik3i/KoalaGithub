export const DEFAULT_USERNAME = 'Shik3i';

const GITHUB_USERNAME_REGEX = /^[a-zA-Z0-9](?:[a-zA-Z0-9]|-(?=[a-zA-Z0-9])){0,38}$/;

export function isValidGitHubUsername(username: string): boolean {
	return GITHUB_USERNAME_REGEX.test(username.trim());
}

export function sanitizeUsername(username: string): string {
	const trimmed = username.trim();
	if (!trimmed || !isValidGitHubUsername(trimmed)) {
		return DEFAULT_USERNAME;
	}
	return trimmed;
}

export function getSanitizedAndEncodedUsername(username: string): string {
	return encodeURIComponent(sanitizeUsername(username));
}

export function updateUsernameUrl(username: string): void {
	if (typeof window === 'undefined') return;
	const url = new URL(window.location.href);
	const cleanUsername = sanitizeUsername(username);

	if (cleanUsername === DEFAULT_USERNAME) {
		url.searchParams.delete('user');
	} else {
		url.searchParams.set('user', cleanUsername);
	}

	window.history.replaceState({}, '', url.toString());
}
