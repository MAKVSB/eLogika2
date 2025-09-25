<script lang="ts">
	import {
		StudyFormEnum,
		CourseItemTypeEnum,
		type CourseItemDTO,
		type CourseItemListResponse
	} from '$lib/api_types';
	import Button from '$lib/components/ui/button/button.svelte';
	import { DataTable } from '$lib/components/ui/data-table';
	import { columns, filters } from './schema';
	import { m } from '$lib/paraglide/messages';
	import { page } from '$app/state';
	import { goto, invalidateAll } from '$app/navigation';
	import { API, decodeBase64UrlToJson, encodeJsonToBase64Url } from '$lib/services/api.svelte';
	import GlobalState from '$lib/shared.svelte';
	import { base } from '$app/paths';
	import type {
		ColumnFiltersState,
		InitialTableState,
		PaginationState,
		SortingState,
		TableState
	} from '@tanstack/table-core';
	import { onMount } from 'svelte';

	let loading = $state(false);
	let {
		mode,
		courseId,
		parentId,
		canGroup = false
	}: {
		mode: StudyFormEnum;
		parentId?: number;
		courseId: string;
		canGroup: boolean;
	} = $props();

	let data: CourseItemDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});
	let encodedParams: string | null = $state(null);

	$effect(() => {
		if (!page.params.courseId) return;
		if (encodedParams) {
			console.log("Transfering 29")
			goto(`?search=${encodedParams}`);
			fetchData();
		} else {
			fetchData();
		}
	});

	const actionsColumn = columns.find((c) => c.uniqueId == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (event: string, id: number) => {
				switch (event) {
					case 'delete':
						if (!confirm('Course item will be deleted permanently.')) {
							return;
						}
						API.request<any, Blob>(
							`/api/v2/courses/${page.params.courseId}/items/${id}`,
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

	type RestRequest = {
		pagination?: PaginationState;
		sorting?: SortingState;
		columnFilters?: ColumnFiltersState;
	};

	async function fetchData() {
		const queryParams: RestRequest = {
			pagination: {
				pageIndex: 0,
				pageSize: 1000
			},
			columnFilters: [
				{
					id: 'StudyForm',
					value: mode
				}
			]
		};

		if (parentId) {
			queryParams.columnFilters?.push({
				id: 'ParentID',
				value: parentId
			});
		} else {
			queryParams.columnFilters?.push({
				id: 'ParentID',
				value: 'NULL'
			});
		}

		await API.request<null, CourseItemListResponse>(
			`/api/v2/courses/${page.params.courseId}/items`,
			{
				searchParams: {
					search: encodeJsonToBase64Url(queryParams)
				}
			}
		)
			.then((res) => {
				data = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {});
	}

	function refetch(state: TableState) {
		const queryParams: RestRequest = {
			...(state.pagination ? { pagination: state.pagination } : {}),
			...(state.sorting ? { sorting: state.sorting } : {}),
			...(state.columnFilters ? { columnFilters: state.columnFilters } : {})
		};
		encodedParams = encodeJsonToBase64Url(queryParams);
	}

	onMount(() => {
		const encodedParams = page.url.searchParams.get('search');
		if (encodedParams) {
			initialState = decodeBase64UrlToJson(encodedParams);
		}
		loading = false;
	});
</script>

<div class="flex items-center gap-4">
	<h1 class="text-xl font-bold">
		{#if mode}
			{m.study_form_enum({ value: mode })}
		{:else}
			Group items
		{/if}
	</h1>
	<div class="flex-1"></div>
	{#if canGroup}
		<Button
			href={`${base}/app/${courseId}/tutor/course/0?studyForm=${mode}&type=${CourseItemTypeEnum.GROUP}${parentId ? `&parentId=${parentId}` : ''}`}
		>
			{m.course_add_group()}
		</Button>
	{/if}
	<Button
		href={`${base}/app/${courseId}/tutor/course/0?studyForm=${mode}&type=${CourseItemTypeEnum.ACTIVITY}${parentId ? `&parentId=${parentId}` : ''}`}
	>
		{m.course_add_activity()}
	</Button>
	<Button
		href={`${base}/app/${courseId}/tutor/course/0?studyForm=${mode}&type=${CourseItemTypeEnum.TEST}${parentId ? `&parentId=${parentId}` : ''}`}
	>
		{m.course_add_test()}
	</Button>
</div>

{#if !loading}
	<DataTable
		{data}
		{columns}
		{filters}
		{refetch}
		{initialState}
		{rowCount}
		paginationEnabled={false}
		queryParam='search'
	/>
{/if}
