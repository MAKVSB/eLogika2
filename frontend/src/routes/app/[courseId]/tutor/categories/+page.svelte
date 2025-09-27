<script lang="ts">
	import { onMount } from 'svelte';

	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import { decodeBase64UrlToJson } from '$lib/services/api.svelte';
	import { type InitialTableState } from '@tanstack/table-core';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import type { CategoryListItemDTO } from '$lib/api_types';
	import { m } from '$lib/paraglide/messages';
	import { base } from '$app/paths';

	let loading: boolean = $state(true);
	let rowItems: CategoryListItemDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});

	let { data } = $props();

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
		<h1 class="mb-8 text-2xl">Categories management</h1>
		<div class="flex gap-2">
			<Button href="{base}/app/{page.params.courseId}/tutor/categories/0">{m.category_add()}</Button>
		</div>
	</div>
	{#if !loading}
		<DataTable data={rowItems} {columns} {filters} {initialState} {rowCount} queryParam='search'/>
	{/if}
</div>
