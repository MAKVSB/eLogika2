<script lang="ts">
	import * as Table from '$lib/components/ui/table';
	import * as Form from '$lib/components/ui/form';
	import { toast } from 'svelte-sonner';
	import { API, ApiError } from '$lib/services/api.svelte';
	import { QuestionFormatEnum, TestInstanceStateEnum, type TestInstanceDTO } from '$lib/api_types';
	import { page } from '$app/state';
	import { Label } from '$lib/components/ui/label';
	import DateRangeField from '$lib/components/ui/date-range-field/date-range-field.svelte';
	import { enumToOptions } from '$lib/utils';
	import TiptapRenderer from '$lib/components/tiptap-renderer/tiptap-renderer.svelte';
	import Tiptap from '$lib/components/tiptap/Tiptap.svelte';
	import { m } from '$lib/paraglide/messages';
	import { parseAbsoluteToLocal } from '@internationalized/date';

	let {
		instanceData
	}: {
		instanceData: TestInstanceDTO;
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
					bonusPointsReason: instanceData.bonusPointsReason,
				}
			}
		)
			.then((res) => {})
			.catch(() => {});
	}

	const getVariantLabel = (n: number) => {
		let label = '';
		while (n >= 0) {
			label = String.fromCharCode('A'.charCodeAt(0) + (n % 26)) + label;
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
		let answerCount = 0;

		for (const question of instanceData.questions ?? []) {
			if (question.questionFormat == QuestionFormatEnum.ABCD) {
				if (question.answers.length > answerCount) {
					answerCount = question.answers.length;
				}
			}
		}
		return answerCount;
	});
</script>

<div class="flex flex-col gap-4 py-4">
	<Form.TextInput
		title={m.testinstance_participants()}
		id="participant"
		name="participant"
		value="{instanceData.participant.username} ({instanceData.participant.firstName} {instanceData
			.participant.familyName})"
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
		<Form.TextInput
			title={m.testinstance_points()}
			id="points"
			name="points"
			type="number"
			value={instanceData.points + instanceData.bonusPoints}
			disabled
			error=""
		></Form.TextInput>
		<Form.TextInput
			title={m.testinstance_bonuspoints()}
			id="bonusPoints"
			name="bonusPoints"
			type="number"
			bind:value={instanceData.bonusPoints}
			error=""
		></Form.TextInput>
		<!--TODO error-->
		<Form.TextArea
			title={m.testinstance_bonuspoints_reason()}
			id="bonusPointsReason"
			name="bonusPointsReason"
			class="sm:col-span-2"
			bind:value={instanceData.bonusPointsReason}
			error=""
		></Form.TextArea>
		<!--TODO error-->
	</div>
	<Form.Checkbox
		title={m.correctanswers_show()}
		name="showcorrect"
		id="showcorrect"
		bind:value={showCorrect}
		error=""
	></Form.Checkbox>
</div>
<div class="bg-green-500 data-[state=checked]:bg-green-500"></div>
<Form.Root onsubmit={handleSubmitCutomValidation}>
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Cell>{m.question()}</Table.Cell>
				<Table.Cell>{m.question_title()}</Table.Cell>
				<Table.Cell>{m.testinstance_tutor_answer_percentage()}</Table.Cell>
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
					<Table.Cell>
						{question.order}
					</Table.Cell>
					<Table.Cell>
						{question.title}
					</Table.Cell>
					<Table.Cell>
						{#if question.questionFormat == QuestionFormatEnum.OPEN}
							<Form.TextInput
								name="q{question.id}"
								id="q{question.id}"
								bind:value={question.textAnswerPercentage}
								onchange={() => {
									question.textAnswerReviewed = true;
								}}
								error=""
							></Form.TextInput><!-- TODO error-->
						{/if}
					</Table.Cell>
					{#each question.answers as answer}
						<Table.Cell>
							<Form.Checkbox
								innerClass={showCorrect
									? answer.selected == answer.correct
										? 'bg-green-500 data-[state=checked]:bg-green-500 dark:bg-green-500 dark:data-[state=checked]:bg-green-500'
										: 'bg-red-500 data-[state=checked]:bg-red-500 dark:bg-red-500 dark:data-[state=checked]:bg-red-500'
									: ''}
								name="q{question.id}-a{answer.id}"
								id="q{question.id}-a{answer.id}"
								bind:value={answer.selected}
								error=""
							></Form.Checkbox><!-- TODO error-->
						</Table.Cell>
					{/each}
				</Table.Row>
			{:else}
				<Table.Row>
					<Table.Cell>{m.no_questions()}</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
	<div>
		{#each instanceData.questions ?? [] as question}
			<div class="flex flex-col gap-4 p-4 border">
				<div>
					<h2 class="text-xl">{m.question()} {question.order}</h2>
					<TiptapRenderer jsonContent={question.content}></TiptapRenderer>
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
													? answer.selected == answer.correct
														? 'bg-green-500 data-[state=checked]:bg-green-500 dark:bg-green-500 dark:data-[state=checked]:bg-green-500'
														: 'bg-red-500 data-[state=checked]:bg-red-500 dark:bg-red-500 dark:data-[state=checked]:bg-red-500'
													: ''}
												name="q{question.id}-a{answer.id}"
												id="q{question.id}-a{answer.id}"
												bind:value={answer.selected}
												error=""
											></Form.Checkbox><!-- TODO error-->
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
							onchange={() => {
								question.textAnswerReviewed = true;
							}}
							error=""
						></Form.TextInput><!-- TODO error-->
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
	<Form.Button text={m.save()} textSubmiting={m.save_progress()}></Form.Button>
</Form.Root>
