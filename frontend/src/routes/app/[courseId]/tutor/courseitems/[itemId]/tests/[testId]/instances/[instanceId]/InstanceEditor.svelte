<script lang="ts">
	import * as Table from '$lib/components/ui/table';
	import * as Form from '$lib/components/ui/form';
	import { toast } from 'svelte-sonner';
	import { API } from '$lib/services/api.svelte';
	import {
		type TestInstanceDTO,
		CourseUserRoleEnum,
		QuestionFormatEnum,
		TestInstanceStateEnum
	} from '$lib/api_types';
	import { page } from '$app/state';
	import { Label } from '$lib/components/ui/label';
	import DateRangeField from '$lib/components/ui/date-range-field/date-range-field.svelte';
	import { displayUserName, enumToOptions } from '$lib/utils';
	import TiptapRenderer from '$lib/components/tiptap-renderer/tiptap-renderer.svelte';
	import Tiptap from '$lib/components/tiptap/Tiptap.svelte';
	import { m } from '$lib/paraglide/messages';
	import { parseAbsoluteToLocal } from '@internationalized/date';
	import GlobalState from '$lib/shared.svelte';

	let {
		instanceData,
		editable
	}: {
		instanceData: TestInstanceDTO;
		editable: boolean;
	} = $props();

	let formErrors = $state();

	let form = {
		errors: {} as Form.ErrorObject,
		isSubmitting: false
	};

	async function handleSubmitCutomValidation() {
		if (!instanceData) {
			return;
		}
		if (instanceData.state != TestInstanceStateEnum.FINISHED) {
			if (!confirm('Saving test instance will transfer in into "finished" state')) {
				return;
			}
		}
		// TODO validation

		await API.request<any, any>(
			`/api/v2/courses/${page.params.courseId}/tests/${page.params.itemId}/instance/${instanceData.id}/tutorsave`,
			{
				method: 'PUT',
				body: {
					id: instanceData.id,
					questions: (instanceData.questions ?? []).map((q) => {
						return {
							id: q.id,
							textAnswer: q.textAnswer,
							textAnswerReviewed: q.textAnswerReviewed,
							textAnswerPercentage: Number(q.textAnswerPercentage),
							answers: q.answers.map((a) => {
								return {
									id: a.id,
									selected: a.selected
								};
							})
						};
					}),
					bonusPoints: instanceData.bonusPoints,
					bonusPointsReason: instanceData.bonusPointsReason
				}
			}
		)
			.then((res) => {})
			.catch(() => {});
	}

	const getVariantLabel = (n: number) => {
		let label = '';
		while (n >= 0) {
			label = String.fromCharCode('a'.charCodeAt(0) + (n % 26)) + label;
			n = Math.floor(n / 26) - 1;
		}
		return label;
	};

	let showCorrect = $state(false);

	let timeRange = $derived({
		start: parseAbsoluteToLocal(instanceData.startedAt),
		end: parseAbsoluteToLocal(instanceData.endedAt ?? instanceData.endsAt)
	});

	let maxAnswerCount = $derived.by(() => {
		let answerCount = 1;

		for (const question of instanceData.questions ?? []) {
			if (question.questionFormat == QuestionFormatEnum.ABCD) {
				if (question.answers.length > answerCount) {
					answerCount = question.answers.length;
				}
			}
		}
		return answerCount;
	});

	let hasTitle: boolean = $derived.by(() => {
		for (const question of instanceData.questions ?? []) {
			if (question.title) {
				return true;
			}
		}
		return false;
	});

	let scrollToQuestion = (id: number) => {
		const el = document.getElementById("q" + id);
		if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' });
	}
</script>

<div class="flex flex-col gap-4">
	<Form.TextInput
		title={m.testinstance_participants()}
		id="participant"
		name="participant"
		value="{displayUserName(instanceData.participant)} ({instanceData.participant.username})"
		disabled
		error=""
	></Form.TextInput>
	<div class="flex flex-col gap-2">
		<Label>{m.testinstance_activetime()}</Label>
		<DateRangeField value={timeRange} granularity="minute" disabled></DateRangeField>
	</div>
	<Form.SingleSelect
		title={m.testinstance_status()}
		id="state"
		name="state"
		bind:value={instanceData.state}
		options={enumToOptions(TestInstanceStateEnum, m.test_instance_state_enum)}
		error=""
		disabled
	></Form.SingleSelect>
	<Form.TextInput
		title={m.testinstance_variant()}
		id="group"
		name="group"
		bind:value={instanceData.group}
		disabled
		error=""
	></Form.TextInput>
	<div class="grid gap-4 sm:grid-cols-2">
		<div class="flex flex-col">
			<Form.TextInput
				title={m.testinstance_points()}
				id="points"
				name="points"
				type="number"
				innerClass="text-4xl sm:text-4xl md:text-4xl lg:text-4xl xl:text-4xl {instanceData.points + instanceData.bonusPoints >= instanceData.pointsMin ? "text-green-500 disabled:opacity-100" : "text-red-500 disabled:opacity-100"}"
				value={instanceData.points + instanceData.bonusPoints}
				disabled
				error=""
			></Form.TextInput>
			{#if !instanceData.pointsFinal}
				<span class="text-red-500">Points are marked as hidden or not final</span>
			{/if}
		</div>
		<Form.TextInput
			title={m.testinstance_bonuspoints()}
			id="bonusPoints"
			name="bonusPoints"
			type="number"
			innerClass="text-4xl sm:text-4xl md:text-4xl lg:text-4xl xl:text-4xl {instanceData.points + instanceData.bonusPoints >= instanceData.pointsMin ? "text-green-500 disabled:opacity-100" : "text-red-500 disabled:opacity-100"}"
			bind:value={instanceData.bonusPoints}
			disabled={!editable}
			error=""
		></Form.TextInput>
		<Form.TextArea
			title={m.testinstance_bonuspoints_reason()}
			id="bonusPointsReason"
			name="bonusPointsReason"
			class="sm:col-span-2"
			bind:value={instanceData.bonusPointsReason}
			disabled={!editable}
			error=""
		></Form.TextArea>
	</div>
	{#if instanceData.showCorrectness}
		<Form.Checkbox
			title={m.correctanswers_show()}
			name="showcorrect"
			id="showcorrect"
			bind:value={showCorrect}
			error=""
		></Form.Checkbox>
	{/if}
</div>
<Form.Root onsubmit={handleSubmitCutomValidation}>
	{#if instanceData.questions}
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Cell>{m.question()}</Table.Cell>
					{#if hasTitle}
						<Table.Cell>
							{m.question_title()}
						</Table.Cell>
					{/if}
					{#each { length: maxAnswerCount } as _, i}
						<Table.Cell>
							{getVariantLabel(i)}
						</Table.Cell>
					{/each}
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each instanceData.questions ?? [] as question}
					<Table.Row>
						<Table.Cell onclick={() => scrollToQuestion(question.id)}>
							{question.order + 1}
						</Table.Cell>
						{#if hasTitle}
							<Table.Cell onclick={() => scrollToQuestion(question.id)}>
								{question.title}
							</Table.Cell>
						{/if}
						{#if question.questionFormat == QuestionFormatEnum.OPEN}
							<Table.Cell colspan={maxAnswerCount}>
								<Form.TextInput
									name="q{question.id}"
									id="q{question.id}"
									bind:value={question.textAnswerPercentage}
									disabled={!editable}
									onchange={() => {
										question.textAnswerReviewed = true;
									}}
									error=""
								></Form.TextInput>
							</Table.Cell>
						{:else if question.questionFormat == QuestionFormatEnum.ABCD}
							{#each question.answers as answer}
								<Table.Cell>
									<Form.Checkbox
										innerClass={showCorrect
											? answer.correct
												? 'bg-green-500 data-[state=checked]:bg-green-500 dark:bg-green-500 dark:data-[state=checked]:bg-green-500'
												: 'bg-red-500 data-[state=checked]:bg-red-500 dark:bg-red-500 dark:data-[state=checked]:bg-red-500'
											: ''}
										name="q{question.id}-a{answer.id}"
										id="q{question.id}-a{answer.id}"
										bind:value={answer.selected}
										disabled={!editable}
										error=""
									></Form.Checkbox>
								</Table.Cell>
							{/each}
						{:else}
						<Table.Cell>
							Invalid question type
						</Table.Cell>
						{/if}
					</Table.Row>
				{:else}
					<Table.Row>
						<Table.Cell>{m.no_questions()}</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	{/if}
	{#if instanceData.showContent}
		<div>
			{#each instanceData.questions ?? [] as question}
				<div class="flex flex-col gap-4 p-4 border" id={"q" + question.id}>
					<div>
						<h2 class="text-xl">
							{m.question()}
							{question.order + 1}
						</h2>
						{#if question.content}
							<TiptapRenderer jsonContent={question.content}></TiptapRenderer>
						{/if}
					</div>
					{#if question.questionFormat == QuestionFormatEnum.ABCD}
						<div>
							<h2 class="text-xl">{m.answers()}</h2>
							<Table.Root>
								<Table.Body>
									{#each question.answers as answer}
										<Table.Row>
											<Table.Cell style="width: 60px;">
												<Form.Checkbox
													innerClass={showCorrect
														? answer.correct
															? 'bg-green-500 data-[state=checked]:bg-green-500 dark:bg-green-500 dark:data-[state=checked]:bg-green-500'
															: 'bg-red-500 data-[state=checked]:bg-red-500 dark:bg-red-500 dark:data-[state=checked]:bg-red-500'
														: ''}
													name="q{question.id}-a{answer.id}"
													id="q{question.id}-a{answer.id}"
													bind:value={answer.selected}
													disabled={!editable}
													error=""
												></Form.Checkbox>
											</Table.Cell>
											<Table.Cell>
												<TiptapRenderer jsonContent={answer.content}></TiptapRenderer>
											</Table.Cell>
										</Table.Row>
									{/each}
								</Table.Body>
							</Table.Root>
						</div>
					{:else if question.questionFormat == QuestionFormatEnum.OPEN}
						<div>
							<h2 class="text-xl">{m.answer()}</h2>
							<Tiptap value={question.textAnswer} disabled></Tiptap>
						</div>
						<div>
							<Form.TextInput
								title={m.testinstance_tutor_answer_percentage()}
								name="q{question.id}"
								id="q{question.id}"
								bind:value={question.textAnswerPercentage}
								disabled={!editable}
								onchange={() => {
									question.textAnswerReviewed = true;
								}}
								error=""
							></Form.TextInput>
						</div>
						{#if showCorrect}
							<div>
								<h2 class="text-xl">Correct answers</h2>
								{#each question.openAnswers?.filter((a) => a.correct) ?? [] as correctAnswer}
									<Tiptap value={correctAnswer.content} disabled></Tiptap>
								{/each}
							</div>
							<div>
								<h2 class="text-xl">Incorrect answers</h2>
								{#each question.openAnswers?.filter((a) => !a.correct) ?? [] as correctAnswer}
									<Tiptap value={correctAnswer.content} disabled></Tiptap>
								{/each}
							</div>
						{/if}
					{/if}
				</div>
			{/each}
		</div>
	{/if}
	{#if GlobalState.activeRole == CourseUserRoleEnum.ADMIN || GlobalState.activeRole == CourseUserRoleEnum.TUTOR || GlobalState.activeRole == CourseUserRoleEnum.GARANT}
		<Form.Button text={m.save()} textSubmiting={m.save_progress()}></Form.Button>
	{/if}
</Form.Root>
