import type { CourseGetByIdResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	if (params.courseId === '0') {
		return {
			creating: true
		};
	}

	const data = API.request<null, CourseGetByIdResponse>(
		`/api/v2/courses/${params.courseId}`,
		{},
		fetch
	);

	return {
		creating: false,
		course: data
	};
}
