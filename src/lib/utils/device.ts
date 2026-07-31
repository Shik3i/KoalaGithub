export function getOrCreateDeviceId(): string {
	if (typeof window === 'undefined') return 'server-side';

	const STORAGE_KEY = 'koala_device_id';
	let deviceId = localStorage.getItem(STORAGE_KEY);

	if (!deviceId) {
		if (typeof crypto !== 'undefined' && crypto.randomUUID) {
			deviceId = crypto.randomUUID();
		} else {
			deviceId = 'dev_' + Math.random().toString(36).substring(2, 15) + Date.now().toString(36);
		}
		localStorage.setItem(STORAGE_KEY, deviceId);
	}

	return deviceId;
}
