<script lang="ts">
	import {
		StudyFormEnum,
		CourseItemTypeEnum,
		type CourseItemDTO,
		type CourseItemListResponse
	} from '$lib/api_types';
	import Button from '$lib/components/ui/button/button.svelte';
	import { DataTable } from '$lib/components/ui/data-table';
	import { initialState as is, tableConfig } from './schema';
	import { m } from '$lib/paraglide/messages';
	import { page } from '$app/state';
	import { invalidateAll } from '$app/navigation';
	import { API } from '$lib/services/api.svelte';
	import { base } from '$app/paths';
	import type {
		ColumnFiltersState,
		PaginationState,
		SortingState
	} from '@tanstack/table-core';
	import { DataTableSearchParams } from '$lib/api_types_static';

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
	let searchParam = $derived('search' + mode);

	const initialState = {
		...is,
		columnFilters: [
			...(is.columnFilters ?? []),
			{
				id: 'StudyForm',
				value: mode
			},
			{
				id: 'ParentID',
				value: parentId ?? 'NULL'
			}
		]
	};

	const actionsColumn = tableConfig.columns.find((c) => c.id == 'actions');
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

	$effect(() => {
		const search =
			page.url.searchParams.get(searchParam) ??
			DataTableSearchParams.fromDataTable(initialState).toURL();

		API.request<null, CourseItemListResponse>(`/api/v2/courses/${page.params.courseId}/items`, {
			searchParams: {
				search
			}
		})
			.then((res) => {
				data = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {});
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
			href={`${base}/app/${courseId}/tutor/courseitems/0?studyForm=${mode}&type=${CourseItemTypeEnum.GROUP}${parentId ? `&parentId=${parentId}` : ''}`}
		>
			{m.course_add_group()}
		</Button>
	{/if}
	<Button
		href={`${base}/app/${courseId}/tutor/courseitems/0?studyForm=${mode}&type=${CourseItemTypeEnum.ACTIVITY}${parentId ? `&parentId=${parentId}` : ''}`}
	>
		{m.course_add_activity()}
	</Button>
	<Button
		href={`${base}/app/${courseId}/tutor/courseitems/0?studyForm=${mode}&type=${CourseItemTypeEnum.TEST}${parentId ? `&parentId=${parentId}` : ''}`}
	>
		{m.course_add_test()}
	</Button>
</div>

<DataTable {data} {rowCount} {initialState} {searchParam} {...tableConfig} />
