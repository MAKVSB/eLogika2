import type { TermsGetByIdResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	if (params.termId === '0') {
		return {
			creating: true
		};
	}

	const data = API.request<null, TermsGetByIdResponse>(
		`/api/v2/courses/${params.courseId}/items/${params.itemId}/terms/${params.termId}`,
		{},
		fetch
	);

	return {
		creating: false,
		term: data
	};
}
