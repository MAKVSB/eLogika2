<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Form from '$lib/components/ui/form';
	import * as Table from '$lib/components/ui/table';
	import { TipTapDefaultContent } from '$lib/constants';
	import { type ErrorObject } from '$lib/components/ui/form/types';
	import { Label } from '$lib/components/ui/label';
	import { m } from '$lib/paraglide/messages';
	import { intToLabel } from '$lib/utils';

	let {
		form = $bindable()
	}: {
		form: any;
	} = $props();

	let changeCounter = $state(0);
	let answerCount = $derived(
		form.fields.answers.length
	);

	function deleteAnswer(index: number) {
		form.fields.answers.splice(index, 1)
		changeCounter += 1;
	}
	
	function addAnswer() {
		form.fields.answers.push({
			id: 0,
			version: 1,
			content: TipTapDefaultContent,
			explanation: TipTapDefaultContent,
			timeToSolve: 30,
			correct: false,
		});
	}

	function setAnswerCount(requestedCount: number) {
		let nonDeletedAnswers = 0;

		for (let answer of form.fields.answers) {
			if (!answer.deleted) {
				nonDeletedAnswers += 1;
			}
			if (nonDeletedAnswers > requestedCount) {
				answer.deleted = true;
			}
		}

		while (nonDeletedAnswers < requestedCount) {
			nonDeletedAnswers += 1;
			addAnswer();
		}

		changeCounter += 1;
	}
</script>

<div class="flex flex-col gap-4 p-2 my-4">
	<div class="flex items-center gap-4 pb-4">
		<h1 class="text-xl text-nowrap">{m.question_answer_count()}</h1>
		<Input type="number" bind:value={answerCount}></Input>
		<Button
			variant="outline"
			onclick={() => {
				setAnswerCount(answerCount);
			}}>{m.question_answer_count_set()}</Button
		>
	</div>

	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head class="w-[10px]">{m.question_answer_order()}</Table.Head>
				<Table.Head class="w-[100px]">{m.question_answer_correctness()}</Table.Head>
				<Table.Head>{m.question_answer_answer()}</Table.Head>
				<Table.Head class="w-[100px]">{m.question_answer_advanced()}</Table.Head>
				<Table.Head class="w-[10px]">{m.question_answer_order()}</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each form.fields.answers as answer, index}
				{#if !answer.deleted}
					{@const answersErrors =
						'answers' in form.errors ? (form.errors.answers as ErrorObject) : {}}
					{@const error =
						String(index) in answersErrors ? (answersErrors[String(index)] as ErrorObject) : {}}
					<Table.Row>
						<Table.Cell class="font-bold">{intToLabel(index)}</Table.Cell>
						<Table.Cell>
							<Form.Checkbox
								name="correct"
								id="correct"
								class="w-full h-8"
								innerClass="border-red-500 bg-red-500 data-[state=checked]:border-green-500 data-[state=checked]:bg-green-500 dark:border-red-500 dark:bg-red-500 dark:data-[state=checked]:border-green-500 dark:data-[state=checked]:bg-green-500"
								bind:value={answer.correct}
								error={error.correct}
							></Form.Checkbox>
						</Table.Cell>
						<Table.Cell class="flex flex-col gap-2">
							{#key changeCounter}
								<Form.Tiptap
									title="Answer"
									name="content"
									id="content"
									bind:value={answer.content}
									error={error.content}
									enableFileUpload
									enableFileLink
								></Form.Tiptap>
								<Form.Tiptap
									title="Explanation"
									name="explanation"
									id="explanation"
									bind:value={answer.explanation}
									error={error.explanation}
									enableFileUpload
									enableFileLink
								></Form.Tiptap>
							{/key}
						</Table.Cell>
						<Table.Cell class="h-full">
							<div class="table h-full">
								<div class="table-cell h-full">
									<div class="flex flex-col items-center justify-between h-full gap-4 lg:flex-row">
										<Form.TextInput
											title="Time to solve (s)"
											name="timeToSolve"
											id="timeToSolve"
											type="text"
											bind:value={answer.timeToSolve}
											error={error.timeToSolve}
										></Form.TextInput>
										<div class="flex flex-col gap-2">
											<Label>&nbsp</Label>
											<Button
												variant="destructive"
												onclick={() => {
													deleteAnswer(index)
												}}
											>
												{m.delete()}
											</Button>
										</div>
									</div>
								</div>
							</div>
						</Table.Cell>
						<Table.Cell class="font-bold">{intToLabel(index)}</Table.Cell>
					</Table.Row>
					{/if}
			{/each}
		</Table.Body>
	</Table.Root>
	<div class="flex justify-center gap-4">
		<Button
			variant="outline"
			onclick={() => {
				addAnswer();
				changeCounter += 1;
			}}>{m.answer_add()}</Button
		>
	</div>
</div>
