<script lang="ts">
	import { onMount } from 'svelte';

	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import {
		API,
		ApiError,
		decodeBase64UrlToJson,
		encodeJsonToBase64Url
	} from '$lib/services/api.svelte';
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
	import GlobalState from '$lib/shared.svelte';
	import { toast } from 'svelte-sonner';
	import type { CategoryListItemDTO, CategoryListResponse } from '$lib/api_types';
	import { m } from '$lib/paraglide/messages';
	import { base } from '$app/paths';
	import { DataTableSearchParams } from '$lib/api_types_static';

	let loading: boolean = $state(true);
	let data: CategoryListItemDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});
	let encodedParams: string | null = $state(null);

	$effect(() => {
		if (!page.params.courseId) return;
		if (encodedParams) {
			console.log("Transfering 21")
			goto(`?search=${encodedParams}`);
			fetchData(encodedParams);
		} else {
			fetchData();
		}
	});

	async function fetchData(encodedFilters?: string) {
		await API.request<null, CategoryListResponse>(
			`/api/v2/courses/${page.params.courseId}/chapters/${page.params.id}/categories`,
			{
				searchParams: {
					...(encodedFilters ? { search: encodedFilters } : {})
				}
			}
		)
			.then((res) => {
				data = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {});
	}

	function refetch(state: TableState) {
		encodedParams = DataTableSearchParams.fromDataTable(state).toURL();
	}

	onMount(async () => {
		const encodedParams = page.url.searchParams.get('search');
		if (encodedParams) {
			initialState = decodeBase64UrlToJson(encodedParams);
		}
		loading = false;
	});
</script>

<div class="flex flex-row justify-between">
	<h1 class="mb-8 text-2xl">Categories</h1>
	<div class="flex gap-2">
		<Button href="{base}/app/{page.params.courseId}/tutor/categories/0">{m.category_add()}</Button>
	</div>
</div>
{#if !loading}
	<DataTable {data} {columns} {filters} {refetch} {initialState} {rowCount} queryParam='search'/>
{/if}
