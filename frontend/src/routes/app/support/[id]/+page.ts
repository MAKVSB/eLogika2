import type { SupportTicketGetByIdResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	if (params.id === '0') {
		return {
			creating: true
		};
	}

	const data = API.request<null, SupportTicketGetByIdResponse>(
		`/api/v2/support/${params.id}`,
		{},
		fetch
	);

	return {
		creating: false,
		ticket: data
	};
}
