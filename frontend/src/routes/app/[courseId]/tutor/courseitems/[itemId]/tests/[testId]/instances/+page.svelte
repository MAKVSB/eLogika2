<script lang="ts">
	import { onMount } from 'svelte';

	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import { API, ApiError, decodeBase64UrlToJson } from '$lib/services/api.svelte';
	import type {
		PrintTestRequest,
		TestEvaluationRequest,
		TestEvaluationResponse,
		TestInstanceListItemDTO,
	} from '$lib/api_types';
	import { type InitialTableState } from '@tanstack/table-core';
	import { page } from '$app/state';
	import Button, { buttonVariants } from '$lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import CreateInstanceDialog from './CreateInstanceDialog/CreateInstanceDialog.svelte';
	import { invalidateAll } from '$app/navigation';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Label } from '$lib/components/ui/label';
	import { m } from '$lib/paraglide/messages';

	let loading: boolean = $state(true);
	let rowItems: TestInstanceListItemDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});

	let { data } = $props();

	let dialogOpen = $state(false);
	let printAnswerSheets = $state(true);
	let separateAnswerSheets = $state(true);

	$effect(() => {
		data.tests
			.then((res) => {
				rowItems = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {});
	});

	onMount(() => {
		const encodedParams = page.url.searchParams.get('search');
		if (encodedParams) {
			initialState = decodeBase64UrlToJson(encodedParams);
		}
		loading = false;
	});

	const actionsColumn = columns.find((c) => c.uniqueId == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (event: string, id: number) => {
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
									testId: Number(page.params.testId),
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
					case 'reevaluate':
						API.request<TestEvaluationRequest, TestEvaluationResponse>(
							`/api/v2/courses/${page.params.courseId}/tests/evaluate`,
							{
								method: 'POST',
								body: {
									courseItemId: Number(page.params.itemId),
									testId: Number(page.params.testId),
									instanceId: id
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
						break;
					case 'delete':
						if (!confirm('Instance will be deleted permanently.')) {
							return;
						}
						API.request<any, Blob>(
							`/api/v2/courses/${page.params.courseId}/tests/${page.params.itemId}/instance/${id}`,
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

	function print() {
		API.request<PrintTestRequest, Blob>(
			`/api/v2/courses/${page.params.courseId}/print/tests`,
			{
				method: 'POST',
				body: {
					courseItemId: Number(page.params.itemId),
					printAnswerSheets: printAnswerSheets,
					separateAnswerSheets: separateAnswerSheets,
					testId: Number(page.params.testId)
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

	function reevaluate() {
		API.request<TestEvaluationRequest, TestEvaluationResponse>(
			`/api/v2/courses/${page.params.courseId}/tests/evaluate`,
			{
				method: 'POST',
				body: {
					testId: Number(page.params.testId),
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
		<h1 class="text-2xl">Generated tests management (instances)</h1>
		<Dialog.Root bind:open={dialogOpen}>
			<Dialog.Trigger class={buttonVariants({ variant: 'default' })}>
				{m.test_instance_create()}
			</Dialog.Trigger>
			{#if dialogOpen}
				<CreateInstanceDialog bind:openState={dialogOpen}></CreateInstanceDialog>
			{/if}
		</Dialog.Root>
	</div>
	{#if !loading}
		<DataTable data={rowItems} {columns} {filters} {initialState} {rowCount} queryParam="search" />
	{/if}
	<div class="flex justify-end gap-4">
		<Button onclick={() => reevaluate()}>{m.test_reevaluate({ type: 'multi' })}</Button>
		<Button onclick={() => print()}>{m.print_test({ type: 'multi' })}</Button>
	</div>
	<div class="flex flex-col gap-4">
		<h3>{m.test_print_settings()}</h3>
		<div class="flex gap-2">
			<Checkbox
				class="rounded-md h-9 w-9"
				id="printAnswerSheets"
				bind:checked={printAnswerSheets}
			/>
			<Label for="printAnswerSheets">{m.print_test_printanswersheets()}</Label>
		</div>
		<div class="flex gap-2">
			<Checkbox
				class="rounded-md h-9 w-9"
				id="separateAnswerSheets"
				bind:checked={separateAnswerSheets}
			/>
			<Label for="separateAnswerSheets">{m.print_test_separateanswersheets()}</Label>
		</div>
	</div>
</div>
