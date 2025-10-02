<script lang="ts">
	import { type CourseItemResultsDTO } from '$lib/api_types';
	import * as Table from '$lib/components/ui/table/index.js';
	import { toast } from 'svelte-sonner';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import { m } from '$lib/paraglide/messages';
	import StudentRow from './StudentRow.svelte';

	let results = $state<CourseItemResultsDTO[]>([]);

	let { data } = $props();
	$effect(() => {
		if (data.results) {
			data.results
				.then((data) => (results = data.data))
				.catch((err) => {
					console.error(err);
					toast.error('Failed to load question');
				});
		}
	});
</script>

{#snippet TableHeader()}
	<Table.Row>
		<Table.Head></Table.Head>
		<Table.Head>{m.user_username()}</Table.Head>
		<Table.Head>{m.user_full_name()}</Table.Head>
		<Table.Head>{m.points_real()}</Table.Head>
		<Table.Head>{m.course_item_passed()}</Table.Head>
	</Table.Row>
{/snippet}

<div class="flex flex-col gap-8 m-8">
	{#await data.results}
		<Pageloader></Pageloader>
	{:then staticResourceData}
		<div class="flex flex-row justify-between">
			<h1 class="text-2xl">
				Course item results:
			</h1>
		</div>
		<Table.Root>
			<Table.Header class="font-medium border-t bg-muted/50">
				{@render TableHeader()}
			</Table.Header>
			<Table.Body>
				{#if staticResourceData.data.length === 0}
					<Table.Row>
						<Table.Cell colspan={7}>{m.no_items_found()}</Table.Cell>
					</Table.Row>
				{:else}
					{#each staticResourceData.data as item}
						<StudentRow userData={item}></StudentRow>
					{/each}
				{/if}
			</Table.Body>
			{#if staticResourceData.data.length > 10}
					<Table.Footer>
						{@render TableHeader()}
					</Table.Footer>
				{/if}
		</Table.Root>
	{/await}
</div>
