<script lang="ts">
	import { page } from '$app/state';
	import type {
		AddStudentRequest,
		AddStudentResponse,
		CourseUserDTO,
		ListCourseUsersResponse
	} from '$lib/api_types';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { API } from '$lib/services/api.svelte';
	import { tableConfig } from './schema';
	import { DataTable } from '$lib/components/ui/data-table';
	import { invalidateAll } from '$app/navigation';
	import { DataTableSearchParams } from '$lib/api_types_static';

	let {
		defaultRole,
		endpoint
	}: {
		defaultRole?: string;
		endpoint: string;
	} = $props();

	let search = $state('');

	let rowItems: CourseUserDTO[] = $state([]);
	let rowCount: number = $state(0);

	const actionsColumn = tableConfig.columns.find((c) => c.id == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (id: number) => {
				await API.request<AddStudentRequest, AddStudentResponse>(endpoint, {
					method: 'POST',
					body: {
						userId: id
					}
				})
					.then((res) => {
						invalidateAll();
					})
					.catch(() => {});

				return true;
			}
		};
	}

	$effect(() => {
		const search =
			page.url.searchParams.get(tableConfig.searchParam) ??
			DataTableSearchParams.fromDataTable(tableConfig.initialState).toURL();

		API.request<null, ListCourseUsersResponse>(
			`/api/v2/courses/${page.params.courseId}/users`,
			{
				searchParams: {
					search
				}
			},
			fetch
		)
			.then((res) => {
				rowItems = res.items;
				rowCount = res.items.length;
			})
			.catch(() => {});
	});
</script>

<Dialog.Content class="max-h-full w-300 overflow-scroll sm:max-h-[90%] sm:max-w-[90%]">
	<Dialog.Header>
		<Dialog.Title>Add tutor to class</Dialog.Title>
	</Dialog.Header>
	<DataTable data={rowItems} {rowCount} {...tableConfig} />
</Dialog.Content>
