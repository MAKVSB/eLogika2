<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, tableConfig } from './schema.svelte';
	import { API } from '$lib/services/api.svelte';
	import type { SupportTicketListItemDTO } from '$lib/api_types';
	import { page } from '$app/state';
	import { m } from '$lib/paraglide/messages';
	import { invalidateAll } from '$app/navigation';

	let rowItems: SupportTicketListItemDTO[] = $state([]);
	let rowCount: number = $state(0);

	let { data } = $props();

	$effect(() => {
		data.data
			.then((res) => {
				rowItems = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {});
	});

	const actionsColumn = columns.find((c) => c.id == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (event: string, id: number) => {
				switch (event) {
					case 'delete':
						if (!confirm(m.question_delete_confirm())) {
							return;
						}
						API.request<any, Blob>(
							`/api/v2/courses/${page.params.courseId}/questions/${id}`,
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
</script>

<div class="flex flex-col gap-8 m-8">
	<div class="flex flex-row justify-between">
		<h1 class="text-2xl">Support tickets</h1>
	</div>
	<div>
		<DataTable data={rowItems} {rowCount} {...tableConfig} />
	</div>
</div>
