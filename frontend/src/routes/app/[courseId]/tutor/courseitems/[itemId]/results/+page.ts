import type { CourseItemListResultsResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	const data = API.request<null, CourseItemListResultsResponse>(
		`/api/v2/courses/${params.courseId}/items/${params.itemId}/results`,
		{},
		fetch
	);

	return {
		results: data
	};
}
