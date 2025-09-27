<script lang="ts">
	import { onMount } from 'svelte';
	import { API, ApiError } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import {
		StudyFormEnum,
		type ClassDTO,
		ClassTypeEnum,
		WeekDayEnum,
		WeekParityEnum,
		type ClassInsertRequest,
		type ClassUpdateRequest,
		type ClassInsertResponse,
		type ClassUpdateResponse,
		type ClassGetByIdResponse
	} from '$lib/api_types';
	import * as Form from '$lib/components/ui/form';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import { ClassInsertRequestSchema } from '$lib/schemas';
	import { enumToOptions } from '$lib/utils.js';
	import { m } from '$lib/paraglide/messages.js';
	import TutorView from './TutorView/TutorView.svelte';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import UserAddDialog from './UserAddDialog/UserAddDialog.svelte';
	import { buttonVariants } from '$lib/components/ui/button';

	let courseId = $derived<string>(page.params.courseId);
	let itemId = $derived<string>(page.params.itemId);
	let dialogOpen = $state(false);

	let { data } = $props();
	$effect(() => {
		if (data.courseItem) {
			data.courseItem.then((data) => setResult(data)).catch(() => {});
		} else {
			form.fields = defaultFormData;
		}
	});

	const defaultFormData: ClassDTO = {
		id: 0,
		version: 0,
		name: '',
		room: '',
		type: ClassTypeEnum.C,
		studyForm: StudyFormEnum.FULLTIME,
		timeFrom: '07:15',
		timeTo: '08:45',
		day: WeekDayEnum.MONDAY,
		weekParity: WeekParityEnum.BOTH,
		studentLimit: 0,
		tutors: []
	};
	let form = $state(Form.createForm(ClassInsertRequestSchema, defaultFormData));

	function setResult(res: ClassGetByIdResponse | ClassInsertResponse | ClassUpdateResponse) {
		form.fields = res.data;
		console.log("Transfering 22")
		goto(String(res.data.id), {
			replaceState: true
		});
	}

	async function handleSubmit(): Promise<any> {
		let request;
		if (data.creating) {
			request = API.request<ClassInsertRequest, ClassInsertResponse>(
				`/api/v2/courses/${page.params.courseId}/classes`,
				{
					method: 'POST',
					body: form.fields
				}
			);
		} else {
			request = API.request<ClassUpdateRequest, ClassUpdateResponse>(
				`/api/v2/courses/${page.params.courseId}/classes/${page.params.classId}`,
				{
					method: 'PUT',
					body: form.fields
				}
			);
		}
		return request.then((res) => setResult(res));
	}
</script>

<div class="flex flex-col gap-8 m-8">
	{#await data.courseItem}
		<Pageloader></Pageloader>
	{:then staticResourceData}
		<div class="flex flex-row justify-between">
			<h1 class="mb-8 text-2xl">
				Class management:
				<b>
					{staticResourceData?.data?.name ?? 'New class'}
				</b>
			</h1>
		</div>
		<Form.Root bind:form onsubmit={handleSubmit} isCreating={data.creating}>
			<div class="flex flex-col gap-4 p-2">
				<Form.TextInput
					title={m.classes_name()}
					name="name"
					id="name"
					type="text"
					bind:value={form.fields.name}
					error={form.errors.name ?? ''}
				></Form.TextInput>
				<Form.TextInput
					title={m.classes_room()}
					name="room"
					id="room"
					type="text"
					bind:value={form.fields.room}
					error={form.errors.room ?? ''}
				></Form.TextInput>
				<Form.SingleSelect
					title={m.classes_type()}
					name="type"
					id="type"
					options={enumToOptions(ClassTypeEnum, m.class_type_enum)}
					bind:value={form.fields.type}
					error={form.errors.type ?? ''}
				></Form.SingleSelect>
				<Form.SingleSelect
					title={m.classes_studyform()}
					name="studyForm"
					id="studyForm"
					options={enumToOptions(StudyFormEnum, m.study_form_enum)}
					bind:value={form.fields.studyForm}
					error={form.errors.studyForm ?? ''}
				></Form.SingleSelect>
				<Form.SingleSelect
					title={m.classes_day()}
					name="day"
					id="day"
					options={enumToOptions(WeekDayEnum, m.week_day_enum)}
					bind:value={form.fields.day}
					error={form.errors.day ?? ''}
				></Form.SingleSelect>
				<Form.SingleSelect
					title={m.classes_weekparity()}
					name="weekParity"
					id="weekParity"
					options={enumToOptions(WeekParityEnum, m.week_parity_enum)}
					bind:value={form.fields.weekParity}
					error={form.errors.weekParity ?? ''}
				></Form.SingleSelect>
				<Form.TextInput
					title={m.classes_studentslimit()}
					name="studentLimit"
					id="studentLimit"
					type="number"
					bind:value={form.fields.studentLimit}
					error={form.errors.studentLimit ?? ''}
				></Form.TextInput>
				<Form.TextInput
					title={m.classes_timefrom()}
					name="timeFrom"
					id="timeFrom"
					type="text"
					bind:value={form.fields.timeFrom}
					error={form.errors.timeFrom ?? ''}
				></Form.TextInput>
				<Form.TextInput
					title={m.classes_timeto()}
					name="timeTo"
					id="timeTo"
					type="text"
					bind:value={form.fields.timeTo}
					error={form.errors.timeTo ?? ''}
				></Form.TextInput>
			</div>
		</Form.Root>

		{#if staticResourceData}
			<div>
				<div class="flex flex-row justify-between">
					<h1 class="mb-8 text-2xl">Tutors</h1>
					<div class="flex gap-2">
						<Dialog.Root bind:open={dialogOpen}>
							<Dialog.Trigger class={buttonVariants({ variant: 'outline' })}
								>{m.classes_addtutor()}</Dialog.Trigger
							>
							{#if dialogOpen}
								<UserAddDialog
									defaultRole="TUTOR"
									endpoint={`api/v2/courses/${page.params.courseId}/classes/${page.params.classId}/tutors`}
								></UserAddDialog>
							{/if}
						</Dialog.Root>
					</div>
				</div>
				<TutorView bind:tutors={form.fields.tutors}></TutorView>
			</div>
		{/if}
	{/await}
</div>
