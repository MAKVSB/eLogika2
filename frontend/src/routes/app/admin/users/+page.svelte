<script lang="ts">
	import { onMount } from 'svelte';

	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import { API, decodeBase64UrlToJson, encodeJsonToBase64Url } from '$lib/services/api.svelte';
	import {
		type ColumnFiltersState,
		type InitialTableState,
		type PaginationState,
		type SortingState,
		type TableState
	} from '@tanstack/table-core';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import type { UserListItemDTO, UserListRequest, UserListResponse } from '$lib/api_types';

	let loading: boolean = $state(true);
	let data: UserListItemDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});
	let encodedParams: string | null = $state(null);

	$effect(() => {
		if (encodedParams) {
			console.log("Transfering 33")
			goto(`?search=${encodedParams}`);
			fetchData(encodedParams);
		} else {
			fetchData();
		}
	});

	type RestRequest = {
		pagination?: PaginationState;
		sorting?: SortingState;
		columnFilters?: ColumnFiltersState;
	};

	async function fetchData(encodedFilters?: string) {
		await API.request<UserListRequest, UserListResponse>('/api/v2/users', {
			searchParams: {
				...(encodedFilters ? { search: encodedFilters } : {})
			}
		})
			.then((res) => {
				data = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {});
	}

	function refetch(state: TableState) {
		const queryParams: RestRequest = {
			...(state.pagination ? { pagination: state.pagination } : {}),
			...(state.sorting ? { sorting: state.sorting } : {}),
			...(state.columnFilters ? { columnFilters: state.columnFilters } : {})
		};
		encodedParams = encodeJsonToBase64Url(queryParams);
	}

	onMount(async () => {
		const encodedParams = page.url.searchParams.get('search');
		if (encodedParams) {
			initialState = decodeBase64UrlToJson(encodedParams);
		}
		loading = false;
	});
</script>

<div class="m-8">
	<div class="flex flex-row justify-between">
		<h1 class="mb-8 text-2xl">Admin user management</h1>
		<!-- TODO import users from edison-->
		<Button href="/app/admin/users/0">Add user</Button>
	</div>
	{#if !loading}
		<DataTable {data} {columns} {filters} {refetch} {initialState} {rowCount} queryParam='search'/>
	{/if}
</div>
