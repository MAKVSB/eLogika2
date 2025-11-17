<script lang="ts">
	import { onMount } from 'svelte';

	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import { API, ApiError, decodeBase64UrlToJson } from '$lib/services/api.svelte';
	import type { ActivityListItemDTO, PrintTestRequest } from '$lib/api_types';
	import { type InitialTableState, type TableState } from '@tanstack/table-core';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import type { SelectOptions } from '$lib/components/ui/form';
	import { FilterTypeEnum } from '$lib/components/ui/data-table/filter';
	import { m } from '$lib/paraglide/messages';
	import { invalidateAll } from '$app/navigation';

	let loading: boolean[] = $state([true, true, true]);
	let rowItems: ActivityListItemDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});

	let termIdFilter: number | undefined = $state(undefined);

	let { data } = $props();

	const actionsColumn = columns.find((c) => c.uniqueId == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (event: string, id: number) => {
				switch (event) {
					case 'delete':
						if (!confirm(m.activity_instance_delete_confirm())) {
							return;
						}
						API.request<any, Blob>(
							`/api/v2/courses/${page.params.courseId}/activities/instance/${id}`,
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
		data.tests
			.then((res) => {
				rowItems = res.items;
				rowCount = res.itemsCount;
				loading[2] = false;
			})
			.catch(() => {});

		data.terms
			.then((res) => {
				const termOptions: SelectOptions = res.items.map((term) => {
					return {
						value: term.id,
						display: term.name
					};
				});

				if (!filters.find((f) => f.accessorKey == 'termId')) {
					filters.push({
						type: FilterTypeEnum.SELECT,
						accessorKey: 'termId',
						values: termOptions,
						placeholder: m.filter_term(),
						emptyValue: 'No filter'
					});
				}

				loading[1] = false;
			})
			.catch(() => {});
	});

	onMount(() => {
		const encodedParams = page.url.searchParams.get('search');
		if (encodedParams) {
			initialState = decodeBase64UrlToJson(encodedParams);
			if (!initialState.columnVisibility) {
				initialState.columnVisibility = {
					termId: false
				};
			}
		}
		loading[0] = false;
	});

	function refetch(state: TableState) {
		const foundTermIdFilter = state.columnFilters.find((f) => f.id == 'termId');
		if (foundTermIdFilter) {
			termIdFilter = foundTermIdFilter.value as number;
		} else {
			termIdFilter = undefined;
		}
	}
</script>

<div class="flex flex-col gap-8 m-8">
	<div class="flex flex-row justify-between">
		<h1 class="text-2xl">Submitted activities and results</h1>
	</div>
	{#if !loading[0] && !loading[1] && !loading[2]}
		<DataTable data={rowItems} {columns} {filters} {initialState} {rowCount} {refetch} queryParam='search'/>
	{/if}
	<div class="flex justify-end gap-4"></div>
</div>
