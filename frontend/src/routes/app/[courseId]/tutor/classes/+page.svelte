<script lang="ts">
	import { onMount } from 'svelte';

	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import { API, ApiError, decodeBase64UrlToJson } from '$lib/services/api.svelte';
	import { CourseUserRoleEnum, type ClassListItemDTO } from '$lib/api_types';
	import { type InitialTableState } from '@tanstack/table-core';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import { m } from '$lib/paraglide/messages';
	import GlobalState from '$lib/shared.svelte';
	import { invalidateAll } from '$app/navigation';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import { base } from '$app/paths';

	let loading: boolean = $state(true);
	let rowItems: ClassListItemDTO[] = $state([]);
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
						if (!confirm('Question will be deleted permanently.')) {
							return;
						}
						API.request<any, Blob>(
							`/api/v2/courses/${page.params.courseId}/classes/${id}`,
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
			.catch(() => {})
			.finally(() => {
				loading = false;
			});
	});

	onMount(() => {
		const encodedParams = page.url.searchParams.get('search');
		if (encodedParams) {
			initialState = decodeBase64UrlToJson(encodedParams);
		}
	});
</script>

<div class="m-8">
	<div class="flex flex-row justify-between">
		<h1 class="mb-8 text-2xl">Class management</h1>
		{#if GlobalState.activeRole && [CourseUserRoleEnum.ADMIN, CourseUserRoleEnum.GARANT].includes(GlobalState.activeRole)}
			<Button href="{base}/app/{page.params.courseId}/tutor/classes/0">{m.class_add()}</Button>
		{/if}
	</div>
	{#if loading}
		<Pageloader></Pageloader>
	{:else}
		<DataTable data={rowItems} {columns} {filters} {initialState} {rowCount} queryParam='class_search'/>
	{/if}
</div>
