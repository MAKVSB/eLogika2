import type { QuestionListResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params, url }) {
	const search = url.searchParams.get('search');

	const data = API.request<null, QuestionListResponse>(
		`/api/v2/courses/${params.courseId}/questions`,
		{
			searchParams: {
				...(search ? { search: search } : {})
			}
		},
		fetch
	);

	return { data };
}
