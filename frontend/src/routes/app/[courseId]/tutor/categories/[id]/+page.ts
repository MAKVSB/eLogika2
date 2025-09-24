import type { CategoryGetByIdResponse } from '$lib/api_types';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	if (params.id === '0') {
		return {
			creating: true
		};
	}

	const data = API.request<null, CategoryGetByIdResponse>(
		`/api/v2/courses/${params.courseId}/categories/${params.id}`,
		{},
		fetch
	);

	return {
		creating: false,
		category: data
	};
}
