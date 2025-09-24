import type { ClassGetByIdResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	if (params.classId === '0') {
		return {
			creating: true
		};
	}

	const data = API.request<null, ClassGetByIdResponse>(
		`/api/v2/courses/${params.courseId}/classes/${params.classId}`,
		{},
		fetch
	);

	return {
		creating: false,
		courseItem: data
	};
}
