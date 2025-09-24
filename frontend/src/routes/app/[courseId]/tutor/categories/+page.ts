import type { CategoryListResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params, url }) {
	const search = url.searchParams.get('search');

	const data = API.request<null, CategoryListResponse>(
		`/api/v2/courses/${params.courseId}/categories`,
		{
			searchParams: {
				...(search ? { search: search } : {})
			}
		},
		fetch
	);

	return { data };
}
