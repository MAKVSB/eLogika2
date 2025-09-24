<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { API, ApiError } from '$lib/services/api.svelte';
	import * as Form from '$lib/components/ui/form';
	import { toast } from 'svelte-sonner';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import { TipTapDefaultContent } from '$lib/constants';
	import { m } from '$lib/paraglide/messages.js';
	import type {
		ActivityInstanceDTO,
		ActivityInstanceGetResponse,
		ActivityInstanceSaveResponse
	} from '$lib/api_types';
	import {
		ActivityInstanceSaveRequestSchema,
		type ActivityInstanceSaveRequest
	} from '$lib/schemas';

	let courseId = $derived<string | null>(page.params.courseId);
	let { data } = $props();

	$effect(() => {
		data.activityInstance.then((data) => setResult(data)).catch(() => {});
	});

	const defaultFormData: ActivityInstanceDTO = {
		id: 0,
		content: TipTapDefaultContent,

		assignmentName: 'string',
		assignmentDescription: TipTapDefaultContent,
		assignmentExpectedResult: TipTapDefaultContent,

		points: 0,
		pointsMax: 0,
		pointsMin: 0,
		editable: true
	};
	let form = $state(Form.createForm(ActivityInstanceSaveRequestSchema, defaultFormData));

	function setResult(res: ActivityInstanceGetResponse | ActivityInstanceSaveResponse) {
		form.fields = res.instanceData;
		console.log("Transfering 18")
		goto(String(res.instanceData.id), {
			replaceState: true
		});
	}

	async function handleSubmit(): Promise<any> {
		let request = API.request<ActivityInstanceSaveRequest, ActivityInstanceSaveResponse>(
			`/api/v2/courses/${page.params.courseId}/activities/instance/${page.params.instanceId}`,
			{
				method: 'PUT',
				body: {
					...form.fields
				}
			}
		);

		return request
			.then((res) => {
				setResult(res);
				toast.success('Saved');
			});
	}
</script>

<div class="m-8">
	{#await data.activityInstance}
		<Pageloader></Pageloader>
	{:then staticResourceData}
		<div class="flex flex-row justify-between">
			<h1 class="mb-8 text-2xl">
				Activity:
				<b>
					{staticResourceData?.instanceData?.assignmentName}
				</b>
			</h1>
		</div>
		<Form.Root bind:form onsubmit={handleSubmit}>
			{#key form.fields}
				<div class="flex flex-col gap-4 p-2">
					<div class="grid grid-cols-4 gap-4">
						<Form.Tiptap
							title="Assignment"
							name="content"
							id="content"
							class="col-span-4"
							bind:value={form.fields.assignmentDescription}
							error={form.errors.assignmentDescription}
							disabled
						></Form.Tiptap>
						<Form.Tiptap
							title="Expected result"
							name="content"
							id="content"
							class="col-span-4"
							bind:value={form.fields.assignmentExpectedResult}
							error={form.errors.assignmentExpectedResult}
							disabled
						></Form.Tiptap>
						<Form.Tiptap
							title="Submission"
							name="content"
							id="content"
							class="col-span-4"
							bind:value={form.fields.content}
							error={form.errors.content}
							disabled={!form.fields.editable}
							enableFileUpload
							enableFileLink
						></Form.Tiptap>
					</div>
				</div>
			{/key}
		</Form.Root>
	{/await}
</div>
