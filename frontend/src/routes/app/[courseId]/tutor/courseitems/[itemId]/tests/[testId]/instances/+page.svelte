<script lang="ts">
	import { onMount } from 'svelte';

	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import { API, ApiError, decodeBase64UrlToJson } from '$lib/services/api.svelte';
	import type { PrintTestRequest, TestInstanceListItemDTO, TestListItemDTO } from '$lib/api_types';
	import { reSplitAlphaNumeric, type InitialTableState } from '@tanstack/table-core';
	import { page } from '$app/state';
	import Button, { buttonVariants } from '$lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import CreateInstanceDialog from './CreateInstanceDialog/CreateInstanceDialog.svelte';
	import { invalidateAll } from '$app/navigation';

	let loading: boolean = $state(true);
	let rowItems: TestInstanceListItemDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});

	let { data } = $props();

	let dialogOpen = $state(false);

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
									printInstances: true,
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

	function print(instances: boolean = false) {
		API.request<PrintTestRequest, Blob>(
			`/api/v2/courses/${page.params.courseId}/print/tests`,
			{
				method: 'POST',
				body: {
					courseItemId: Number(page.params.itemId),
					printInstances: instances,
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
</script>

<div class="flex flex-col gap-8 m-8">
	<div class="flex flex-row justify-between">
		<h1 class="text-2xl">Generated tests management (instances)</h1>
		<Dialog.Root bind:open={dialogOpen}>
			<Dialog.Trigger class={buttonVariants({ variant: 'default' })}>
				Create instance
			</Dialog.Trigger>
			{#if dialogOpen}
				<CreateInstanceDialog bind:openState={dialogOpen}></CreateInstanceDialog>
			{/if}
		</Dialog.Root>
	</div>
	{#if !loading}
		<DataTable data={rowItems} {columns} {filters} {initialState} {rowCount} queryParam='search'/>
	{/if}
	<div class="flex justify-end gap-4">
		<Button onclick={() => print(true)}>Print all instances</Button>
		<Button onclick={() => print(false)}>Print all tests</Button>
	</div>
</div>
