<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import { type InitialTableState } from '@tanstack/table-core';
	import Button from '$lib/components/ui/button/button.svelte';
	import type { UserListItemDTO } from '$lib/api_types';
	import { base } from '$app/paths';

	let { data } = $props();
	let isLoading: boolean = $state(true);

	$effect(() => {
		data.users
			.then((res) => {
				rowItems = res.items;
				rowCount = res.itemsCount;
				isLoading = false;
			})
			.catch(() => {});
	});

	let rowItems: UserListItemDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});
</script>

<div class="m-8">
	<div class="flex flex-row justify-between">
		<h1 class="mb-8 text-2xl">Admin user management</h1>
		<div class="flex gap-2">
			<Button href="{base}/app/admin/users/0">Add user</Button>
		</div>
	</div>
	{#if !isLoading}
		<DataTable data={rowItems} {columns} {filters} {initialState} {rowCount} queryParam="search" />
	{/if}
</div>
