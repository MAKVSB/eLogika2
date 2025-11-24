import type { ListJoinedStudentsResponse } from '$lib/api_types.js';
import { DataTableSearchParams } from '$lib/api_types_static';
import { API } from '$lib/services/api.svelte.js';
import { tableConfig } from './schema';

export async function load({ fetch, params, url }) {
	const search =
		url.searchParams.get(tableConfig.searchParam) ??
		DataTableSearchParams.fromDataTable(tableConfig.initialState).toURL();
	const showHistory = url.searchParams.get('showHistory');

	const data = API.request<null, ListJoinedStudentsResponse>(
		`/api/v2/courses/${params.courseId}/items/${params.itemId}/terms/${params.termId}/students`,
		{
			searchParams: {
				...(search ? { search } : {}),
				...(showHistory ? { showHistory } : {})
			}
		},
		fetch
	);

	return {
		students: data,
		showHistory: showHistory != null
	};
}
