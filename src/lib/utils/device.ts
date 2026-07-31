const STORAGE_KEY = 'koala_device_id';
const UUID_V4 = /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

export function getStoredDeviceId(): string | null {
	if (typeof window === 'undefined') return null;
	try {
		const value = localStorage.getItem(STORAGE_KEY);
		return value && UUID_V4.test(value) ? value.toLowerCase() : null;
	} catch {
		return null;
	}
}

export function getOrCreateDeviceId(): string | null {
	const stored = getStoredDeviceId();
	if (stored) return stored;
	if (typeof window === 'undefined' || typeof crypto === 'undefined') return null;

	let deviceId: string;
	if (typeof crypto.randomUUID === 'function') {
		deviceId = crypto.randomUUID();
	} else if (typeof crypto.getRandomValues === 'function') {
		const bytes = crypto.getRandomValues(new Uint8Array(16));
		bytes[6] = (bytes[6] & 0x0f) | 0x40;
		bytes[8] = (bytes[8] & 0x3f) | 0x80;
		const hex = Array.from(bytes, (byte) => byte.toString(16).padStart(2, '0')).join('');
		deviceId = `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20)}`;
	} else {
		return null;
	}

	try {
		localStorage.setItem(STORAGE_KEY, deviceId);
		return deviceId;
	} catch {
		return null;
	}
}

export function clearDeviceId(): void {
	if (typeof window === 'undefined') return;
	try {
		localStorage.removeItem(STORAGE_KEY);
	} catch {
		// Storage may be unavailable in hardened browser contexts.
	}
}
