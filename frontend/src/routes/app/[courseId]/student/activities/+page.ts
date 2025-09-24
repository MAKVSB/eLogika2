import type { ListAvailableActivitiesResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, params }) {
	const data = API.request<null, ListAvailableActivitiesResponse>(
		`/api/v2/courses/${params.courseId}/activities/available`,
		{},
		fetch
	);

	return {
		tests: data
	};
}
