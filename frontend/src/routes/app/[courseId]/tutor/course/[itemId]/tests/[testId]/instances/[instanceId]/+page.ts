import type { TestInstanceTutorGetResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	const test = API.request<null, TestInstanceTutorGetResponse>(
		`/api/v2/courses/${params.courseId}/tests/${params.itemId}/instance/${params.instanceId}`,
		{},
		fetch
	);

	return {
		test
	};
}
