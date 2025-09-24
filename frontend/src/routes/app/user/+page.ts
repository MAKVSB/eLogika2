import type { UserGetByIdResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	const data = API.request<null, UserGetByIdResponse>(`/api/v2/users/self`, {}, fetch);

	return {
		userData: data
	};
}
