<script lang="ts">
	import { onMount } from 'svelte';

	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters, initialState, searchParam, tableConfig } from './schema';
	import { API } from '$lib/services/api.svelte';
	import type {
		PrintTestRequest,
		TestEvaluationRequest,
		TestEvaluationResponse,
		TestListItemDTO
	} from '$lib/api_types';
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
	import { Label } from '$lib/components/ui/label';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import Loader from '$lib/components/ui/loader/loader.svelte';

	let rowItems: TestListItemDTO[] = $state([]);
	let rowCount: number = $state(0);

	let termIdFilter: number | undefined = $state(undefined);

	let printRunning = $state(false);

	let { data } = $props();

	let dialogOpen = $state(false);
	let printAnswerSheets = $state(true);
	let separateAnswerSheets = $state(false);

	const actionsColumn = columns.find((c) => c.id == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (event: string, id: number, otherParams: any) => {
				switch (event) {
					case 'print':
						API.request<PrintTestRequest, Blob>(
							`/api/v2/courses/${page.params.courseId}/print/tests`,
							{
								method: 'POST',
								body: {
									courseItemId: Number(page.params.itemId),
									printAnswerSheets: printAnswerSheets,
									separateAnswerSheets: separateAnswerSheets,
									testId: Number(id)
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
					case 'reevaluate':
						API.request<TestEvaluationRequest, TestEvaluationResponse>(
							`/api/v2/courses/${page.params.courseId}/tests/evaluate`,
							{
								method: 'POST',
								body: {
									courseItemId: Number(page.params.itemId),
									testId: Number(id)
								}
							},
							fetch
						)
							.then((res) => {
								if (res.success) {
									toast.success('Test re-evaluated successfuly');
								}
								invalidateAll();
							})
							.catch((err) => {
								console.log(err);
							});
					case 'delete':
						if (!confirm(m.test_delete_confirm())) {
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

				const existingFilter = filters.find((f) => f.accessorKey == 'termId');
				if (existingFilter) {
					if (existingFilter.type == FilterTypeEnum.SELECT) {
						existingFilter.values = termOptions;
					}
				}
			})
			.catch(() => {});
	});

	function refetch(state: TableState) {
		const foundTermIdFilter = state.columnFilters.find((f) => f.id == 'termId');
		if (foundTermIdFilter) {
			termIdFilter = foundTermIdFilter.value as number;
		} else {
			termIdFilter = undefined;
		}
	}

	function print() {
		printRunning = true;
		API.request<PrintTestRequest, Blob>(
			`/api/v2/courses/${page.params.courseId}/print/tests`,
			{
				method: 'POST',
				body: {
					termId: Number(termIdFilter),
					courseItemId: Number(page.params.itemId),
					printAnswerSheets: printAnswerSheets,
					separateAnswerSheets: separateAnswerSheets
				}
			},
			fetch
		)
			.then((res) => {
				const url = URL.createObjectURL(res);
				window.open(url); // opens in new tab
			})
			.catch(() => {})
			.finally(() => {
				printRunning = false;
			});
	}

	function reevaluate() {
		API.request<TestEvaluationRequest, TestEvaluationResponse>(
			`/api/v2/courses/${page.params.courseId}/tests/evaluate`,
			{
				method: 'POST',
				body: {
					termId: Number(termIdFilter),
					courseItemId: Number(page.params.itemId)
				}
			},
			fetch
		)
			.then((res) => {
				if (res.success) {
					toast.success('Test re-evaluated successfuly');
				}
				invalidateAll();
			})
			.catch((err) => {
				console.log(err);
			});
	}
</script>

<div class="flex flex-col gap-8 m-8">
	<div class="flex flex-row justify-between">
		<h1 class="text-2xl">Generated tests management</h1>
	</div>
	<DataTable data={rowItems} {rowCount} {refetch} {...tableConfig} />
	<div class="flex justify-end gap-4">
		<Dialog.Root bind:open={dialogOpen}>
			<Dialog.Trigger class={buttonVariants({ variant: 'default' })} disabled={!termIdFilter}>
				{m.tests_generate()}
			</Dialog.Trigger>
			{#if dialogOpen && termIdFilter}
				<GeneratorDialog bind:openState={dialogOpen} termId={termIdFilter}></GeneratorDialog>
			{/if}
		</Dialog.Root>
		<Button onclick={() => reevaluate()}>{m.test_reevaluate({ type: 'multi' })}</Button>
		<Button onclick={() => print()} disabled={!termIdFilter || printRunning}>
			{#if printRunning}
				<Loader></Loader>
			{/if}
			{m.print_test({ type: 'multi' })}
		</Button>
	</div>
	<div class="flex flex-col items-end gap-4">
		<h3 class="text-xl">{m.test_print_settings()}</h3>
		<div class="flex gap-2">
			<Label for="printAnswerSheets">{m.print_test_printanswersheets()}</Label>
			<Checkbox
				class="rounded-md h-9 w-9"
				id="printAnswerSheets"
				bind:checked={printAnswerSheets}
			/>
		</div>
		<div class="flex gap-2">
			<Label for="separateAnswerSheets">{m.print_test_separateanswersheets()}</Label>
			<Checkbox
				class="rounded-md h-9 w-9"
				id="separateAnswerSheets"
				bind:checked={separateAnswerSheets}
			/>
		</div>
	</div>
</div>
