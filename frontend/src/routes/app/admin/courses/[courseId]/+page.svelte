<script lang="ts">
	import { API } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import {
		type CourseInsertRequest,
		type CourseInsertResponse,
		SemesterEnum,
		type CourseDTO,
		type CourseUpdateResponse,
		type CourseGetByIdResponse
	} from '$lib/api_types';
	import * as Form from '$lib/components/ui/form';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import { TipTapDefaultContent } from '$lib/constants';
	import { enumToOptions } from '$lib/utils';
	import { m } from '$lib/paraglide/messages';
	import { CourseInsertRequestSchema } from '$lib/schemas';
	import type { ErrorObject } from '$lib/components/ui/form';
	import { base } from '$app/paths';

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

	async function handleSubmit(): Promise<any> {
		let request = API.request<CourseInsertRequest, CourseInsertResponse>(`/api/v2/courses`, {
				method: 'POST',
				body: form.fields
			});

		return request.then((res) => {
			toast.success("Created succesfully")
			goto(base + `/app/${res.data.id}/tutor/course`);
		});
	}
</script>

<div class="flex flex-col gap-8 m-8">
		<div class="flex flex-row justify-between">
			<h1 class="text-2xl">
				Course management:
				<b>
					New course
				</b>
			</h1>
		</div>
		<Form.Root bind:form onsubmit={handleSubmit} isCreating>
			<div class="flex flex-col gap-4 p-2">
				<div class="grid grid-cols-1 gap-4 md:grid-cols-3">
					<Form.TextInput
						title={m.course_name()}
						name="name"
						id="name"
						type="text"
						class="col-span-2"
						bind:value={form.fields.name}
						error={form.errors.name ?? ''}
					></Form.TextInput>
					<Form.TextInput
						title={m.course_shortname()}
						name="shortname"
						id="shortname"
						type="text"
						bind:value={form.fields.shortname}
						error={form.errors.shortname ?? ''}
					></Form.TextInput>
				</div>
				<div class="grid grid-cols-1 gap-4 md:grid-cols-3">
					<Form.YearInput
						title={m.academic_year()}
						name="year"
						id="year"
						bind:value={form.fields.year}
						error={form.errors.year ?? ''}
					></Form.YearInput>
					<Form.SingleSelect
						title={m.semester()}
						name="semester"
						id="semester"
						bind:value={form.fields.semester}
						options={enumToOptions(SemesterEnum, m.semester_enum)}
						error={form.errors.semester ?? ''}
					></Form.SingleSelect>
					<Form.Checkbox
						title={m.public()}
						name="public"
						id="public"
						bind:value={form.fields.public}
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
						error={form.errors.pointsMin ?? ''}
					></Form.TextInput>
					<Form.TextInput
						title={m.course_points_max()}
						name="pointsMax"
						id="pointsMax"
						type="number"
						bind:value={form.fields.pointsMax}
						error={form.errors.pointsMax ?? ''}
					></Form.TextInput>
				</div>
				{#key form.fields}
					<Form.Tiptap
						title={m.course_content()}
						name="content"
						id="content"
						bind:value={form.fields.content}
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
							error={((form.errors.importOptions as ErrorObject) ?? {}).code ?? ''}
						></Form.TextInput>
						<Form.TextInput
							title={m.course_import_semesterdate()}
							name="import_date"
							id="import_date"
							bind:value={form.fields.importOptions.date}
							error={((form.errors.importOptions as ErrorObject) ?? {}).date ?? ''}
						></Form.TextInput>
					</div>
				</div>
			</div>
		</Form.Root>
</div>
