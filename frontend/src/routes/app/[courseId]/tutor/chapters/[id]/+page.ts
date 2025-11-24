import type { ChapterGetByIdResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params, url }) {
	if (params.id === '0') {
		return {
			creating: true
		};
	}

	const data = API.request<null, ChapterGetByIdResponse>(
		`/api/v2/courses/${params.courseId}/chapters/${params.id}`,
		{},
		fetch
	);

	return {
		creating: false,
		course: data
	};
}
