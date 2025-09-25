import type { UserListRequest, UserListResponse } from '$lib/api_types.js';
import { API } from '$lib/services/api.svelte.js';

export async function load({ fetch, url }) {
	const search = url.searchParams.get('search');

	const data = API.request<UserListRequest, UserListResponse>(
		'/api/v2/users',
		{
			searchParams: {
				...(search ? { search: search } : {})
			}
		},
		fetch
	);

	return {
		users: data
	};
}
