<script lang="ts">
	import { onMount } from 'svelte';

	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import type { JoinedStudentDTO } from '$lib/api_types';
	import { type InitialTableState } from '@tanstack/table-core';
	import * as Form from "$lib/components/ui/form"
	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	let loading: boolean = $state(true);
	let rowItems: JoinedStudentDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});

	let { data } = $props();

	let showHistory = $state(data.showHistory)

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

<div class="flex flex-col gap-8 m-8">
	<div class="flex flex-row justify-between">
		<h1 class="text-2xl">Signed students</h1>
		<div class="flex gap-2">
		<!-- <Button href="{base}/app/{page.params.courseId}/tutor/questions/0">{m.quesstions_add()}</Button> -->
		</div>
	</div>
		<div class="flex flex-row justify-between">
			<Form.Checkbox
				title="Show history"
				name="showHistory"
				id="showHistory"
				bind:value={showHistory}
				error=""
				onclick={() => {
					const newUrl = new URL(page.url);
					if (page.url.searchParams.get("showHistory") == null) {
						newUrl.searchParams.set("showHistory", "true");
					} else {
						newUrl.searchParams.delete("showHistory");
					}
					goto(newUrl, { replaceState: true });
				}}
			>
			</Form.Checkbox>
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
