<script lang="ts">
	import { page } from '$app/state';
	import type { TestInstanceEventDTO, TestInstanceGetTelemetryResponse } from '$lib/api_types';
	import { DataTable } from '$lib/components/ui/data-table';
	import { API, ApiError, decodeBase64UrlToJson } from '$lib/services/api.svelte';
	import { onMount } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { columns, filters } from './schema';
	import type { InitialTableState } from '@tanstack/table-core';

	let rowItems: TestInstanceEventDTO[] = $state([]);
	let rowCount = $state(0);
	let loading = $state(false);
	let initialState: InitialTableState = $state({});

	const loadEvents = () => {
		API.request<null, TestInstanceGetTelemetryResponse>(
			`/api/v2/courses/${page.params.courseId}/tests/${page.params.itemId}/instance/${page.params.instanceId}/telemetry`
		)
			.then((res) => {
				rowItems = res.items;
				rowCount = res.items.length;
				loading = false;
			})
			.catch(() => {});
	};

	onMount(async () => {
		await loadEvents();
		const encodedParams = page.url.searchParams.get('search');
		if (encodedParams) {
			initialState = decodeBase64UrlToJson(encodedParams);
		}
	});
</script>

<!-- {JSON.stringify(rowItems)} -->
<DataTable
	data={rowItems}
	{columns}
	{filters}
	{initialState}
	{rowCount}
	paginationEnabled={false}
	queryParam='search'
/>

<!-- {"EventType":"","EventData":null,"PageID":"3379fb37-354d-49e9-a849-aba04e5640f3","User":{"id":1,"firstName":"Daniel","familyName":"Admin","username":"MAK0065","email":"mak0065@vsb.cz"}} -->
