<script lang="ts">
	import { API, ApiError } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import {
		type CourseInsertRequest,
		type CourseInsertResponse,
		SemesterEnum,
		type CourseDTO,
		type CourseUpdateRequest,
		type CourseUpdateResponse,
		type CourseGetByIdResponse,
		CourseUserRoleEnum
	} from '$lib/api_types';
	import * as Form from '$lib/components/ui/form';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import { TipTapDefaultContent } from '$lib/constants';
	import { enumToOptions } from '$lib/utils';
	import { m } from '$lib/paraglide/messages';
	import { CourseInsertRequestSchema } from '$lib/schemas';
	import type { ErrorObject } from '$lib/components/ui/form';
	import GlobalState from '$lib/shared.svelte';

	let { data } = $props();

	$effect(() => {
		if (data.course) {
			data.course.then((data) => setResult(data)).catch(() => {});
		} else {
			form.fields = defaultFormData;
		}
	});

	const defaultFormData: CourseDTO = {
		id: 0,
		version: 0,
		name: '',
		content: TipTapDefaultContent,
		shortname: '',
		public: false,
		year: new Date().getFullYear(),
		semester: SemesterEnum.SUMMER,
		pointsMax: 100,
		pointsMin: 51,
		importOptions: {
			date: '',
			code: ''
		}
	};
	let form = $state(Form.createForm(CourseInsertRequestSchema, defaultFormData));

	function setResult(res: CourseGetByIdResponse | CourseInsertResponse | CourseUpdateResponse) {
		form.fields = res.data;
		console.log('Transfering 23');
		goto(page.url, {
			replaceState: true
		});
	}

	async function handleSubmit(): Promise<any> {
		let request;
		if (data.creating) {
			request = API.request<CourseInsertRequest, CourseInsertResponse>(`/api/v2/courses`, {
				method: 'POST',
				body: form.fields
			});
		} else {
			request = API.request<CourseUpdateRequest, CourseUpdateResponse>(
				`/api/v2/courses/${page.params.courseId}`,
				{
					method: 'PUT',
					body: form.fields
				}
			);
		}

		return request.then((res) => {
			setResult(res)
			if (data.creating) {
				toast.success("Created succesfully")
			} else {
				toast.success("Saved succesfully")
			}
		});
	}

	let disabled = $state(
		GlobalState.activeRole != CourseUserRoleEnum.ADMIN &&
			GlobalState.activeRole != CourseUserRoleEnum.GARANT
	);
</script>

<div class="flex flex-col gap-8 m-8">
	{#await data.course}
		<Pageloader></Pageloader>
	{:then staticResourceData}
		<div class="flex flex-row justify-between">
			<h1 class="text-2xl">
				Course management:
				<b>
					{staticResourceData?.data?.name ?? 'New course'}
				</b>
			</h1>
		</div>
		<Form.Root bind:form onsubmit={handleSubmit} isCreating={data.creating}>
			<div class="flex flex-col gap-4 p-2">
				<div class="grid grid-cols-1 gap-4 md:grid-cols-3">
					<Form.TextInput
						title={m.course_name()}
						name="name"
						id="name"
						type="text"
						class="col-span-2"
						bind:value={form.fields.name}
						{disabled}
						error={form.errors.name ?? ''}
					></Form.TextInput>
					<Form.TextInput
						title={m.course_shortname()}
						name="shortname"
						id="shortname"
						type="text"
						bind:value={form.fields.shortname}
						{disabled}
						error={form.errors.shortname ?? ''}
					></Form.TextInput>
				</div>
				<div class="grid grid-cols-1 gap-4 md:grid-cols-3">
					<Form.YearInput
						title={m.academic_year()}
						name="year"
						id="year"
						bind:value={form.fields.year}
						{disabled}
						error={form.errors.year ?? ''}
					></Form.YearInput>
					<Form.SingleSelect
						title={m.semester()}
						name="semester"
						id="semester"
						bind:value={form.fields.semester}
						{disabled}
						options={enumToOptions(SemesterEnum, m.semester_enum)}
						error={form.errors.semester ?? ''}
					></Form.SingleSelect>
					<Form.Checkbox
						title={m.public()}
						name="public"
						id="public"
						bind:value={form.fields.public}
						{disabled}
						error={form.errors.public ?? ''}
					></Form.Checkbox>
				</div>
				<div class="grid grid-cols-1 gap-4 md:grid-cols-3">
					<Form.TextInput
						title={m.course_points_min()}
						name="pointsMin"
						id="pointsMin"
						type="number"
						bind:value={form.fields.pointsMin}
						{disabled}
						error={form.errors.pointsMin ?? ''}
					></Form.TextInput>
					<Form.TextInput
						title={m.course_points_max()}
						name="pointsMax"
						id="pointsMax"
						type="number"
						bind:value={form.fields.pointsMax}
						{disabled}
						error={form.errors.pointsMax ?? ''}
					></Form.TextInput>
				</div>
				{#key form.fields}
					<Form.Tiptap
						title={m.course_content()}
						name="content"
						id="content"
						bind:value={form.fields.content}
						{disabled}
						error={form.errors.content}
						enableFileUpload
						enableFileLink
					></Form.Tiptap>
				{/key}
				<div class="flex flex-col gap-4">
					<h3>{m.import_options()}:</h3>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
						<Form.TextInput
							title={m.course_import_code()}
							name="import_code"
							id="impot_code"
							bind:value={form.fields.importOptions.code}
							{disabled}
							error={((form.errors.importOptions as ErrorObject) ?? {}).code ?? ''}
						></Form.TextInput>
						<Form.TextInput
							title={m.course_import_semesterdate()}
							name="import_date"
							id="import_date"
							bind:value={form.fields.importOptions.date}
							{disabled}
							error={((form.errors.importOptions as ErrorObject) ?? {}).date ?? ''}
						></Form.TextInput>
					</div>
				</div>
			</div>
		</Form.Root>
	{/await}
</div>
