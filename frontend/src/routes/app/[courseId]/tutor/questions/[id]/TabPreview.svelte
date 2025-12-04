<script lang="ts">
	import { type QuestionGetByIdResponse } from '$lib/api_types';
	import TiptapRenderer from '$lib/components/tiptap-renderer/tiptap-renderer.svelte';
	import PageLoader from '$lib/components/ui/loader/pageloader.svelte';
	import { m } from '$lib/paraglide/messages';
	import * as Table from '$lib/components/ui/table/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox';

	let {
		question
	}: {
		question?: Promise<QuestionGetByIdResponse>;
	} = $props();
</script>

<div class="text-center text-red-500">{m.question_preview_warning()}</div>

<div class="m-4">
	{#if question}
		{#await question}
			<PageLoader></PageLoader>
		{:then questionData} 
			<h2 class="text-2xl">{m.question()}:</h2>
			<div class="p-2 border">
				<TiptapRenderer jsonContent={questionData.data.content!}></TiptapRenderer>
			</div>
			<h2 class="text-2xl">{m.answers()}:</h2>
			<div class="flex flex-col gap-4">

			<Table.Root>
				<Table.Body>
					{#each questionData.data.answers as answer}
						<Table.Row>
							<Table.Cell>
								<div class="p-4">
									<Checkbox disabled></Checkbox>
								</div>
							</Table.Cell>

							<Table.Cell class="w-full">
								<TiptapRenderer jsonContent={answer.content!}></TiptapRenderer>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
			</div>
		{/await}
	{/if}
</div>