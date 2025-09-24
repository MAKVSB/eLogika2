import type { TemplateGetByIdResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	if (params.id === '0') {
		return {
			creating: true
		};
	}

	const data = API.request<null, TemplateGetByIdResponse>(
		`/api/v2/courses/${params.courseId}/templates/${params.id}`,
		{},
		fetch
	);

	return {
		creating: false,
		template: data
	};
}
