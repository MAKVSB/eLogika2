<script lang="ts">
	import { onMount } from 'svelte';

	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import { API, decodeBase64UrlToJson } from '$lib/services/api.svelte';
	import { type InitialTableState } from '@tanstack/table-core';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import type { TemplateListItemDTO } from '$lib/api_types';
	import { m } from '$lib/paraglide/messages';
	import { invalidateAll } from '$app/navigation';
	import { base } from '$app/paths';

	let loading: boolean = $state(true);
	let rowItems: TemplateListItemDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});

	let { data } = $props();

	const actionsColumn = columns.find((c) => c.uniqueId == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (event: string, id: number) => {
				switch (event) {
					case 'delete':
						if (!confirm('Template will be deleted permanently.')) {
							return;
						}
						API.request<any, Blob>(
							`/api/v2/courses/${page.params.courseId}/templates/${id}`,
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
		data.data
			.then((res) => {
				rowItems = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {});
	});

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
		<h1 class="mb-8 text-2xl">Templates management</h1>
		<Button href="{base}/app/{page.params.courseId}/tutor/templates/0">{m.template_add()}</Button>
	</div>
	{#if !loading}
		<DataTable data={rowItems} {columns} {filters} {initialState} {rowCount} queryParam='search'/>
	{/if}
</div>
