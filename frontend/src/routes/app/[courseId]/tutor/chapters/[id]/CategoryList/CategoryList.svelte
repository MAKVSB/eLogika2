<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { tableConfig } from './schema';
	import { API } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import type { CategoryListItemDTO, CategoryListResponse } from '$lib/api_types';
	import { m } from '$lib/paraglide/messages';
	import { base } from '$app/paths';
	import { DataTableSearchParams } from '$lib/api_types_static';

	let rowItems: CategoryListItemDTO[] = $state([]);
	let rowCount: number = $state(0);

	$effect(() => {
		const search =
			page.url.searchParams.get(tableConfig.searchParam) ??
			DataTableSearchParams.fromDataTable(tableConfig.initialState).toURL();

		API.request<null, CategoryListResponse>(
			`/api/v2/courses/${page.params.courseId}/chapters/${page.params.id}/categories`,
			{
				searchParams: {
					...(search ? { search } : {})
				}
			}
		)
			.then((res) => {
				rowItems = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {});
	});
</script>

<div class="flex flex-row justify-between">
	<h1 class="mb-8 text-2xl">Categories</h1>
	<div class="flex gap-2">
		<Button href="{base}/app/{page.params.courseId}/tutor/categories/0">{m.category_add()}</Button>
	</div>
</div>

<DataTable data={rowItems} {rowCount} {...tableConfig} />
