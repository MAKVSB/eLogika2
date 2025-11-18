import type { ListJoinedStudentsResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params, url }) {
	const search = url.searchParams.get('search');
	const showHistory = url.searchParams.get('showHistory');

	const data = API.request<null, ListJoinedStudentsResponse>(
		`/api/v2/courses/${params.courseId}/items/${params.itemId}/terms/${params.termId}/students`,
		{
			searchParams: {
				...(search ? { search } : {}),
				...(showHistory ? { showHistory } : {})
			}
		},
		fetch
	);

	return {
		students: data,
		showHistory: showHistory != null
	};
}
