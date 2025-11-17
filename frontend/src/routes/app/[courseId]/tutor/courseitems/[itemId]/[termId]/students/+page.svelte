<script lang="ts">
	import { onMount } from 'svelte';

	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import type { JoinedStudentDTO } from '$lib/api_types';
	import { type InitialTableState } from '@tanstack/table-core';

	let loading: boolean = $state(true);
	let rowItems: JoinedStudentDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});

	let { data } = $props();

	$effect(() => {
		data.students
			.then((res) => {
				rowItems = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {});
	});

	onMount(() => {
		loading = false;
	});
</script>

<div class="m-8">
	<div class="flex flex-row justify-between">
		<h1 class="mb-8 text-2xl">Signed students</h1>
		<div class="flex gap-2">
		<!-- <Button href="{base}/app/{page.params.courseId}/tutor/questions/0">{m.quesstions_add()}</Button> -->
		</div>
	</div>
	{#if !loading}
		<DataTable
			data={rowItems}
			{columns}
			{filters}
			{initialState}
			{rowCount}
			paginationEnabled={false}
			queryParam='search'
		/>
	{/if}
</div>
