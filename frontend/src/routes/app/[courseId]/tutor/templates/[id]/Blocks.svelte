<script lang="ts">
	import {
		AnswerDistributionEnum,
		QuestionFormatEnum,
		type TemplateCreatorDTO,
		type TemplateCreatorResponse
	} from '$lib/api_types';
	import { Button } from '$lib/components/ui/button';
	import { type ErrorObject } from '$lib/components/ui/form/types';
	import { Input } from '$lib/components/ui/input';
	import { onMount } from 'svelte';
	import Block from './Block.svelte';
	import { API } from '$lib/services/api.svelte';
	import { toast } from 'svelte-sonner';
	import { page } from '$app/state';
	import * as Accordion from '$lib/components/ui/accordion/index.js';
	import { m } from '$lib/paraglide/messages';
	import Loader from '$lib/components/ui/loader/loader.svelte';

	let {
		form = $bindable()
	}: {
		form: any;
	} = $props();

	let blocksCount = $state(0);

	function addBlock() {
		form.fields.blocks.push({
			title: '',
			showName: true,
			difficultyFrom: 0,
			difficultyTo: 100,
			weight: 0,
			questionFormat: QuestionFormatEnum.ABCD,
			questionCount: 4,
			answerCount: 4,
			AnswerDistribution: AnswerDistributionEnum.MINIMUM_ONE_CORRECT_ONE_INCORRECT,
			wrongAnswerPercentage: 100,
			mixInsideBlock: true,
			chapters: [],
			segments: []
		});
	}

	function setBlockCount(requestedCount: number) {
		let nonDeletedBlocks = 0;

		for (let block of form.fields.blocks) {
			if (!block.deleted) {
				nonDeletedBlocks += 1;
			}
			if (nonDeletedBlocks > requestedCount) {
				block.deleted = true;
			}
		}

		while (nonDeletedBlocks < requestedCount) {
			nonDeletedBlocks += 1;
			addBlock();
		}
	}

	let isLoading = $state(true);
	let templateCreatorData: TemplateCreatorDTO[] = $state([]);
	async function LoadTemplateCreatorData(courseId: string) {
		await API.request<null, TemplateCreatorResponse>(
			`/api/v2/courses/${courseId}/templates/creator`
		)
			.then((res) => {
				templateCreatorData = res.chapters;
				isLoading = false;
			})
			.catch(() => {});
	}

	onMount(async () => {
		await LoadTemplateCreatorData(page.params.courseId);
		blocksCount = form.fields.blocks.length;
	});
</script>

{#if isLoading}
	{m.template_loading()}
	<Loader></Loader>
{:else}
	<div class="flex items-center gap-4 pb-4">
		<h1 class="text-xl text-nowrap">{m.template_blocks()}</h1>
		<div class="flex-1"></div>
		<div class="flex flex-col">
			<Input
				type="number"
				bind:value={blocksCount}
				onchange={() => {
					if (blocksCount < 0) {
						blocksCount = 0;
					}
				}}
			></Input>
			{#if form.errors.blocks && typeof form.errors.blocks == 'string'}
				<p class="text-sm text-red-500">{form.errors.blocks}</p>
			{/if}
		</div>
		<Button
			variant="outline"
			onclick={() => {
				setBlockCount(blocksCount);
			}}>Set</Button
		>
	</div>
	<!-- TODO fix any -->
	{#if form.fields.blocks.filter((b: any) => !b.deleted).length > 0}
		<Accordion.Root type="single" class="w-full border rounded-md">
			{#each form.fields.blocks as block, index}
				{@const err = (form?.errors?.blocks as ErrorObject) ?? {}}
				{@const blockErr = err[index] ?? {}}

				{#if !block.deleted}
					<Accordion.Item value={String(index)}>
						<Accordion.Trigger class="items-center px-4 py-2">
							<div class="flex items-center justify-between w-full">
								<p>{m.template_block()} {index + 1}</p>
								<Button
									variant="destructive"
									onclick={() => {
										form.fields.blocks.splice(index, 1);
									}}>{m.delete()}</Button
								>
							</div>
						</Accordion.Trigger>
						<Accordion.Content class="flex flex-col gap-2 p-4 border-t rounded-none">
							<Block bind:fields={form.fields.blocks[index]} errors={blockErr} {templateCreatorData}
							></Block>
						</Accordion.Content>
					</Accordion.Item>
				{/if}
			{/each}
		</Accordion.Root>
	{/if}
{/if}
