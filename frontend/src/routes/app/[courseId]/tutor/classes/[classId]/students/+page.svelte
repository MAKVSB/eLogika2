<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { tableConfig } from './schema';
	import { API } from '$lib/services/api.svelte';
	import type {
		ClassImportStudentsResponse,
		ClassUserDTO,
		RemoveStudentRequest,
		RemoveStudentResponse
	} from '$lib/api_types';
	import { page } from '$app/state';
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import UserAddDialog from '../UserAddDialog/UserAddDialog.svelte';
	import { m } from '$lib/paraglide/messages';
	import { invalidate } from '$app/navigation';

	let rowItems: ClassUserDTO[] = $state([]);
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
				await API.request<RemoveStudentRequest, RemoveStudentResponse>(
					`api/v2/courses/${page.params.courseId}/classes/${page.params.classId}/students`,
					{
						method: 'DELETE',
						body: {
							userId: id
						}
					}
				)
					.then((res) => {
						rowItems = res.students;
						rowCount = res.students.length;
					})
					.catch(() => {});

				return true;
			}
		};
	}

	async function importStudents() {
		await API.request<null, ClassImportStudentsResponse>(
			`api/v2/courses/${page.params.courseId}/classes/${page.params.classId}/students/import`,
			{
				method: 'POST'
			}
		)
			.then((res) => {
				invalidate((url) => {
					return url.href.endsWith(
						`/api/v2/courses/${page.params.courseId}/classes/${page.params.classId}/students`
					);
				});
			})
			.catch(() => {});

		return true;
	}
</script>

<div class="m-8">
	<div class="flex flex-row justify-between">
		<h1 class="mb-8 text-2xl">Class students</h1>
		<div class="flex gap-2">
			<Button variant="outline" onclick={() => importStudents()}>{m.class_students_import()}</Button
			>
			<Dialog.Root bind:open={dialogOpen}>
				<Dialog.Trigger class={buttonVariants({ variant: 'outline' })}
					>{m.class_student_add()}</Dialog.Trigger
				>
				{#if dialogOpen}
					<UserAddDialog
						defaultRole="STUDENT"
						endpoint={`api/v2/courses/${page.params.courseId}/classes/${page.params.classId}/students`}
					></UserAddDialog>
				{/if}
			</Dialog.Root>
		</div>
	</div>
	<DataTable data={rowItems} {rowCount} {...tableConfig} />
</div>
