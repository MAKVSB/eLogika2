<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters, tableConfig } from './schema';
	import type { QuestionVersionDTO, QuestionSelectVersionResponse } from '$lib/api_types';
	import { type InitialTableState } from '@tanstack/table-core';
	import { page } from '$app/state';
	import { API } from '$lib/services/api.svelte';
	import { toast } from 'svelte-sonner';

	let {
		versions
	}: {
		versions?: QuestionVersionDTO[];
	} = $props();

	let rowCount: number = $derived(versions?.length ?? 0);

	const actionsColumn = columns.find((c) => c.id == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (event: string, id: number) => {
				switch (event) {
					case 'selectversion':
						await API.request<null, QuestionSelectVersionResponse>(
							`api/v2/courses/${page.params.courseId}/questions/${id}/selectversion`,
							{
								method: 'PATCH'
							}
						)
							.then((res) => {
								versions = res.data.versions;
								toast.success('Active version changed');
							})
							.catch(() => {});
						break;
				}

				return true;
			}
		};
	}
</script>

{#if versions}
	<DataTable
		data={versions}
		{rowCount}
		{...tableConfig}
	/>
{/if}
