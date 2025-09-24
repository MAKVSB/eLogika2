import GlobalState from '$lib/shared.svelte';
import { toast } from 'svelte-sonner';

type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';

interface RequestOptions<T> {
	method?: HttpMethod;
	headers?: Record<string, string>;
	body?: T;
	retry?: boolean;
	credentials?: RequestCredentials;
	searchParams?: Record<string, string>;
	keepalive?: boolean;
}

export class ApiError extends Error {
	public data: any;

	constructor(message: string, data: any) {
		super(message);
		this.data = data; // this property is defined in parent
	}
}

class Api {
	constructor() {}

	baseUrl = import.meta.env.VITE_API_URL;

	public async login(
		formData: { email: string; password: string },
		fetchFunc?: any
	): Promise<boolean> {
		let f = fetch;
		if (fetchFunc) {
			f = fetchFunc;
		}

		try {
			const res = await f(this.baseUrl + '/api/v2/auth/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(formData),
				credentials: 'include'
			});
			if (!res.ok) throw new Error('Failed to login');
			const data = await res.json();
			GlobalState.accessToken = data.accessToken;

			return true;
		} catch (err) {
			console.error(err);
			GlobalState.accessToken = null;
			GlobalState.loggedUser = null;
			return false;
		}
	}

	public async refreshAccessToken(fetchFunc?: any): Promise<boolean> {
		let f = fetch;
		if (fetchFunc) {
			f = fetchFunc;
		}
		try {
			const res = await f(this.baseUrl + '/api/v2/auth/refresh', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					...(GlobalState.accessToken ? { Authorization: `Bearer ${GlobalState.accessToken}` } : {})
				},
				body: JSON.stringify({}),
				credentials: 'include'
			});
			if (!res.ok) throw new Error('Failed to refresh token');
			const data = await res.json();
			GlobalState.accessToken = data.accessToken;
			return true;
		} catch (err) {
			GlobalState.accessToken = null;
			GlobalState.loggedUser = null;
			return false;
		}
	}

	public async request<T = any, U = any>(
		url: string,
		options: RequestOptions<T> = {},
		fetchFunc?: any
	): Promise<U> {
		let f = fetch;
		if (fetchFunc) {
			f = fetchFunc;
		}

		console.log(options);
		console.log(JSON.stringify(options.body));

		const fetchUrl = new URL(url, this.baseUrl);
		fetchUrl.search = new URLSearchParams(options.searchParams).toString();

		const config: RequestInit = {
			method: options.method || 'GET',
			headers: {
				...(options.body instanceof FormData ? {} : { 'Content-Type': 'application/json' }),
				...(GlobalState.accessToken ? { Authorization: `Bearer ${GlobalState.accessToken}` } : {}),
				...(options.headers || {}),
				'X-AS-ROLE': (GlobalState.activeRole ?? 'STUDENT') as unknown as string
			},
			body:
				options.body instanceof FormData
					? options.body
					: options.body
						? JSON.stringify(options.body)
						: undefined,
			...(options.credentials ? { credentials: options.credentials } : {}),
			keepalive: options.keepalive
		};

		console.log(config);

		let res = await f(fetchUrl, config);

		if (res.status === 401 && options.retry !== false) {
			const refreshed = await this.refreshAccessToken();
			if (refreshed) {
				return this.request<T, U>(url, { ...options, retry: false });
			}
			throw new Error('Unauthorized: Token refresh failed');
		}
		if (!res.ok) {
			const errResBody = await res.json();
			if ('message' in errResBody) {
				toast.error(errResBody.message, {
					description: errResBody.details
				});
				throw new ApiError(`API error`, errResBody);
			} else {
				throw new ApiError(`API error`, errResBody);
			}
		}
		const contentType = res.headers.get('Content-Type') || '';

		if (contentType && contentType.includes('application/json')) {
			return res.json() as unknown as U;
		}
		console.log(contentType);
		if (contentType.startsWith('image/') || contentType.includes('application/')) {
			return (await res.blob()) as unknown as U;
		}
		return res.text() as unknown as U;
	}
}

export const API = new Api();

export function encodeJsonToBase64Url<T>(object: T): string {
	const jsonString = JSON.stringify(object);
	const utf8Bytes = new TextEncoder().encode(jsonString);
	const base64 = btoa(String.fromCharCode(...utf8Bytes)); // Correctly encode UTF-8
	return base64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, ''); // make it URL-safe
}

export function decodeBase64UrlToJson<T>(searchParamString: string): T {
	const urlSafe = searchParamString.replace(/-/g, '+').replace(/_/g, '/');

	let base64 = urlSafe.replace(/-/g, '+').replace(/_/g, '/');
	while (base64.length % 4 !== 0) {
		base64 += '=';
	}
	const binaryString = atob(base64);
	const utf8Bytes = Uint8Array.from(binaryString, (c) => c.charCodeAt(0));
	const jsonString = new TextDecoder().decode(utf8Bytes);

	return JSON.parse(jsonString);
}
