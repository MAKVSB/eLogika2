import type { ClassListResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params, url }) {
	const search = url.searchParams.get('class_search');

	const data = API.request<null, ClassListResponse>(
		`/api/v2/courses/${params.courseId}/classes`,
		{
			searchParams: {
				...(search ? { search: search } : {})
			}
		},
		fetch
	);

	return { data };
}
