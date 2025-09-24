import type { TestInstanceListResponse, TestListResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	const testData = API.request<null, TestInstanceListResponse>(
		`/api/v2/courses/${params.courseId}/tests/${params.itemId}/instances/${params.testId}`,
		{},
		fetch
	);

	return {
		tests: testData
	};
}
