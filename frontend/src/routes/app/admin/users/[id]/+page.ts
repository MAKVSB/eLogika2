import type { UserGetByIdResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	if (params.id === '0') {
		return {
			creating: true
		};
	}

	const data = API.request<null, UserGetByIdResponse>(`/api/v2/users/${params.id}`, {}, fetch);

	return {
		creating: false,
		question: data
	};
}
