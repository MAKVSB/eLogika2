import type { CategoryListResponse } from '$lib/api_types.js';
import { DataTableSearchParams } from '$lib/api_types_static.js';
import { API } from '$lib/services/api.svelte.js';
import { tableConfig } from './schema';

export async function load({ fetch, params, url }) {
	const search =
		url.searchParams.get(tableConfig.searchParam) ??
		DataTableSearchParams.fromDataTable(tableConfig.initialState).toURL();

	const data = API.request<null, CategoryListResponse>(
		`/api/v2/courses/${params.courseId}/categories`,
		{
			searchParams: {
				search
			}
		},
		fetch
	);

	return { data };
}
