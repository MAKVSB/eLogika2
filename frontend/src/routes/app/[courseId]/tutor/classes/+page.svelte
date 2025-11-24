<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { tableConfig } from './schema';
	import { API } from '$lib/services/api.svelte';
	import {
		CourseUserRoleEnum,
		type ClassImportClassesResponse,
		type ClassListItemDTO
	} from '$lib/api_types';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages';
	import GlobalState from '$lib/shared.svelte';
	import { invalidate, invalidateAll } from '$app/navigation';
	import { base } from '$app/paths';
	import { toast } from 'svelte-sonner';
	import Loader from '$lib/components/ui/loader/loader.svelte';

	let rowItems: ClassListItemDTO[] = $state([]);
	let rowCount: number = $state(0);

	let importRunning = $state(false);

	let { data } = $props();

	const actionsColumn = tableConfig.columns.find((c) => c.id == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (event: string, id: number) => {
				switch (event) {
					case 'delete':
						if (!confirm(m.class_delete_confirm())) {
							return;
						}
						API.request<any, Blob>(
							`/api/v2/courses/${page.params.courseId}/classes/${id}`,
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

	$effect(() => {
		data.data
			.then((res) => {
				rowItems = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {});
	});

	const importClasses = async () => {
		importRunning = true;
		await API.request<null, ClassImportClassesResponse>(
			`api/v2/courses/${page.params.courseId}/classes/import`,
			{
				method: 'POST'
			}
		)
			.then((res) => {
				invalidate((url) => {
					return url.href.endsWith(`/api/v2/courses/${page.params.courseId}/classes`);
				});
				toast.success('Import finished');
			})
			.catch(() => {})
			.finally(() => {
				importRunning = false;
			});

		return true;
	};
</script>

<div class="m-8">
	<div class="flex flex-row justify-between">
		<h1 class="mb-8 text-2xl">Class management</h1>
		<div class="flex gap-2">
			{#if GlobalState.activeRole && [CourseUserRoleEnum.ADMIN, CourseUserRoleEnum.GARANT].includes(GlobalState.activeRole)}
				<Button onclick={() => importClasses()} disabled={importRunning}>
					{#if importRunning}
						<Loader></Loader>
					{/if}
					Import classes
				</Button>
				<Button href="{base}/app/{page.params.courseId}/tutor/classes/0">{m.class_add()}</Button>
			{/if}
		</div>
	</div>
	<DataTable data={rowItems} {rowCount} {...tableConfig} />
</div>
