<script lang="ts">
	import { page } from '$app/state';
	import {
		TestInstanceFormEnum,
		type CourseUserDTO,
		type ListCourseUsersResponse,
		type TestInstanceCreateRequest,
		type TestInstanceCreateResponse
	} from '$lib/api_types';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { API } from '$lib/services/api.svelte';
	import { tableConfig } from './schema';
	import { DataTable } from '$lib/components/ui/data-table';
	import * as Form from '$lib/components/ui/form';
	import { goto, invalidateAll } from '$app/navigation';
	import { enumToOptions } from '$lib/utils';
	import { base } from '$app/paths';
	import { DataTableSearchParams } from '$lib/api_types_static';

	let {
		openState = $bindable()
	}: {
		openState: boolean;
	} = $props();

	let rowItems: CourseUserDTO[] = $state([]);
	let rowCount: number = $state(0);

	let instanceForm = $state(TestInstanceFormEnum.OFFLINE);

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
				rowCount = res.itemsCount;
			})
			.catch(() => {});
	});

	const actionsColumn = tableConfig.columns.find((c) => c.id == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (id: number) => {
				await API.request<TestInstanceCreateRequest, TestInstanceCreateResponse>(
					`/api/v2/courses/${page.params.courseId}/tests/${page.params.itemId}/instances/${page.params.testId}/create`,
					{
						method: 'POST',
						body: {
							form: instanceForm,
							userId: id
						}
					}
				)
					.then((res: TestInstanceCreateResponse) => {
						invalidateAll();
						console.log('Transfering 28');
						goto(
							base +
								`/app/${page.params.courseId}/tutor/courseitems/${page.params.itemId}/tests/${page.params.testId}/instances/${res.instanceId}`
						);
					})
					.catch(() => {});
				return true;
			}
		};
	}
</script>

<Dialog.Content class="max-h-full w-300 overflow-scroll sm:max-h-[90%] sm:max-w-[90%]">
	<Dialog.Header>
		<Dialog.Title>Create instance for student</Dialog.Title>
	</Dialog.Header>

	<Form.SingleSelect
		title="Instance form"
		name="instanceForm"
		id="instandeForm"
		bind:value={instanceForm}
		options={enumToOptions(TestInstanceFormEnum)}
		error=""
	></Form.SingleSelect>

	<DataTable data={rowItems} {rowCount} {...tableConfig} />
</Dialog.Content>
