<script lang="ts">
	import { invalidate } from '$app/navigation';
	import { page } from '$app/state';
	import type { CourseItemResultDTO, CourseItemSelectResultResponse } from '$lib/api_types';
	import { Button } from '$lib/components/ui/button';
	import * as Table from '$lib/components/ui/table/index.js';
	import { m } from '$lib/paraglide/messages';
	import { getLocale } from '$lib/paraglide/runtime';
	import { API } from '$lib/services/api.svelte';

	let {
		results,
		studentId
	}: {
		results: CourseItemResultDTO[];
		studentId: number;
	} = $props();

	const toggleSelect = (itemId: number, resultId: number) => {
		API.request<null, CourseItemSelectResultResponse>(
			`/api/v2/courses/${page.params.courseId}/items/${itemId}/results/${resultId}`,
			{
				method: 'PUT'
			}
		).then(() => {
			invalidate((url) => {
				// Match /api/v2/courses/<anything>/items/<anything>/results
				return /^\/api\/v2\/courses\/[^/]+\/items\/[^/]+\/results$/.test(url.pathname);
			});
		});
	};
</script>

{#snippet TableHeader()}
	<Table.Row>
		<Table.Head>Course item name</Table.Head>
		<Table.Head>{m.term_name()}</Table.Head>
		<Table.Head>StartedAt</Table.Head>
		<Table.Head>{m.points_real()}</Table.Head>
		<!-- TODO <Table.Head>{m.course_item_passed()}</Table.Head> -->
		<Table.Head>{m.result_selected()}</Table.Head>
		<Table.Head>{m.actions()}</Table.Head>
	</Table.Row>
{/snippet}

<Table.Root>
	<Table.Header class="font-medium border-t bg-muted/50">
		{@render TableHeader()}
	</Table.Header>
	<Table.Body>
		{#if results.length === 0}
			<Table.Row>
				<Table.Cell colspan={7}>No items found</Table.Cell>
			</Table.Row>
		{:else}
			{#each results as result}
				<Table.Row>
					<Table.Cell>{result.courseItemName}</Table.Cell>
					<Table.Cell>{result.termName}</Table.Cell>
					<Table.Cell>
						{new Date(result.startedAt).toLocaleString(getLocale())}
					</Table.Cell>
					<Table.Cell>
						{result.points}
						{#if !result.final}
							*
						{/if}
					</Table.Cell>
					<!-- <Table.Cell>{m.yes_no({ value: String(result.passed) })}</Table.Cell> -->
					<Table.Cell>{m.yes_no({ value: String(result.selected) })}</Table.Cell>
					<Table.Cell>
						{#if result.selected}
							<Button
								variant="destructive"
								onclick={() => toggleSelect(result.courseItemId, result.id)}
							>
								Unselect
							</Button>
						{:else}
							<Button onclick={() => toggleSelect(result.courseItemId, result.id)}>Select</Button>
						{/if}
					</Table.Cell>
				</Table.Row>
			{/each}
		{/if}
	</Table.Body>
</Table.Root>
