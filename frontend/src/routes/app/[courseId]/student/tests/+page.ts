import type { ListAvailableTestsResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	const data = API.request<null, ListAvailableTestsResponse>(
		`/api/v2/courses/${params.courseId}/tests/available`,
		{},
		fetch
	);

	return {
		tests: data
	};
}
