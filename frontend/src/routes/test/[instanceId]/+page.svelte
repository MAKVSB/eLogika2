<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import AppHeader from '$lib/components/app-header.svelte';
	import TestSidebar from '$lib/components/test-sidebar.svelte';
	import TiptapRenderer from '$lib/components/tiptap-renderer/tiptap-renderer.svelte';
	import Tiptap from '$lib/components/tiptap/Tiptap.svelte';
	import * as Table from '$lib/components/ui/table/index.js';
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import {
		QuestionFormatEnum,
		TestInstanceEventTypeEnum,
		TestInstanceStateEnum,
		type TestInstanceDTO,
		type TestInstanceGetResponse,
		type TestInstanceQuestion,
		type TestInstanceQuestionAnswerDTO,
		type TestInstanceQuestionDTO,
		type TestInstanceStartResponse
	} from '$lib/api_types';
	import Button from '$lib/components/ui/button/button.svelte';
	import { onDestroy, onMount } from 'svelte';
	import { API, ApiError } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import { m } from '$lib/paraglide/messages';
	import ExamLogger from './ExamLogger';
	import { cn } from '$lib/utils';

	function formatTime(seconds: number): string {
		const h = Math.floor(seconds / 3600);
		const m = Math.floor((seconds % 3600) / 60);
		const s = seconds % 60;
		return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
	}

	let instanceData: TestInstanceDTO | undefined = $state();
	let endTime: number = $state(Date.now());
	let isLoading = $derived(instanceData === undefined);

	onMount(async () => {
		loadTest();
	});

	const loadTest = () => {
		API.request<null, TestInstanceGetResponse>(`/api/v2/tests/${page.params.instanceId}`)
			.then((res) => {
				instanceData = res.instanceData;
				endTime = Math.floor(new Date(res.instanceData.endsAt).getTime() / 1000);
				if (
					instanceData &&
					(instanceData.state == TestInstanceStateEnum.ACTIVE ||
						instanceData.state == TestInstanceStateEnum.READY)
				) {
					if (logger.running() == false) {
						logger.start();
					}
				}
			})
			.catch(() => {});
	};

	const startTest = async () => {
		logger.record(TestInstanceEventTypeEnum.TESTSTART);
		await API.request<null, TestInstanceStartResponse>(
			`/api/v2/tests/${page.params.instanceId}/start`,
			{
				method: 'PUT'
			}
		)
			.then(() => loadTest())
			.catch(() => {});
	};

	const saveAnswers = async (instanceQuestionData: TestInstanceQuestionDTO) => {
		await API.request<TestInstanceQuestion, TestInstanceGetResponse>(
			`/api/v2/tests/${page.params.instanceId}/save`,
			{
				method: 'PUT',
				body: {
					id: instanceQuestionData.id,
					answers: instanceQuestionData.answers.map((a: TestInstanceQuestionAnswerDTO) => {
						return {
							id: a.id,
							selected: a.selected
						};
					}),
					textAnswer: instanceQuestionData.textAnswer
				}
			}
		);
	};

	const finishTest = async () => {
		if (instanceData && instanceData.questions) {
			await API.request<any, TestInstanceGetResponse>(
				`/api/v2/tests/${page.params.instanceId}/finish`,
				{
					method: 'PUT',
					body: {
						id: instanceData.id,
						questions: instanceData?.questions.map((q) => {
							return {
								id: q.id,
								textAnswer: q.textAnswer,
								answers: q.answers.map((a) => {
									return {
										id: a.id,
										selected: a.selected
									};
								})
							};
						})
					}
				}
			)
				.then((res) => {
					logger.stop();
					if (instanceData) {
						instanceData.state = TestInstanceStateEnum.FINISHED;
					}
				})
				.catch(() => {});
		}
	};

	let currentTime = $state(Math.floor(Date.now() / 1000));
	onMount(() => {
		const interval = setInterval(() => {
			currentTime = Math.floor(Date.now() / 1000);
		}, 1000);

		return () => {
			clearInterval(interval);
		};
	});
	let timeLeft = $derived.by(() => {
		const timeDiff = endTime - currentTime;
		if (timeDiff < 0) {
			if (instanceData?.state == TestInstanceStateEnum.ACTIVE) {
				finishTest();
			}
		}
		return formatTime(endTime - currentTime);
	});
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
			selectedQuestion && saveAnswers(selectedQuestion);
			questionIndex = newQuestionIndex % instanceData.questions.length;
			logger.record(TestInstanceEventTypeEnum.QUESTIONSWITCHED, {
				questionId: instanceData.questions[questionIndex].id,
				questionOrder: questionIndex
			});
		}
	};

	//#region Logger
	let logger: ExamLogger;

	onMount(() => {
		logger = new ExamLogger({
			instanceId: page.params.instanceId,
			includeDebugWidget: true,
			endpoint: import.meta.env.VITE_API_URL + `/api/v2/tests/${page.params.instanceId}/telemetry`
		});
		return () => {
			logger.stop();
		};
	});
	//#endregion
</script>

<div class="[--header-height:calc(--spacing(14))]">
	<Sidebar.Provider class="flex flex-col">
		<AppHeader />
		{#if !isLoading && instanceData}
			{#if instanceData.state == TestInstanceStateEnum.READY}
				<Sidebar.Inset>
					<div class="flex flex-col items-center h-full gap-4 pt-40">
						<h1 class="text-2xl">Your test is ready</h1>
						<p>Time limit for this test is: {instanceData.timeLimit} minutes</p>
						<Button variant="default" onclick={startTest}>Start test</Button>
					</div>
				</Sidebar.Inset>
			{:else if instanceData.state == TestInstanceStateEnum.FINISHED}
				<Sidebar.Inset>
					<div class="flex flex-col items-center justify-center h-full gap-4 pt-40 text-center">
						<h1 class="text-2xl">Test has finished</h1>
						<p>Time limit for this test was: {instanceData.timeLimit} minutes</p>
						<p>You can now close this window. Results will be soon available in eLogika</p>
					</div>
				</Sidebar.Inset>
			{:else if instanceData.state == TestInstanceStateEnum.ACTIVE}
				<div class="flex flex-1">
					<TestSidebar bind:questionIndex bind:instanceData />
					<Sidebar.Inset>
						<div class="flex flex-col flex-1 w-full mx-auto contain-inline-size">
							<div class="flex flex-col h-full gap-8 m-8">
								<div class="text-right">
									{m.testwriter_timeleft()}
									<span class="text-red-500">
										{timeLeft}
									</span>
								</div>
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
									<Button
										variant={'destructive'}
										onclick={() => {
											finishTest();
										}}
									>
										{m.testwriter_saveandsend()}
									</Button>
								</div>
							</div>
						</div>
					</Sidebar.Inset>
				</div>
			{/if}
		{/if}
	</Sidebar.Provider>
</div>
