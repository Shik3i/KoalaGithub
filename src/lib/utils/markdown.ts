import { sanitizeUsername } from './username';

export function renderTemplate(template: string, username: string, theme: string): string {
	const cleanUsername = sanitizeUsername(username);
	const encodedUsername = encodeURIComponent(cleanUsername);
	const encodedTheme = encodeURIComponent(theme);

	return template
		.replaceAll('{username}', encodedUsername)
		.replaceAll('{theme}', encodedTheme);
}
