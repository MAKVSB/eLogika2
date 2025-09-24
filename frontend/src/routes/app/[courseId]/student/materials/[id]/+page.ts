import type { ChapterGetByIdResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	const data = API.request<null, ChapterGetByIdResponse>(
		`/api/v2/courses/${params.courseId}/chapters/${params.id}`,
		{},
		fetch
	);

	return {
		chapter: data
	};
}
