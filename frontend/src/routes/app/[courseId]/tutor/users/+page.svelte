<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { tableConfig } from './schema';
	import { API } from '$lib/services/api.svelte';
	import {
		CourseUserRoleEnum,
		type CourseUserDTO,
		type RemoveCourseUserRequest,
		type RemoveCourseUserResponse
	} from '$lib/api_types';
	import { page } from '$app/state';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { m } from '$lib/paraglide/messages';
	import { invalidate } from '$app/navigation';
	import UserAddDialog from './UserAddDialog/UserAddDialog.svelte';
	import GlobalState from '$lib/shared.svelte';
	import { toast } from 'svelte-sonner';

	let rowItems: CourseUserDTO[] = $state([]);
	let rowCount: number = $state(0);
	let dialogOpen = $state(false);

	let { data } = $props();

	$effect(() => {
		data.data
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
							return url.pathname.endsWith(`courses/${page.params.courseId}/users`);
						});
						toast.success(m.toast_course_user_deleted());
					})
					.catch(() => {});

				return true;
			}
		};
	}

	const rolesColumn = tableConfig.columns.find((c) => c.id == 'roles');
	if (rolesColumn) {
		rolesColumn.meta = {
			...(rolesColumn.meta ?? {}),
			clickEventHandler: async (
				event: string,
				id: number,
				params: { role: CourseUserRoleEnum }
			) => {
				switch (event) {
					case 'remove_role':
						if (!confirm(m.user_role_remove_confirm())) {
							return;
						}
						await API.request<RemoveCourseUserRequest, RemoveCourseUserResponse>(
							`api/v2/courses/${page.params.courseId}/users`,
							{
								method: 'DELETE',
								body: {
									userId: id,
									role: params.role
								}
							}
						)
							.then((res) => {
								invalidate((url) => {
									return url.pathname.endsWith(`courses/${page.params.courseId}/users`);
								});
								toast.success(
									m.toast_course_user_role_removed({
										role: m.course_user_role_enum({ value: params.role })
									})
								);
							})
							.catch(() => {});

						return true;
				}
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
					<Dialog.Trigger class={buttonVariants({ variant: 'outline' })}>
						{m.course_user_add()}
					</Dialog.Trigger>
					{#if dialogOpen}
						<UserAddDialog
							endpoint={`api/v2/courses/${page.params.courseId}/users`}
						></UserAddDialog>
					{/if}
				</Dialog.Root>
			{/if}
		</div>
	</div>
	<DataTable data={rowItems} {rowCount} {...tableConfig} />
</div>
