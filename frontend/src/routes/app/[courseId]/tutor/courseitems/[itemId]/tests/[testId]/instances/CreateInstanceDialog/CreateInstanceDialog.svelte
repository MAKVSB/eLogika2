<script lang="ts">
	import { page } from '$app/state';
	import {
		CourseUserRoleEnum,
		TestInstanceFormEnum,
		type CourseUserDTO,
		type ListCourseUsersResponse,
		type TestInstanceCreateRequest,
		type TestInstanceCreateResponse
	} from '$lib/api_types';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { API } from '$lib/services/api.svelte';
	import { columns, filters } from './schema';
	import { DataTable } from '$lib/components/ui/data-table';
	import type { InitialTableState, TableState } from '@tanstack/table-core';
	import * as Form from '$lib/components/ui/form';
	import { goto, invalidateAll } from '$app/navigation';
	import { enumToOptions } from '$lib/utils';
	import { base } from '$app/paths';

	$effect(() => {
		loadData(page.url.searchParams.get('createInstanceSearch'))
	})

	let {
		openState = $bindable()
	}: {
		openState: boolean;
	} = $props();

	let rowItems: CourseUserDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({
		columnFilters: [
			{
				id: 'role',
				value: CourseUserRoleEnum.STUDENT
			}
		]
	});

	let instanceForm = $state(TestInstanceFormEnum.OFFLINE);

	const loadData = (search: string | null) => {
		API.request<null, ListCourseUsersResponse>(
			`/api/v2/courses/${page.params.courseId}/users`,
			{
				searchParams: {
					...(search ? { search: search } : {})
				}
			},
			fetch
		)
			.then((res) => {
				rowItems = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {});
	};

	const actionsColumn = columns.find((c) => c.uniqueId == 'actions');
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
							userId: id,
						}
					}
				)
					.then((res: TestInstanceCreateResponse) => {
						invalidateAll();
						console.log("Transfering 28")
						goto(
							base+`/app/${page.params.courseId}/tutor/courseitems/${page.params.itemId}/tests/${page.params.testId}/instances/${res.instanceId}`
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

	<DataTable
		data={rowItems}
		{columns}
		{filters}
		{rowCount}
		{initialState}
		queryParam="createInstanceSearch"
		replaceState
	/>
</Dialog.Content>
