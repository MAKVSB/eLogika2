<script lang="ts">
	import { page } from '$app/state';
	import type { TestInstanceEventDTO, TestInstanceGetTelemetryResponse } from '$lib/api_types';
	import { DataTable } from '$lib/components/ui/data-table';
	import { API } from '$lib/services/api.svelte';
	import { toast } from 'svelte-sonner';
	import { tableConfig } from './schema';
	import { DataTableSearchParams } from '$lib/api_types_static';

	let rowItems: TestInstanceEventDTO[] = $state([]);
	let rowCount = $state(0);

	$effect(() => {
		const search =
			page.url.searchParams.get(tableConfig.searchParam) ??
			DataTableSearchParams.fromDataTable(tableConfig.initialState).toURL();

		API.request<null, TestInstanceGetTelemetryResponse>(
			`/api/v2/courses/${page.params.courseId}/tests/${page.params.itemId}/instance/${page.params.instanceId}/telemetry`,
			{
				searchParams: {
					search
				}
			}
		)
			.then((res) => {
				rowItems = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {});
	});
</script>

<DataTable data={rowItems} {rowCount} {...tableConfig} />
