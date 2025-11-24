<script lang="ts">
	import { page } from '$app/state';
	import type { ClassUserDTO, RemoveTutorRequest, RemoveTutorResponse } from '$lib/api_types';
	import { DataTable } from '$lib/components/ui/data-table';
	import { API } from '$lib/services/api.svelte';
	import { tableConfig } from './schema';

	const actionsColumn = tableConfig.columns.find((c) => c.id == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (id: number) => {
				await API.request<RemoveTutorRequest, RemoveTutorResponse>(
					`api/v2/courses/${page.params.courseId}/classes/${page.params.classId}/tutors`,
					{
						method: 'DELETE',
						body: {
							userId: id
						}
					}
				)
					.then((res) => {
						tutors = res.tutors;
					})
					.catch(() => {});

				return true;
			}
		};
	}

	let {
		tutors = $bindable()
	}: {
		tutors: ClassUserDTO[];
	} = $props();

	let rowItems: ClassUserDTO[] = $derived(tutors);
	let rowCount: number = $derived(tutors.length);
</script>

<DataTable data={rowItems} {rowCount} {...tableConfig}></DataTable>
