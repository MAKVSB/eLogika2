<script lang="ts">
	import { onMount } from 'svelte';

	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import { API, ApiError, decodeBase64UrlToJson } from '$lib/services/api.svelte';
	import type { PrintTestRequest, TestListItemDTO } from '$lib/api_types';
	import { type InitialTableState, type TableState } from '@tanstack/table-core';
	import { page } from '$app/state';
	import Button, { buttonVariants } from '$lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import type { SelectOptions } from '$lib/components/ui/form';
	import { FilterTypeEnum } from '$lib/components/ui/data-table/filter';
	import { m } from '$lib/paraglide/messages';
	import GeneratorDialog from './GeneratorDialog/GeneratorDialog.svelte';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { invalidateAll } from '$app/navigation';

	let loading: boolean[] = $state([true, true, true]);
	let rowItems: TestListItemDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({
		pagination: {
			pageIndex: 0,
			pageSize: 25
		}
	});

	let termIdFilter: number | undefined = $state(undefined);

	let { data } = $props();

	let dialogOpen = $state(false);

	const actionsColumn = columns.find((c) => c.uniqueId == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (event: string, id: number, otherParams: any) => {
				switch (event) {
					case 'print':
						console.log(otherParams)
						API.request<PrintTestRequest, Blob>(
							`/api/v2/courses/${page.params.courseId}/print/tests`,
							{
								method: 'POST',
								body: {
									courseItemId: Number(page.params.itemId),
									printInstances: otherParams.instances ?? false,
									testId: Number(id),
									instanceId: id
								}
							},
							fetch
						)
							.then((res) => {
								const url = URL.createObjectURL(res);
								window.open(url); // opens in new tab
							})
					     .catch(() => {});
						break;
					case 'delete':
						if (!confirm('Test and all its instances will be deleted permanently.')) {
							return;
						}
						API.request<any, Blob>(
							`/api/v2/courses/${page.params.courseId}/tests/${page.params.itemId}/instances/${id}`,
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
				loading[2] = true;
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

	function print(instances: boolean = false) {
		API.request<PrintTestRequest, Blob>(
			`/api/v2/courses/${page.params.courseId}/print/tests`,
			{
				method: 'POST',
				body: {
					termId: termIdFilter,
					courseItemId: Number(page.params.itemId),
					printInstances: instances
				}
			},
			fetch
		)
			.then((res) => {
				const url = URL.createObjectURL(res);
				window.open(url); // opens in new tab
			})
			.catch(() => {});
	}
</script>

<div class="flex flex-col gap-8 m-8">
	<div class="flex flex-row justify-between">
		<h1 class="text-2xl">Generated tests management</h1>
	</div>
	{#if !loading[0] && !loading[1] && !loading[2]}
		<DataTable data={rowItems} {columns} {filters} {initialState} {rowCount} {refetch} queryParam='search'/>
	{/if}
	<div class="flex justify-end gap-4">
		<Dialog.Root bind:open={dialogOpen}>
			<Dialog.Trigger class={buttonVariants({ variant: 'default' })} disabled={!termIdFilter}>
				Generate tests
			</Dialog.Trigger>
			{#if dialogOpen && termIdFilter}
				<GeneratorDialog bind:openState={dialogOpen} termId={termIdFilter}></GeneratorDialog>
			{/if}
		</Dialog.Root>
		<Button onclick={() => print(true)} disabled={!termIdFilter}>Print all instances</Button>
		<Button onclick={() => print(false)} disabled={!termIdFilter}>Print all tests</Button>
	</div>
</div>
