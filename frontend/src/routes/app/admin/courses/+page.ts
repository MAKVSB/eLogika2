import type { CourseListRequest, CourseListResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, url }) {
	const search = url.searchParams.get('search');

	const data = API.request<CourseListRequest, CourseListResponse>(
		`/api/v2/courses`,
		{
			searchParams: {
				...(search ? { search: search } : {})
			}
		},
		fetch
	);

	return {
		courses: data
	};
}
