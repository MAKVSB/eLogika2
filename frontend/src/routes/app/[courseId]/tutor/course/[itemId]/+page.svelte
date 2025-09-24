<script lang="ts">
	import { API } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import {
		StudyFormEnum,
		type CourseItemDTO,
		CourseItemTypeEnum,
		QuestionTypeEnum,
		type CourseItemInsertRequest,
		type CourseItemInsertResponse,
		type CourseItemGetByIdResponse,
		type CourseItemUpdateRequest,
		type CourseItemUpdateResponse,
		CourseUserRoleEnum,
		EvaluateByAttemptEnum
	} from '$lib/api_types';
	import * as Form from '$lib/components/ui/form';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import { TipTapDefaultContent } from '$lib/constants';
	import type { ErrorObject } from '$lib/components/ui/form/types';
	import CourseItems from '../CourseItems/CourseItems.svelte';
	import { CourseItemInsertRequestSchema } from '$lib/schemas';
	import TestItem from './TestItem.svelte';
	import ActivityItem from './ActivityItem.svelte';
	import GroupItem from './GroupItem.svelte';
	import Terms from './Terms/terms.svelte';
	import GlobalState from '$lib/shared.svelte';
	import { deepMerge, enumToOptions } from '$lib/utils';
	import { m } from '$lib/paraglide/messages';

	let courseId = $derived<string>(page.params.courseId);
	let itemId = $derived<string>(page.params.itemId);

	function parseStudyFormEnum(str: string | null): StudyFormEnum {
		if (str && str == StudyFormEnum.COMBINED) {
			return StudyFormEnum.COMBINED;
		} else if (str && str == StudyFormEnum.FULLTIME) {
			return StudyFormEnum.FULLTIME;
		} else {
			return StudyFormEnum.FULLTIME;
		}
	}

	function parseItemTypeEnum(str: string | null): CourseItemTypeEnum {
		if (str && str == CourseItemTypeEnum.TEST) {
			return CourseItemTypeEnum.TEST;
		} else if (str && str == CourseItemTypeEnum.ACTIVITY) {
			return CourseItemTypeEnum.ACTIVITY;
		} else if (str && str == CourseItemTypeEnum.GROUP) {
			return CourseItemTypeEnum.GROUP;
		} else {
			return CourseItemTypeEnum.TEST;
		}
	}

	function parseParentData(): number | undefined {
		const parentId = page.url.searchParams.get('parentId');
		if (!parentId) {
			return undefined;
		}
		const parentIdNum = Number(parentId);
		if (isNaN(parentIdNum)) {
			return undefined;
		}
		return parentIdNum;
	}

	let { data } = $props();
	$effect(() => {
		if (data.courseItem) {
			data.courseItem
				.then((data) => setResult(data))
				.catch((err) => {
					console.error(err);
					toast.error('Failed to load question');
				});
		} else {
			form.fields = defaultFormData();
		}
	});

	function defaultFormData(): CourseItemDTO {
		const type = parseItemTypeEnum(page.url.searchParams.get('type'));
		const studyForm = parseStudyFormEnum(page.url.searchParams.get('studyForm'));
		const formData: CourseItemDTO & CourseItemInsertRequest = {
			id: 0,
			version: 0,
			name: '',
			type: type,
			pointsMin: 0,
			pointsMax: 0,
			mandatory: true,
			studyForm: studyForm,
			maxAttempts: 1,
			allowNegative: false,
			editable: true,
			managedBy:
				GlobalState.activeRole == CourseUserRoleEnum.TUTOR
					? CourseUserRoleEnum.TUTOR
					: CourseUserRoleEnum.GARANT,
			evaluateByAttempt: EvaluateByAttemptEnum.LAST
		};

		const parentId = parseParentData();
		if (parentId) {
			formData.parentId = parentId;
		}

		switch (type) {
			case CourseItemTypeEnum.ACTIVITY:
				formData.activityDetail = {
					id: 0,
					description: TipTapDefaultContent,
					expectedResult: TipTapDefaultContent
				};
				break;
			case CourseItemTypeEnum.TEST:
				formData.testDetail = {
					id: 0,
					testType: QuestionTypeEnum.EXAM,
					testTemplateId: 0,
					timeLimit: 60,
					showResults: false,
					showTest: false,
					allowOffline: false,
					isPaper: false,
					ipRanges: ''
				};
				break;
			case CourseItemTypeEnum.GROUP:
				formData.groupDetail = {
					id: 0,
					choice: false,
					chooseMin: 1,
					chooseMax: 1
				};
				break;
			default:
				break;
		}
		return formData;
	}
	let form = $state(Form.createForm(CourseItemInsertRequestSchema, defaultFormData()));

	function setResult(
		res: CourseItemGetByIdResponse | CourseItemInsertResponse | CourseItemUpdateResponse
	) {
		form.fields = res.data;
		console.log("Transfering 24")
		goto(String(res.data.id), {
			replaceState: true
		});
	}

	async function handleSubmit(): Promise<any> {
		let request;
		if (data.creating) {
			request = API.request<CourseItemInsertRequest, CourseItemInsertResponse>(
				`/api/v2/courses/${page.params.courseId}/items`,
				{
					method: 'POST',
					body: {
						...form.fields,
						parentId: form.fields.parentId ?? undefined
					}
				}
			);
		} else {
			request = API.request<CourseItemUpdateRequest, CourseItemUpdateResponse>(
				`/api/v2/courses/${page.params.courseId}/items/${page.params.itemId}`,
				{
					method: 'PUT',
					body: {
						...form.fields
					}
				}
			);
		}
		return request
			.then((res) => setResult(res));
	}
</script>

<div class="flex flex-col gap-8 m-8">
	{#await data.courseItem}
		<Pageloader></Pageloader>
	{:then staticResourceData}
		<div class="flex flex-row justify-between">
			<h1 class="mb-8 text-2xl">
				Course item management:
				<b>
					{staticResourceData?.data?.name ?? 'New course item'}
				</b>
			</h1>
		</div>
		<Form.Root
			bind:form
			onsubmit={handleSubmit}
			isCreating={data.creating}
			hideDefaultbutton={!form.fields.editable}
		>
			<div class="flex flex-col gap-4 p-2">
				<Form.TextInput
					title={m.courseitem_name()}
					name="name"
					id="name"
					type="text"
					bind:value={form.fields.name}
					error={form.errors.name ?? ''}
					disabled={!form.fields.editable}
				></Form.TextInput>
				<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
					<Form.TextInput
						title={m.courseitem_min_points()}
						name="pointsMin"
						id="pointsMin"
						type="number"
						bind:value={form.fields.pointsMin}
						error={form.errors.pointsMin ?? ''}
						disabled={!form.fields.editable}
					></Form.TextInput>
					<Form.TextInput
						title={m.courseitem_max_points()}
						name="pointsMax"
						id="pointsMax"
						type="number"
						bind:value={form.fields.pointsMax}
						error={form.errors.pointsMax ?? ''}
						disabled={!form.fields.editable}
					></Form.TextInput>
					<Form.Checkbox
						title={m.courseitem_allownegative()}
						name="allowNegative"
						id="allowNegative"
						bind:value={form.fields.allowNegative}
						error={form.errors.allowNegative}
						disabled={!form.fields.editable}
					></Form.Checkbox>
					<Form.SingleSelect
						title={m.courseitem_evaluatebyattempt()}
						name="evaluateByAttempt"
						id="evaluateByAttempt"
						bind:value={form.fields.evaluateByAttempt}
						error={form.errors.evaluateByAttempt}
						options={enumToOptions(EvaluateByAttemptEnum)}
						disabled={!form.fields.editable}
					></Form.SingleSelect>
				</div>
				<div class="grid grid-cols-3 gap-4">
					<Form.TextInput
						title={m.courseitem_maxattempts()}
						name="maxAttempts"
						id="maxAttempts"
						type="number"
						class="col-span-2"
						bind:value={form.fields.maxAttempts}
						error={form.errors.maxAttempts ?? ''}
						disabled={!form.fields.editable}
					></Form.TextInput>
					<Form.Checkbox
						title={m.course_item_mandatory()}
						name="mandatory"
						id="mandatory"
						bind:value={form.fields.mandatory}
						error={form.errors.mandatory}
						disabled={!form.fields.editable}
					></Form.Checkbox>
				</div>
				{#if form.fields.testDetail && form.fields.type == CourseItemTypeEnum.TEST}
					<TestItem
						bind:fields={form.fields.testDetail}
						errors={(form.errors.testDetail ?? {}) as ErrorObject}
						disabled={!form.fields.editable}
					></TestItem>
				{:else if form.fields.activityDetail && form.fields.type == CourseItemTypeEnum.ACTIVITY}
					<ActivityItem
						bind:fields={form.fields.activityDetail}
						errors={(form.errors.activityDetail ?? {}) as ErrorObject}
						disabled={!form.fields.editable}
					></ActivityItem>
				{:else if form.fields.groupDetail && form.fields.type == CourseItemTypeEnum.GROUP}
					<GroupItem
						bind:fields={form.fields.groupDetail}
						errors={(form.errors.groupDetail ?? {}) as ErrorObject}
						disabled={!form.fields.editable}
					></GroupItem>
				{/if}
			</div>
		</Form.Root>
		{#if staticResourceData}
			{#if form.fields.type == CourseItemTypeEnum.GROUP}
				<CourseItems
					canGroup={!!form.fields.parentId}
					{courseId}
					parentId={staticResourceData?.data.id}
					mode={form.fields.studyForm}
				></CourseItems>
			{/if}
			{#if staticResourceData.data.editable}
				<Terms {courseId} {itemId}></Terms>
			{/if}
		{/if}
	{/await}
</div>
