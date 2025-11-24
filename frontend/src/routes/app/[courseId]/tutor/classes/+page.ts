import type { ClassListResponse } from '$lib/api_types.js';
import { DataTableSearchParams } from '$lib/api_types_static';
import { API } from '$lib/services/api.svelte.js';
import { tableConfig } from './schema';

export async function load({ fetch, params, url }) {
	const search =
		url.searchParams.get(tableConfig.searchParam) ??
		DataTableSearchParams.fromDataTable(tableConfig.initialState).toURL();

	const data = API.request<null, ClassListResponse>(
		`/api/v2/courses/${params.courseId}/classes`,
		{
			searchParams: {
				search
			}
		},
		fetch
	);

	return { data };
}
