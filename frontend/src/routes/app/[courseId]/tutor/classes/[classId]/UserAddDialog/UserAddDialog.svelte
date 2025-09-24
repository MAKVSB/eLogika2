<script lang="ts">
	import { page } from '$app/state';
	import type {
		AddStudentRequest,
		AddStudentResponse,
		CourseUserDTO,
		ListCourseUsersResponse
	} from '$lib/api_types';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import {
		API,
		encodeJsonToBase64Url
	} from '$lib/services/api.svelte';
	import { onMount } from 'svelte';
	import { columns, filters } from './schema';
	import { DataTable } from '$lib/components/ui/data-table';
	import type {
		ColumnFiltersState,
		InitialTableState,
		SortingState,
		TableState
	} from '@tanstack/table-core';
	import Loader from '$lib/components/ui/loader/loader.svelte';
	import { invalidateAll } from '$app/navigation';

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
	let initialState: InitialTableState = $state({});
	let loading = $state(true);

	const actionsColumn = columns.find((c) => c.uniqueId == 'actions');
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

	const loadData = () => {
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
				rowCount = res.items.length;
			})
			.catch(() => {});
	};

	type RestRequest = {
		sorting?: SortingState;
		columnFilters?: ColumnFiltersState;
	};

	function refetch(state: TableState) {
		const queryParams: RestRequest = {
			...(state.sorting ? { sorting: state.sorting } : {}),
			...(state.columnFilters ? { columnFilters: state.columnFilters } : {})
		};
		search = encodeJsonToBase64Url(queryParams);
		loadData();
	}

	onMount(() => {
		initialState.columnFilters = [
			{
				id: 'roles',
				value: defaultRole
			}
		];
		search = encodeJsonToBase64Url(initialState);
		loadData();
		loading = false;
	});
</script>

<Dialog.Content class="max-h-full w-300 overflow-scroll sm:max-h-[90%] sm:max-w-[90%]">
	<Dialog.Header>
		<Dialog.Title>Add tutor to class</Dialog.Title>
	</Dialog.Header>
	{#if !loading}
		<DataTable
			data={rowItems}
			{columns}
			{filters}
			{refetch}
			{rowCount}
			queryParam='usersearch'
			{initialState}
		/>
	{:else}
		<Loader></Loader>
	{/if}
</Dialog.Content>
