import type { QuestionGetByIdResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	if (params.id === '0') {
		return {
			creating: true
		};
	}

	const data = API.request<null, QuestionGetByIdResponse>(
		`/api/v2/courses/${params.courseId}/questions/${params.id}`,
		{},
		fetch
	);

	return {
		creating: false,
		question: data
	};
}
