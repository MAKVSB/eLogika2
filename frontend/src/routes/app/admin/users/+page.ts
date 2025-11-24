import type { UserListRequest, UserListResponse } from '$lib/api_types.js';
import { DataTableSearchParams } from '$lib/api_types_static';
import { API } from '$lib/services/api.svelte.js';
import { tableConfig } from './schema';

export async function load({ fetch, url }) {
	const search =
		url.searchParams.get(tableConfig.searchParam) ??
		DataTableSearchParams.fromDataTable(tableConfig.initialState).toURL();

	const data = API.request<UserListRequest, UserListResponse>(
		'/api/v2/users',
		{
			searchParams: {
				search
			}
		},
		fetch
	);

	return {
		users: data
	};
}
