import type { CourseItemGetByIdResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	if (params.itemId === '0') {
		return {
			creating: true
		};
	}

	const data = API.request<null, CourseItemGetByIdResponse>(
		`/api/v2/courses/${params.courseId}/items/${params.itemId}`,
		{},
		fetch
	);

	return {
		creating: false,
		courseItem: data
	};
}
