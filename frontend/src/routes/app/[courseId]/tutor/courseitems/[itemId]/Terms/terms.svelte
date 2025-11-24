<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { tableConfig } from './schema';
	import { API } from '$lib/services/api.svelte';
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import type { TermDTO, TermsListRequest, TermsListResponse } from '$lib/api_types';
	import { base } from '$app/paths';
	import { m } from '$lib/paraglide/messages';
	import { DataTableSearchParams } from '$lib/api_types_static';

	let data: TermDTO[] = $state([]);
	let rowCount: number = $state(0);

	let {
		courseId,
		itemId
	}: {
		courseId: string;
		itemId: string;
	} = $props();

	const actionsColumn = tableConfig.columns.find((c) => c.id == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (event: string, id: number) => {
				switch (event) {
					case 'delete':
						if (!confirm(m.term_delete_confirm())) {
							return;
						}
						API.request<any, Blob>(
							`/api/v2/courses/${page.params.courseId}/items/${page.params.itemId}/terms/${id}`,
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
		const search =
			page.url.searchParams.get(tableConfig.searchParam) ??
			DataTableSearchParams.fromDataTable(tableConfig.initialState).toURL();

		API.request<TermsListRequest, TermsListResponse>(
			`/api/v2/courses/${courseId}/items/${itemId}/terms`,
			{
				searchParams: {
					search
				}
			}
		)
			.then((res) => {
				data = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {});
	});
</script>

<div>
	<div class="flex flex-row justify-between">
		<h1 class="mb-8 text-2xl">Terms</h1>
		<div class="flex gap-2">
			<Button href="{base}/app/{courseId}/tutor/courseitems/{itemId}/0">Add term</Button>
		</div>
	</div>
	<DataTable {data} {rowCount} {...tableConfig} />
</div>
