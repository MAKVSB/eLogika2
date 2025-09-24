<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import AppHeader from '$lib/components/app-header.svelte';
	import TiptapRenderer from '$lib/components/tiptap-renderer/tiptap-renderer.svelte';
	import Tiptap from '$lib/components/tiptap/Tiptap.svelte';
	import * as Table from '$lib/components/ui/table/index.js';
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import {
		QuestionFormatEnum,
		TestInstanceStateEnum,
		type TestInstanceDTO,
		type TestInstanceGetResponse,
	} from '$lib/api_types';
	import Button from '$lib/components/ui/button/button.svelte';
	import { onMount } from 'svelte';
	import { API } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import { m } from '$lib/paraglide/messages';

	let instanceData: TestInstanceDTO | undefined = $state();
	let isLoading = $derived(instanceData === undefined);

	onMount(async () => {
		loadTest();
	});

	const loadTest = () => {
		API.request<null, TestInstanceGetResponse>(`/api/v2/tests/${page.params.instanceId}`)
			.then((res) => {
				instanceData = res.instanceData;
			})
			.catch(() => {});
	};

	let questionIndex = $state(0);
	let selectedQuestion = $derived.by(() => {
		if (instanceData) {
			if (instanceData.questions) {
				if (questionIndex <= instanceData.questions.length) {
					return instanceData?.questions[questionIndex];
				}
			}
		}
		return undefined;
	});

	let canSwitchQuestion = (direction: 'next' | 'prev', currentIndex: number): boolean => {
		if (instanceData?.questions) {
			if (direction == 'prev') {
				return currentIndex != 0;
			} else {
				return currentIndex != instanceData.questions.length - 1;
			}
		}
		return false;
	};

	let switchQuestion = (newQuestionIndex: number) => {
		if (instanceData && instanceData.questions) {
			selectedQuestion;
			questionIndex = newQuestionIndex % instanceData.questions.length;
		}
	};
</script>

<div class="[--header-height:calc(--spacing(14))]">
	<Sidebar.Provider class="flex flex-col">
		<AppHeader />
		{#if !isLoading && instanceData}
			{#if instanceData.state != TestInstanceStateEnum.FINISHED}
				<Sidebar.Inset>
					<div class="flex flex-col items-center justify-center h-full gap-4 pt-40 text-center">
						<h1 class="text-2xl">Test is not finished, or its validity has expired</h1>
					</div>
				</Sidebar.Inset>
			{:else}
					<Sidebar.Inset>
						<div class="flex flex-col flex-1 w-full mx-auto contain-inline-size">
							<div class="flex flex-col h-full gap-8 m-8">
								{#if selectedQuestion}
									<div>
										<h1 class="font-bold">{m.testwriter_question()}:</h1>
										<div class="p-4 border">
											<TiptapRenderer jsonContent={selectedQuestion.content!}></TiptapRenderer>
										</div>
									</div>
									<div>
										<h1 class="font-bold">{m.testwriter_answer()}:</h1>
										{#if selectedQuestion.questionFormat == QuestionFormatEnum.ABCD}
											<Table.Root>
												<Table.Body>
													{#each selectedQuestion.answers as answer}
														<Table.Row>
															<Table.Cell>
																<div class="p-4">
																	<Checkbox bind:checked={answer.selected}></Checkbox>
																</div>
															</Table.Cell>

															<Table.Cell class="w-full">
																<TiptapRenderer jsonContent={answer.content}></TiptapRenderer>
															</Table.Cell>
														</Table.Row>
													{/each}
												</Table.Body>
											</Table.Root>
										{:else if selectedQuestion.questionFormat == QuestionFormatEnum.OPEN}
											{#key selectedQuestion.id}
												<Tiptap bind:value={selectedQuestion.textAnswer}></Tiptap>
											{/key}
										{/if}
									</div>
								{:else}
									Failed to load instance data. Please refresh the page
								{/if}
								<div class="flex-grow"></div>
								<div class="flex justify-evenly">
									<Button
										disabled={!canSwitchQuestion('prev', questionIndex)}
										onclick={() => {
											switchQuestion(questionIndex - 1);
										}}
									>
										{m.testwriter_previous()}
									</Button>
									<Button
										disabled={!canSwitchQuestion('next', questionIndex)}
										onclick={() => {
											switchQuestion(questionIndex + 1);
										}}
									>
										{m.testwriter_next()}
									</Button>
								</div>
							</div>
						</div>
					</Sidebar.Inset>
			{/if}
		{/if}
	</Sidebar.Provider>
</div>
