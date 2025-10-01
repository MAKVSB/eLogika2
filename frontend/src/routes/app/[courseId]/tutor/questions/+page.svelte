<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import { API } from '$lib/services/api.svelte';
	import type {
		QuestionCheckResponse,
		QuestionListItemDTO,
		QuestionToggleActiveResponse
	} from '$lib/api_types';
	import { type InitialTableState } from '@tanstack/table-core';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import { m } from '$lib/paraglide/messages';
	import { invalidateAll } from '$app/navigation';
	import { base } from '$app/paths';
	import { FilterTypeEnum } from '$lib/components/ui/data-table/filter';

	let loading: boolean = $state(true);
	let rowItems: QuestionListItemDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({
		columnVisibility: {
			chapterId: false,
			categoryId: false
		}
	});

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

	$effect(() => {
		data.chapterData
			.then((res) => {
				const filter = filters.find((f) => f.accessorKey == 'chapterId');
				if (filter && filter.type == FilterTypeEnum.SELECT) {
					filter.values = res.items.map((chapter) => {
						return {
							value: chapter.id,
							display: chapter.name
						};
					});
				}
			})
			.catch(() => {})
			.finally(() => {
				loading = false;
			});
	});

	$effect(() => {
		data.categoryData
			.then((res) => {
				const filter = filters.find((f) => f.accessorKey == 'categoryId');
				if (filter && filter.type == FilterTypeEnum.SELECT) {
					filter.values = res.items.map((chapter) => {
						return {
							value: chapter.id,
							display: chapter.name
						};
					});
				}
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
						if (!confirm('Question will be deleted permanently.')) {
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

	const activeColumn = columns.find((c) => c.uniqueId == 'active');
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
		const search = page.url.searchParams.get('search');

		API.request<any, Blob>(
			`/api/v2/courses/${page.params.courseId}/print/questions`,
			{
				method: 'POST',
				searchParams: {
					...(search ? { search: search } : {})
				}
			},
			fetch
		)
			.then((res) => {
				const url = URL.createObjectURL(res);
				window.open(url); // opens in new tab
			})
			.catch(() => {});
	}
</script>

<div class="flex flex-col gap-8 m-8">
	<div class="flex flex-row justify-between">
		<h1 class="text-2xl">Question management</h1>
		<Button href="{base}/app/{page.params.courseId}/tutor/questions/0">{m.quesstions_add()}</Button>
	</div>
	<div>
		{#if !loading}
			<DataTable
				data={rowItems}
				{columns}
				{filters}
				{initialState}
				{rowCount}
				queryParam="search"
			/>
		{/if}
	</div>
	<div class="flex justify-end gap-4">
		<Button onclick={() => print()}>{m.print_filtered_questions()}</Button>
	</div>
</div>
