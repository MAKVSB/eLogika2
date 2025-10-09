<script lang="ts">
	import { API } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import * as Form from '$lib/components/ui/form';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import {
		QuestionTypeEnum,
		QuestionFormatEnum,
		type QuestionAdminDTO,
		type QuestionCheckedByDTO,
		type QuestionCheckResponse,
		type QuestionGetByIdResponse,
		type QuestionInsertResponse,
		type QuestionUpdateResponse,
		type QuestionInsertRequest,
		type QuestionUpdateRequest
	} from '$lib/api_types';
	import { Button } from '$lib/components/ui/button';
	import GlobalState from '$lib/shared.svelte';
	import { TipTapDefaultContent } from '$lib/constants';
	import { m } from '$lib/paraglide/messages';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import TabAnswer from './TabAnswer.svelte';
	import TabSteps from './TabSteps.svelte';
	import TabQuestion from './TabQuestion.svelte';
	import { QuestionInsertRequestSchema } from '$lib/schemas';
	import { displayUserName } from '$lib/utils';

	let courseId = $derived<string>(page.params.courseId);
	let { data } = $props();

	let isLoaded = $state(false);

	$effect(() => {
		if (data.question) {
			data.question
				.then((data) => setResult(data))
				.catch((err) => {
					console.error(err);
					toast.error('Failed to load question');
				});
		} else {
			form.fields = defaultFormData;
			isLoaded = true;
		}
	});

	const defaultFormData: QuestionAdminDTO = {
		id: 0,
		version: 0,
		title: '',
		content: TipTapDefaultContent,
		timeToRead: 30,
		timeToProcess: 30,
		questionType: QuestionTypeEnum.EXAM,
		questionFormat: QuestionFormatEnum.ABCD,
		createdBy: GlobalState.loggedUser!,
		active: true,
		chapterId: 0,
		categoryId: 0,
		steps: [],
		answers: [],
		checkedBy: []
	};
	let form = $state(Form.createForm(QuestionInsertRequestSchema, defaultFormData));

	function setResult(
		res: QuestionGetByIdResponse | QuestionInsertResponse | QuestionUpdateResponse
	) {
		form.fields = res.data;
		isLoaded = true;
		console.log('Transfering 30');
		goto(String(res.data.id), {
			replaceState: true
		});
	}

	async function handleSubmit(e: any) {
		let request;
		if (data.creating) {
			request = API.request<QuestionInsertRequest, QuestionInsertResponse>(
				`/api/v2/courses/${courseId}/questions`,
				{
					method: 'POST',
					body: form.fields
				}
			);
		} else {
			if (!e.submitter.dataset.newVersion) {
				if (
					!window.confirm(
						'This question may be already used in test. Saving without new version will override already generated tests.'
					)
				) {
					return;
				}
			}
			request = API.request<QuestionUpdateRequest, QuestionUpdateResponse>(
				`/api/v2/courses/${courseId}/questions/${page.params.id}`,
				{
					method: 'PUT',
					body: {
						...form.fields,
						asNewVersion: e.submitter.dataset.newVersion ? true : false
					}
				}
			);
		}

		await request
			.then((res) => {
				setResult(res);
				toast.success('Saved');
			})
			.catch(() => {});
	}

	async function toggleCheck(check: boolean) {
		if (!data.creating) {
			await API.request<null, QuestionCheckResponse>(
				`api/v2/courses/${courseId}/questions/${page.params.id}/check`,
				{
					method: check ? 'POST' : 'DELETE'
				}
			)
				.then((res) => {
					form.fields.checkedBy = res.checkedBy;
					toast.success('Saved');
				})
				.catch(() => {});
		}
	}

	function printUserName(user: QuestionCheckedByDTO, last: boolean) {
		return  displayUserName(user) + `${last ? '' : ', '}`;
	}
</script>

{#snippet additionalButtons()}
	<Form.Button
		text={data.creating ? m.create() : 'Save as new version'}
		textSubmiting={data.creating ? m.create_progress() : m.save_progress()}
		data-new-version
	></Form.Button>
	{#if !data.creating}
		<Form.Button
			isSubmitting={form.isSubmitting}
			text="Save without creating version"
			textSubmiting={m.save_progress()}
		></Form.Button>
		{#if form.fields.checkedBy.find((usr) => usr.id == GlobalState.loggedUser?.id)}
			<Button variant="destructive" onclick={() => toggleCheck(false)}
				>{m.question_check_action_uncheck()}</Button
			>
		{:else}
			<Button variant="outline" onclick={() => toggleCheck(true)} class="bg-green-500">
				{m.question_check_action_check()}
			</Button>
		{/if}
	{/if}
{/snippet}

<div class="m-8">
	{#await data.question}
		<Pageloader></Pageloader>
	{:then staticResourceData}
		{#if isLoaded}
			<div class="flex flex-row justify-between">
				<h1 class="mb-8 text-2xl">
					Question management:
					<b>
						{staticResourceData?.data?.title ?? 'New question'}
					</b>
				</h1>
			</div>
			<Form.Root
				bind:form
				onsubmit={handleSubmit}
				isCreating={data.creating}
				{additionalButtons}
				hideDefaultbutton
			>
				<Tabs.Root value="question">
					<Tabs.List>
						<Tabs.Trigger value="question">{m.question_tab_main()}</Tabs.Trigger>
						<Tabs.Trigger value="answers">{m.question_tab_answers()}</Tabs.Trigger>
						<Tabs.Trigger value="steps">{m.question_tab_steps()}</Tabs.Trigger>
					</Tabs.List>
					<Tabs.Content value="question">
						<TabQuestion bind:form></TabQuestion>
					</Tabs.Content>
					<Tabs.Content value="answers">
						<TabAnswer bind:form></TabAnswer>
					</Tabs.Content>
					<Tabs.Content value="steps">
						<TabSteps bind:form {courseId}></TabSteps>
					</Tabs.Content>
				</Tabs.Root>
			</Form.Root>

			<div class="flex flex-col gap-8 p-2">
				{#if staticResourceData?.data.id}
					<div class="flex gap-4 p-4 my-4 border rounded-md grow">
						{m.question_check_checkedby()}:
						{#each form.fields.checkedBy as user, index}
							{printUserName(user, staticResourceData?.data.checkedBy.length - 1 == index)}
						{/each}
					</div>
				{/if}
			</div>
		{/if}
	{/await}
</div>
