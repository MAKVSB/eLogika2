import type { TestInstanceListResponse, TestListResponse } from '$lib/api_types.js';
import { DataTableSearchParams } from '$lib/api_types_static';
import { API } from '$lib/services/api.svelte.js';
import { tableConfig } from './schema';

export async function load({ fetch, params, url }) {
	const search =
		url.searchParams.get(tableConfig.searchParam) ??
		DataTableSearchParams.fromDataTable(tableConfig.initialState).toURL();

	const testData = API.request<null, TestInstanceListResponse>(
		`/api/v2/courses/${params.courseId}/tests/${params.itemId}/instances/${params.testId}`,
		{
			searchParams: {
				search
			}
		},
		fetch
	);

	return {
		tests: testData
	};
}
