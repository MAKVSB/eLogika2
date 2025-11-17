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
	import { goto, invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import type { TermDTO, TermsListRequest, TermsListResponse } from '$lib/api_types';
	import { base } from '$app/paths';
	import { m } from '$lib/paraglide/messages';
	import { DataTableSearchParams } from '$lib/api_types_static';

	let loading: boolean = $state(true);
	let data: TermDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});
	let encodedParams: string | null = $state(null);

	let {
		courseId,
		itemId
	}: {
		courseId: string;
		itemId: string;
	} = $props();

	const actionsColumn = columns.find((c) => c.uniqueId == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (event: string, id: number) => {
				switch (event) {
					case 'delete':
						if (!confirm(m.term_delete_confirm())) {
							return;
						}
						API.request<any, Blob>(
							`/api/v2/courses/${page.params.courseId}/items/${page.params.itemId}/terms/${id}`,
							{
								method: 'DELETE'
							},
							fetch
						)
							.then((res) => {
								invalidateAll();
							})
							.catch(() => {});
						break;
				}

				return true;
			}
		};
	}

	$effect(() => {
		if (encodedParams) {
			console.log("Transfering 27")
			goto(`?search=${encodedParams}`);
			fetchData(encodedParams);
		} else {
			fetchData();
		}
	});

	async function fetchData(encodedFilters?: string) {
		await API.request<TermsListRequest, TermsListResponse>(
			`/api/v2/courses/${courseId}/items/${itemId}/terms`,
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

<div>
	<div class="flex flex-row justify-between">
		<h1 class="mb-8 text-2xl">Terms</h1>
		<div class="flex gap-2">
			<Button href="{base}/app/{courseId}/tutor/courseitems/{itemId}/0">Add term</Button>
		</div>
	</div>
	{#if !loading}
		<DataTable {data} {columns} {filters} {refetch} {initialState} {rowCount} queryParam='search'/>
	{/if}
</div>
