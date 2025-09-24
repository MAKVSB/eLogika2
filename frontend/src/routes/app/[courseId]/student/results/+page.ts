import type { StudentCourseItemListResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	const data = API.request<null, StudentCourseItemListResponse>(
		`/api/v2/courses/${params.courseId}/items/students`,
		{},
		fetch
	);

	return {
		courseItems: data
	};
}
