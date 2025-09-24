<script lang="ts">
	import { page } from '$app/state';
	import {
		QuestionFormatEnum,
		QuestionTypeEnum,
		type QuestionListItemDTO,
		type QuestionListRequest,
		type QuestionListResponse
	} from '$lib/api_types';
	import * as Dialog from '$lib/components/ui/dialog/';
	import { API, encodeJsonToBase64Url } from '$lib/services/api.svelte';
	import type { ColumnFiltersState, PaginationState, SortingState } from '@tanstack/table-core';
	import { onMount } from 'svelte';
	import * as Form from '$lib/components/ui/form';
	import { enumToOptions } from '$lib/utils';
	import { m } from '$lib/paraglide/messages';

	let {
		isOpen = $bindable(false),
		selectedQuestions = $bindable()
	}: {
		isOpen: boolean;
		selectedQuestions: any[];
	} = $props();

	let questionTypeFilter = $state(null);
	let questionFormatFilter = $state(null);
	let chapterFilter = $state(null);
	let categoryFilter = $state(null);

	let updateCounter = $state(0);

	type RestRequest = {
		pagination?: PaginationState;
		sorting?: SortingState;
		columnFilters?: ColumnFiltersState;
	};

	let questions: QuestionListItemDTO[] = $state([]);

	onMount(refetchData);

	function refetchData() {
		const queryParams: RestRequest = {
			pagination: {
				pageIndex: 0,
				pageSize: 5000
			},
			columnFilters: []
		};

		if (chapterFilter) {
			queryParams.columnFilters?.push({
				id: 'ChapterID',
				value: chapterFilter
			});
		}

		if (categoryFilter) {
			queryParams.columnFilters?.push({
				id: 'CategoryID',
				value: categoryFilter
			});
		}

		if (questionTypeFilter) {
			queryParams.columnFilters?.push({
				id: 'QuestionType',
				value: questionTypeFilter
			});
		}

		if (questionFormatFilter) {
			queryParams.columnFilters?.push({
				id: 'QuestionFormat',
				value: questionFormatFilter
			});
		}

		API.request<QuestionListRequest, QuestionListResponse>(
			`/api/v2/courses/${page.params.courseId}/questions`,
			{
				searchParams: {
					search: encodeJsonToBase64Url(queryParams)
				}
			}
		)
			.then((res) => {
				questions = res.items;
			})
			.catch(() => {});
	}
</script>

<Dialog.Root bind:open={isOpen}>
	<Dialog.Content class="max-h-full w-300 overflow-scroll sm:max-h-[90%] sm:max-w-[90%]">
		<Dialog.Header>
			<Dialog.Title>Select questions</Dialog.Title>
		</Dialog.Header>

		<div class="grid gap-4 md:grid-cols-2">
			<Form.SingleSelect
				title="Question type"
				placeholder={m.filter_questiontype()}
				bind:value={questionTypeFilter}
				name="questionType"
				id="questionType"
				options={enumToOptions(QuestionTypeEnum, m.question_type_enum)}
				nullable={true}
				onchange={refetchData}
				error=""
			></Form.SingleSelect>
			<Form.SingleSelect
				title="Question format"
				placeholder={m.filter_questionformat()}
				bind:value={questionFormatFilter}
				name="questionFormat"
				id="questionFormat"
				options={enumToOptions(QuestionFormatEnum, m.question_format_enum)}
				nullable={true}
				onchange={refetchData}
				error=""
			></Form.SingleSelect>
		</div>

		<div class="grid gap-4 md:grid-cols-2">
			{#each questions as question, index}
				<div class="flex items-center gap-4">
					{#key updateCounter}
						<Form.Checkbox
							value={selectedQuestions.includes(question.questionGroupId)}
							name="name"
							id="qid-{question.id}"
							error=""
							onclick={() => {
								if (selectedQuestions.includes(question.questionGroupId)) {
									selectedQuestions = selectedQuestions.filter((item) => item != question.questionGroupId);
									updateCounter++;
								} else {
									selectedQuestions.push(question.questionGroupId);
									updateCounter++;
								}
							}}
						></Form.Checkbox>
					{/key}
					<p class="col-span-8">[{question.questionType}] {question.title}</p>
				</div>
			{/each}
		</div>
	</Dialog.Content>
</Dialog.Root>
