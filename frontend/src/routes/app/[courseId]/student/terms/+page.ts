import type { StudentTermsListResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	const data = API.request<null, StudentTermsListResponse>(
		`/api/v2/courses/${params.courseId}/terms`,
		{},
		fetch
	);

	return {
		terms: data
	};
}
