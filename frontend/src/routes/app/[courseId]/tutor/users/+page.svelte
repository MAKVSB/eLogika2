<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import { API } from '$lib/services/api.svelte';
	import {
		CourseUserRoleEnum,
		type CourseUserDTO,
		type RemoveCourseUserRequest,
		type RemoveCourseUserResponse
	} from '$lib/api_types';
	import { type InitialTableState } from '@tanstack/table-core';
	import { page } from '$app/state';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { m } from '$lib/paraglide/messages';
	import { invalidate } from '$app/navigation';
	import UserAddDialog from './UserAddDialog/UserAddDialog.svelte';
	// import UserAddDialog from '../UserAddDialog/UserAddDialog.svelte';
	import GlobalState from '$lib/shared.svelte';

	let loading: boolean = $state(true);
	let rowItems: CourseUserDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});
	let dialogOpen = $state(false);

	let { data } = $props();

	$effect(() => {
		data.data
			.then((res) => {
				rowItems = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {})
			.finally(() => {
				loading = false;
			});
	});

	const actionsColumn = columns.find((c) => c.uniqueId == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (id: number) => {
				await API.request<RemoveCourseUserRequest, RemoveCourseUserResponse>(
					`api/v2/courses/${page.params.courseId}/users`,
					{
						method: 'DELETE',
						body: {
							userId: id
						}
					}
				)
					.then((res) => {
						invalidate((url) => {
							return url.href.endsWith(`courses/${page.params.courseId}/users`);
						});
					})
					.catch(() => {});

				return true;
			}
		};
	}
</script>

<div class="m-8">
	<div class="flex flex-row justify-between">
		<h1 class="mb-8 text-2xl">Course users</h1>
		<div>
			{#if GlobalState.activeRole == CourseUserRoleEnum.ADMIN}
				<Dialog.Root bind:open={dialogOpen}>
					<Dialog.Trigger class={buttonVariants({ variant: 'outline' })}
						>{m.class_student_add()}</Dialog.Trigger
					>
					{#if dialogOpen}
						<UserAddDialog endpoint={`api/v2/courses/${page.params.courseId}/users`}
						></UserAddDialog>
					{/if}
				</Dialog.Root>
			{/if}
		</div>
	</div>
	{#if !loading}
		<DataTable data={rowItems} {columns} {filters} {initialState} {rowCount} queryParam="search" />
	{/if}
</div>
