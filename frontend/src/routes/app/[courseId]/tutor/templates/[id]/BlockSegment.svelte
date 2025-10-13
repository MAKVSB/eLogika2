<script lang="ts">
	import { CategoryFilterEnum, StepSelectionEnum, type TemplateCreatorDTO } from '$lib/api_types';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import * as Form from '$lib/components/ui/form';
	import { Label } from '$lib/components/ui/label';
	import { m } from '$lib/paraglide/messages';
	import { enumToOptions } from '$lib/utils';
	import { formatError } from 'zod/v4';
	import QuestionSelectModal from './QuestionSelectModal.svelte';

	let {
		fields = $bindable(),
		errors = $bindable({}),
		templateCreatorData,
		delete: deleteChapter
	}: {
		fields: any;
		errors?: any;
		templateCreatorData: TemplateCreatorDTO[];
		delete: any;
	} = $props();

	let questionSelectorOpen = $state(false);
	let changeCount = $state(0);

	let chapterOptions: Form.SelectOptions = $derived.by(() => {
		return templateCreatorData.map((v) => {
			return {
				display: v.name,
				value: v.id
			};
		});
	});

	let categoryOptions: Form.SelectOptions = $derived.by(() => {
		const chapterData = templateCreatorData.find((chap) => chap.id == fields.chapterId);
		if (chapterData) {
			return chapterData.categories.map((v) => {
				return {
					display: v.name,
					value: v.id
				};
			});
		} else {
			return [];
		}
	});

	let categorySteps = $derived.by(() => {
		const chapterData = templateCreatorData.find((chap) => chap.id == fields.chapterId);
		if (chapterData) {
			const categoryData = chapterData.categories.find((cat) => cat.id == fields.categoryId);
			if (categoryData) {
				return categoryData.steps;
			} else {
				return [];
			}
		} else {
			return [];
		}
	});

	let filterByOptions: Form.SelectOptions = $derived.by(() => {
		return Object.values(CategoryFilterEnum).map((value) => {
			let enabled = true;

			if (!fields.categoryId && [CategoryFilterEnum.S, CategoryFilterEnum.SQOR].includes(value)) {
				enabled = false;
			}
			return {
				value: value,
				display: m.block_category_filter_steps_enum({ value: value }),
				enabled,
				...(!enabled ? { tooltip: 'Select chapter and category to enable this option' } : {})
			};
		});
	});

	function onChangeChapter() {
		if (fields.chapterId == null) {
			fields.categoryId = null;
		} else {
			const chapterData = templateCreatorData.find((chap) => chap.id == fields.chapterId);
			const categoryData = chapterData?.categories.find((cat) => cat.id == fields.categoryId);
			if (!categoryData) {
				fields.categoryId = null;
			}
		}
		onChangeCategory();
	}

	function onChangeCategory() {
		if (fields.categoryId == null) {
			if ([CategoryFilterEnum.ALL, CategoryFilterEnum.S].includes(fields.filterBy)) {
				fields.filterBy = CategoryFilterEnum.ALL;
			}
			if ([CategoryFilterEnum.Q, CategoryFilterEnum.SQOR].includes(fields.filterBy)) {
				fields.filterBy = CategoryFilterEnum.Q;
			}
		} else {
			if (fields.steps.length != 0) {
				const chapterData = templateCreatorData.find((chap) => chap.id == fields.chapterId);
				const categoryData = chapterData?.categories.find((cat) => cat.id == fields.categoryId);
				const firstStepData = categoryData?.steps.find((s) => s.id == fields.steps[0]);
				if (!firstStepData) {
					fields.steps = [];
				}
			}
		}
		onChangeFilterBy();
	}

	function onChangeFilterBy() {
		switch (fields.filterBy) {
			case CategoryFilterEnum.ALL:
				fields.steps = [];
				fields.questions = [];
				break;
			case CategoryFilterEnum.Q:
				fields.steps = [];
				break;
			case CategoryFilterEnum.S:
				fields.questions = [];
				if (fields.steps.length != 0) {
					const chapterData = templateCreatorData.find((chap) => chap.id == fields.chapterId);
					const categoryData = chapterData?.categories.find((cat) => cat.id == fields.categoryId);
					const firstStepData = categoryData?.steps.find((s) => s.id == fields.steps[0]);
					if (!firstStepData) {
						fields.steps = [];
					}
				}
				break;
			default:
				break;
		}
	}
</script>

<div class="flex flex-row border">
	<div class="grid w-full grid-cols-1 p-2 md:grid-cols-2 xl:grid-cols-4 gap-x-4 gap-y-2">
		<Form.TextInput
			bind:value={fields.questionCount}
			error={errors.questionCount}
			id="questionCount"
			name="questionCount"
			type="number"
			title={m.template_block_chapter_questioncount()}
		></Form.TextInput>

		<Form.SingleSelect
			title={m.filter_chapter()}
			name="chapter"
			id="chapter"
			bind:value={fields.chapterId}
			placeholder={m.template_segment_chapter_select_placeholder()}
			onchange={onChangeChapter}
			error=""
			nullable
			options={chapterOptions}
		></Form.SingleSelect>
		<Form.SingleSelect
			title={m.filter_category()}
			name="chapter"
			id="chapter"
			bind:value={fields.categoryId}
			placeholder={m.template_segment_category_select_placeholder()}
			error=""
			nullable
			onchange={onChangeCategory}
			disabled={!fields.chapterId}
			options={categoryOptions}
		></Form.SingleSelect>

		<Form.SingleSelect
			title={m.template_block_chapter_category_filterby()}
			name="filterBy"
			id="filterBy"
			bind:value={fields.filterBy}
			error={errors.filterBy}
			onchange={onChangeFilterBy}
			options={filterByOptions}
		></Form.SingleSelect>


        <hr class="col-span-1 md:col-span-2 xl:col-span-4" />
        

		{#if [CategoryFilterEnum.S, CategoryFilterEnum.SQOR].includes(fields.filterBy)}
			<div class="grid gap-2">
				<Label>{m.template_block_chapter_category_steps()}</Label>
				<ul class="grid gap-2">
					{#key changeCount}
						{#each categorySteps as step}
							<li class="flex gap-4">
								<Checkbox
									id="step-{step.id}"
									checked={fields.steps.indexOf(step.id) !== -1}
									onclick={() => {
										const found = fields.steps.indexOf(step.id);
										if (found == -1) {
											fields.steps.push(step.id);
										} else {
											fields.steps.splice(found, 1);
										}
										changeCount += 1;
									}}
								></Checkbox>
								<Label for="step-{step.id}">{step.name} {step.id} ({step.difficulty})</Label>
							</li>
						{/each}
					{/key}
				</ul>
			</div>
			<Form.SingleSelect
				title={m.template_block_chapter_category_stepsmode()}
				name="stepsMode"
				id="stepsMode"
				bind:value={fields.stepsMode}
				error={errors.stepsMode}
				options={enumToOptions(StepSelectionEnum, m.conjunctions_enum)}
			></Form.SingleSelect>
		{:else}
			<div class="md:col-span-2"></div>
		{/if}
		{#if [CategoryFilterEnum.Q, CategoryFilterEnum.SQOR].includes(fields.filterBy)}
			<div class="flex flex-col gap-2 md:col-span-2">
				<Label>&nbsp;</Label>
				<Button
					class="w-full"
					onclick={() => {
						questionSelectorOpen = true;
					}}>{m.template_block_chapter_category_question_select_button()}</Button
				>
				{#if questionSelectorOpen}
					<QuestionSelectModal bind:isOpen={questionSelectorOpen} bind:selectedQuestions={fields.questions}></QuestionSelectModal>
				{/if}
			</div>
		{:else}
			<div class="md:col-span-2"></div>
		{/if}
	</div>
    <Button class="m-2" variant="destructive" onclick={deleteChapter}>{m.delete()}</Button>
</div>
