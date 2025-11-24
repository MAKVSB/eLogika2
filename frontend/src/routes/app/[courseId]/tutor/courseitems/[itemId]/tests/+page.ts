import type { TermsListRecursiveResponse, TestListResponse } from '$lib/api_types.js';
import { DataTableSearchParams } from '$lib/api_types_static';
import { API } from '$lib/services/api.svelte.js';
import { tableConfig } from './schema';

export async function load({ fetch, params, url }) {
	const search =
		url.searchParams.get(tableConfig.searchParam) ??
		DataTableSearchParams.fromDataTable(tableConfig.initialState).toURL();

	const testData = API.request<null, TestListResponse>(
		`/api/v2/courses/${params.courseId}/tests/${params.itemId}`,
		{
			searchParams: {
				search
			}
		},
		fetch
	);

	const termData = API.request<null, TermsListRecursiveResponse>(
		`/api/v2/courses/${params.courseId}/items/${params.itemId}/terms/recursive`,
		{},
		fetch
	);

	return {
		tests: testData,
		terms: termData
	};
}
