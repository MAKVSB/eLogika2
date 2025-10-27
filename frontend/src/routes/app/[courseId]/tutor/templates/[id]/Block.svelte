<script lang="ts">
	import {
		AnswerDistributionEnum,
		CategoryFilterEnum,
		QuestionFormatEnum,
		StepSelectionEnum,
		type TemplateCreatorDTO
	} from '$lib/api_types';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Form from '$lib/components/ui/form';
	import { type ErrorObject } from '$lib/components/ui/form/types';
	import { m } from '$lib/paraglide/messages';
	import { enumToOptions } from '$lib/utils';
	import BlockSegment from './BlockSegment.svelte';

	let {
		fields = $bindable(),
		errors = $bindable({}),
		templateCreatorData
	}: {
		fields: any;
		errors?: any;
		templateCreatorData: TemplateCreatorDTO[];
	} = $props();

	function addSegment() {
		fields.segments.push({
			questionCount: 1,
			chapterId: null,
			categoryId: null,
			steps: [],
			questions: [],
			stepsMode: StepSelectionEnum.NC,
			filterBy: CategoryFilterEnum.ALL,
			deleted: false
		});

		console.log(fields.segments);
	}

	let questionCountDeriver = $derived.by(() => {
		if (fields.segments.length != 0) {
			return fields.segments.reduce((acc: number, item: any) => {
				const n = Number(item.questionCount);
				if (isNaN(n)) {
					return acc;
				} else {
					return acc + n;
				}
			}, 0);
		} else {
			return fields.questionCount;
		}
	});

	$effect(() => {
		fields.questionCount = questionCountDeriver;
	});
</script>

<div class="grid grid-cols-12 max-w-160 gap-x-4 gap-y-2">
	<Form.TextInput
		title={m.template_block_title()}
		name="title"
		id="title"
		type="text"
		class="col-span-8 sm:col-span-8"
		bind:value={fields.title}
		error={errors.title}
	></Form.TextInput>
	<Form.Checkbox
		title={m.template_block_showname()}
		name="showName"
		id="showName"
		class="col-span-4 sm:col-span-4"
		bind:value={fields.showName}
		error={errors.showName}
	></Form.Checkbox>

	<Form.TextInput
		title={m.template_block_difficultyfrom()}
		name="difficultyFrom"
		id="difficultyFrom"
		type="number"
		class="col-span-6 sm:col-span-6"
		bind:value={fields.difficultyFrom}
		error={errors.difficultyFrom}
	></Form.TextInput>
	<Form.TextInput
		title={m.template_block_difficultyto()}
		name="difficultyTo"
		id="difficultyTo"
		type="number"
		class="col-span-6 sm:col-span-6"
		bind:value={fields.difficultyTo}
		error={errors.difficultyTo}
	></Form.TextInput>

	<Form.TextInput
		title={m.template_block_weight()}
		name="weight"
		id="weight"
		type="number"
		class="col-span-12 sm:col-span-12"
		bind:value={fields.weight}
		error={errors.weight}
	></Form.TextInput>

	<Form.TextInput
		title={m.template_block_questioncount()}
		name="questionCount"
		id="questionCount"
		type="number"
		class="col-span-4 sm:col-span-4"
		bind:value={fields.questionCount}
		disabled
		error={errors.questionCount}
	></Form.TextInput>
	<Form.SingleSelect
		title={m.template_block_questionformat()}
		name="questionFormat"
		id="questionFormat"
		class="col-span-8 sm:col-span-8"
		bind:value={fields.questionFormat}
		error={errors.questionFormat}
		options={enumToOptions(QuestionFormatEnum, m.question_format_enum)}
	></Form.SingleSelect>

	<Form.TextInput
		title={m.template_block_answercount()}
		name="answerCount"
		id="answerCount"
		type="number"
		class="col-span-4 sm:col-span-4"
		bind:value={fields.answerCount}
		error={errors.answerCount}
		disabled={fields.questionFormat == QuestionFormatEnum.OPEN}
	></Form.TextInput>
	<Form.SingleSelect
		title={m.template_block_answerdistribution()}
		name="answerDistribution"
		id="answerDistribution"
		class="col-span-8 sm:col-span-8"
		bind:value={fields.answerDistribution}
		error={errors.answerDistribution}
		options={enumToOptions(AnswerDistributionEnum, m.answer_distribution_enum)}
		disabled={fields.questionFormat == QuestionFormatEnum.OPEN}
	></Form.SingleSelect>

	<Form.TextInput
		title={m.template_block_penalisationpercentage()}
		name="wrongAnswerPercentage"
		id="wrongAnswerPercentage"
		type="number"
		class="col-span-12 sm:col-span-6"
		bind:value={fields.wrongAnswerPercentage}
		error={errors.wrongAnswerPercentage}
	></Form.TextInput>
	<Form.Checkbox
		title={m.template_block_allowemptyanswers()}
		tooltip={m.template_block_allowemptyanswers_tooltip()}
		name="allowEmptyAnswers"
		id="allowEmptyAnswers"
		type="number"
		class="col-span-12 sm:col-span-6"
		bind:value={fields.allowEmptyAnswers}
		error={errors.allowEmptyAnswers}
	></Form.Checkbox>

	<Form.Checkbox
		title={m.template_block_mix()}
		name="mixInsideBlock"
		id="mixInsideBlock"
		class="col-span-4 sm:col-span-4"
		bind:value={fields.mixInsideBlock}
		error={errors.mixInsideBlock}
	></Form.Checkbox>
</div>
<hr class="my-2" />

{#if fields.segments.filter((c: any) => !c.deleted).length > 0}
	<p class="text-lg font-bold">{m.template_block_segments()}:</p>
	{#each fields.segments as segment, index}
		{#if !segment.deleted}
			{@const err = errors.segments ?? ({} as ErrorObject)}
			<BlockSegment
				bind:fields={fields.segments[index]}
				{templateCreatorData}
				errors={err[index] ?? ({} as ErrorObject)}
				delete={() => {
					fields.segments.splice(index, 1);
				}}
			></BlockSegment>
		{/if}
	{/each}
{:else}
	<p class="text-lg font-bold text-red-500">
		{m.template_block_nochapter()}
	</p>
{/if}

<hr class="my-2" />
<div class="flex justify-between">
	<Button onclick={() => addSegment()}>{m.template_block_segment_add()}</Button>
</div>
