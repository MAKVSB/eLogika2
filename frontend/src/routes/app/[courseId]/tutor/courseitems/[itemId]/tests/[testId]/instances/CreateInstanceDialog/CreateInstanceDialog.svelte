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
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { API, ApiError, encodeJsonToBase64Url } from '$lib/services/api.svelte';
	import { onMount } from 'svelte';
	import { columns, filters } from './schema';
	import { DataTable } from '$lib/components/ui/data-table';
	import type { ColumnFiltersState, InitialTableState, PaginationState, RowSelectionState, SortingState, TableState } from '@tanstack/table-core';
	import Loader from '$lib/components/ui/loader/loader.svelte';
	import * as Form from '$lib/components/ui/form';
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import { Input } from '$lib/components/ui/input';
	import { toast } from 'svelte-sonner';
	import { goto, invalidateAll } from '$app/navigation';
	import { enumToOptions } from '$lib/utils';
	import { base } from '$app/paths';

	type RestRequest = {
		pagination?: PaginationState;
		sorting?: SortingState;
		columnFilters?: ColumnFiltersState;
	};

	let {
		openState = $bindable()
	}: {
		openState: boolean;
	} = $props();

	let rowItems: CourseUserDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});
	let loading = $state(true);

	let instanceForm = $state(TestInstanceFormEnum.OFFLINE);

	const loadData = () => {
		initialState.columnFilters = [
			{
				id: 'role',
				value: CourseUserRoleEnum.STUDENT
			}
		];
		const search = encodeJsonToBase64Url(initialState);

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

	function refetch(state: TableState) {
		const queryParams: RestRequest = {
			...(state.pagination ? { pagination: state.pagination } : {}),
			...(state.sorting ? { sorting: state.sorting } : {}),
			...(state.columnFilters ? { columnFilters: state.columnFilters } : {})
		};
		const search = encodeJsonToBase64Url(queryParams);
		
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

	}

	onMount(() => {
		loadData();
		loading = false;
	});
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

	{#if !loading}
		<DataTable
			data={rowItems}
			{columns}
			{filters}
			{rowCount}
			{initialState}
			{refetch}
		/>
	{:else}
		<Loader></Loader>
	{/if}
</Dialog.Content>
