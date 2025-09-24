import type { TermsListRecursiveResponse, ActivityListResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params, url }) {
	const search = url.searchParams.get('search');

	const testData = API.request<null, ActivityListResponse>(
		`/api/v2/courses/${params.courseId}/activities/${params.itemId}`,
		{
			searchParams: {
				...(search ? { search: search } : {})
			}
		},
		fetch
	);

	const termData = API.request<null, TermsListRecursiveResponse>(
		`/api/v2/courses/${params.courseId}/items/${params.itemId}/terms/recursive`,
		{},
		fetch
	);

	return {
		tests: testData,
		terms: termData
	};
}
