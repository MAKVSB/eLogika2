<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { tableConfig } from './schema';
	import Button from '$lib/components/ui/button/button.svelte';
	import type { CourseListItemDTO } from '$lib/api_types';
	import { base } from '$app/paths';

	let { data } = $props();

	$effect(() => {
		data.courses
			.then((res) => {
				rowItems = res.items;
				rowCount = res.itemsCount;
			})
			.catch(() => {});
	});

	let rowItems: CourseListItemDTO[] = $state([]);
	let rowCount: number = $state(0);

</script>

<div class="m-8">
	<div class="flex flex-row justify-between">
		<h1 class="mb-8 text-2xl">Courses management</h1>
		<div class="flex gap-2">
			<Button href="{base}/app/admin/courses/0">Add course</Button>
		</div>
	</div>
	<DataTable data={rowItems} {rowCount} {...tableConfig}/>
</div>
