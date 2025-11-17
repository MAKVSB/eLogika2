import type { ListJoinedStudentsResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params, url }) {
	const search = url.searchParams.get('search');

	const data = API.request<null, ListJoinedStudentsResponse>(
		`/api/v2/courses/${params.courseId}/items/${params.itemId}/terms/${params.termId}/students`,
		{
			searchParams: {
				...(search ? { search: search } : {})
			}
		},
		fetch
	);

	return {
		students: data
	};
}
