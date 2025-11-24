<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { tableConfig } from './schema.svelte';
	import { API } from '$lib/services/api.svelte';
	import type {
		QuestionCheckResponse,
		QuestionListItemDTO,
		QuestionToggleActiveResponse
	} from '$lib/api_types';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import { m } from '$lib/paraglide/messages';
	import { invalidateAll } from '$app/navigation';
	import { base } from '$app/paths';
	import { FilterTypeEnum } from '$lib/components/ui/data-table/filter';
	import Loader from '$lib/components/ui/loader/loader.svelte';
	import { DataTableSearchParams } from '$lib/api_types_static';

	let rowItems: QuestionListItemDTO[] = $state([]);
	let rowCount: number = $state(0);

	let printRunning = $state(false);

	let { data } = $props();

	$effect(() => {
		data.data
			.then((res) => {
				rowItems = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {});
	});

	$effect(() => {
		data.chapterData
			.then((res) => {
				const filter = tableConfig.filters.find((f) => f.accessorKey == 'chapterId');
				if (filter && filter.type == FilterTypeEnum.SELECT) {
					filter.values = res.items.map((chapter) => {
						return {
							value: chapter.id,
							display: chapter.name
						};
					});
				}
			})
			.catch(() => {});
	});

	$effect(() => {
		data.categoryData
			.then((res) => {
				const filter = tableConfig.filters.find((f) => f.accessorKey == 'categoryId');
				if (filter && filter.type == FilterTypeEnum.SELECT) {
					filter.values = res.items.map((chapter) => {
						return {
							value: chapter.id,
							display: chapter.name
						};
					});
				}
			})
			.catch(() => {});
	});

	const actionsColumn = tableConfig.columns.find((c) => c.id == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (event: string, id: number) => {
				switch (event) {
					case 'check':
					case 'uncheck':
						await API.request<null, QuestionCheckResponse>(
							`api/v2/courses/${page.params.courseId}/questions/${id}/check`,
							{
								method: event == 'check' ? 'POST' : 'DELETE'
							}
						)
							.then((res) => {
								let row = rowItems.find((val) => val.id == id);
								if (row) {
									row.checkedBy = res.checkedBy;
									toast.success('Saved');
								}
							})
							.catch(() => {});
						break;
					case 'print':
						API.request<any, Blob>(
							`/api/v2/courses/${page.params.courseId}/print/questions/${id}`,
							{
								method: 'POST'
							},
							fetch
						)
							.then((res) => {
								const url = URL.createObjectURL(res);
								window.open(url); // opens in new tab
							})
							.catch(() => {});
						break;
					case 'delete':
						if (!confirm(m.question_delete_confirm())) {
							return;
						}
						API.request<any, Blob>(
							`/api/v2/courses/${page.params.courseId}/questions/${id}`,
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

	const activeColumn = tableConfig.columns.find((c) => c.id == 'active');
	if (activeColumn) {
		activeColumn.meta = {
			...(activeColumn.meta ?? {}),
			clickEventHandler: async (id: number) => {
				await API.request<null, QuestionToggleActiveResponse>(
					`api/v2/courses/${page.params.courseId}/questions/${id}/toggleActive`,
					{
						method: 'PATCH'
					}
				)
					.then((res) => {
						let row = rowItems.find((val) => val.id == id);
						if (row) {
							row.active = res.active;
							toast.success('Saved');
						}
					})
					.catch(() => {});

				return true;
			}
		};
	}

	function print() {
		printRunning = true;
		const search =
			page.url.searchParams.get(tableConfig.searchParam) ??
			DataTableSearchParams.fromDataTable(tableConfig.initialState).toURL();

		API.request<any, Blob>(
			`/api/v2/courses/${page.params.courseId}/print/questions`,
			{
				method: 'POST',
				searchParams: {
					search
				}
			},
			fetch
		)
			.then((res) => {
				const url = URL.createObjectURL(res);
				window.open(url); // opens in new tab
			})
			.catch(() => {})
			.finally(() => {
				printRunning = false;
			});
	}
</script>

<div class="flex flex-col gap-8 m-8">
	<div class="flex flex-row justify-between">
		<h1 class="text-2xl">Question management</h1>
		<Button href="{base}/app/{page.params.courseId}/tutor/questions/0">{m.quesstions_add()}</Button>
	</div>
	<div>
		<DataTable data={rowItems} {rowCount} {...tableConfig} />
	</div>
	<div class="flex justify-end gap-4">
		<Button onclick={() => print()} disabled={printRunning}>
			{#if printRunning}
				<Loader></Loader>
			{/if}
			{m.print_filtered_questions()}
		</Button>
	</div>
</div>
