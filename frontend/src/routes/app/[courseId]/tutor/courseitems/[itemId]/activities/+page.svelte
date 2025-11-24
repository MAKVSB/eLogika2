<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { filters, tableConfig } from './schema.svelte';
	import { API, decodeBase64UrlToJson } from '$lib/services/api.svelte';
	import type { ActivityListItemDTO } from '$lib/api_types';
	import { type TableState } from '@tanstack/table-core';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import type { SelectOptions } from '$lib/components/ui/form';
	import { FilterTypeEnum } from '$lib/components/ui/data-table/filter';
	import { m } from '$lib/paraglide/messages';
	import { invalidateAll } from '$app/navigation';
	import { DataTableSearchParams } from '$lib/api_types_static';

	let rowItems: ActivityListItemDTO[] = $state([]);
	let rowCount: number = $state(0);

	let termIdFilter: number | undefined = $state(undefined);

	let { data } = $props();

	const actionsColumn = tableConfig.columns.find((c) => c.id == 'actions');
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
			})
			.catch(() => {});
	});

	$effect(() => {
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
			})
			.catch(() => {});
	});

	$effect(() => {
		const search =
			page.url.searchParams.get(tableConfig.searchParam) ??
			DataTableSearchParams.fromDataTable(tableConfig.initialState).toURL();
		const state = decodeBase64UrlToJson(search) as TableState;
		const foundTermIdFilter = (state.columnFilters ?? []).find((f) => f.id == 'termId');
		if (foundTermIdFilter) {
			termIdFilter = foundTermIdFilter.value as number;
		} else {
			termIdFilter = undefined;
		}
	});
</script>

<div class="flex flex-col gap-8 m-8">
	<div class="flex flex-row justify-between">
		<h1 class="text-2xl">Submitted activities and results</h1>
	</div>
	<DataTable data={rowItems} {rowCount} {...tableConfig} />
	<div class="flex justify-end gap-4"></div>
</div>
