<script lang="ts">
	import type { CourseItemResultDTO } from '$lib/api_types';
	import * as Table from '$lib/components/ui/table/index.js';
	import { m } from '$lib/paraglide/messages';

	let {
		results,
		studentId
	}: {
		results: CourseItemResultDTO[];
		studentId: number;
	} = $props();
</script>

{#snippet TableHeader()}
	<Table.Row>
		<Table.Head>Course item name</Table.Head>
		<Table.Head>Term name</Table.Head>
		<Table.Head>{m.points_real()}</Table.Head>
		<!-- TODO <Table.Head>{m.course_item_passed()}</Table.Head> -->
		<Table.Head>Selected</Table.Head>
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
						{result.points}
						{#if !result.final}
							*
						{/if}
					</Table.Cell>
					<!-- <Table.Cell>{m.yes_no({ value: String(result.passed) })}</Table.Cell> -->
					<Table.Cell>{m.yes_no({ value: String(result.selected) })}</Table.Cell>
				</Table.Row>
			{/each}
		{/if}
	</Table.Body>
</Table.Root>
