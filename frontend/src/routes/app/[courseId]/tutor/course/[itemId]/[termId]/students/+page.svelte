<script lang="ts">
	import { onMount } from 'svelte';

	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import { API, ApiError, decodeBase64UrlToJson } from '$lib/services/api.svelte';
	import type { JoinedStudentDTO } from '$lib/api_types';
	import { type InitialTableState } from '@tanstack/table-core';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import { m } from '$lib/paraglide/messages';

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
		const encodedParams = page.url.searchParams.get('search');
		if (encodedParams) {
			initialState = decodeBase64UrlToJson(encodedParams);
		}
		loading = false;
	});
</script>

<div class="m-8">
	<div class="flex flex-row justify-between">
		<h1 class="mb-8 text-2xl">Signed students</h1>
		<!-- <Button href="{base}/app/{page.params.courseId}/tutor/questions/0">{m.quesstions_add()}</Button> -->
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
