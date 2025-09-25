<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import {
		type InitialTableState,
	} from '@tanstack/table-core';
	import Button from '$lib/components/ui/button/button.svelte';
	import type { CourseListItemDTO } from '$lib/api_types';
	import { base } from '$app/paths';

	let { data } = $props();
	let isLoading: boolean = $state(true);

	$effect(() => {
		data.courses
			.then((res) => {
				rowItems = res.items;
				rowCount = res.itemsCount;
				isLoading = false;
			})
			.catch(() => {});
	});

	let rowItems: CourseListItemDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});

</script>

<div class="m-8">
	<div class="flex flex-row justify-between">
		<h1 class="mb-8 text-2xl">Courses management</h1>
		<Button href="{base}/app/admin/courses/0">Add course</Button>
	</div>
	{#if !isLoading}
		<DataTable data={rowItems} {columns} {filters} {initialState} {rowCount} queryParam='search'/>
	{/if}
</div>
