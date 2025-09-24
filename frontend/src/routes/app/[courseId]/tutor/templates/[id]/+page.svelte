<script lang="ts">
	import { API, ApiError } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import * as Form from '$lib/components/ui/form';
	import { goto } from '$app/navigation';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import {
		type TemplateDTO,
		type TemplateGetByIdResponse,
		type TemplateUpdateRequest,
		type TemplateUpdateResponse,
		type TemplateInsertResponse,
		type TemplateInsertRequest
	} from '$lib/api_types';
	import Blocks from './Blocks.svelte';
	import { m } from '$lib/paraglide/messages';
	import GlobalState from '$lib/shared.svelte';
	import { TemplateInsertRequestSchema } from '$lib/schemas';

	let courseId = $derived<string | null>(page.params.courseId);
	let { data } = $props();

	$effect(() => {
		if (data.template) {
			data.template.then((data) => setResult(data)).catch(() => {});
		} else {
			form.fields = defaultFormData;
		}
	});

	const defaultFormData = {
		id: 0,
		version: 0,
		title: '',
		description: '',
		mixBlocks: false,
		mixEverything: false,
		blocks: [],
		createdBy: GlobalState.loggedUser!
	};
	let form = $state(Form.createForm(TemplateInsertRequestSchema, defaultFormData));

	function setResult(
		res: TemplateGetByIdResponse | TemplateInsertResponse | TemplateUpdateResponse
	) {
		form.fields = res.data;
		console.log("Transfering 31")
		goto(String(res.data.id), {
			replaceState: true
		});
	}

	async function handleSubmit(): Promise<any> {
		let request;
		if (data.creating) {
			request = API.request<TemplateInsertRequest, TemplateInsertResponse>(
				`/api/v2/courses/${courseId}/templates`,
				{
					method: 'POST',
					body: form.fields
				}
			);
		} else {
			request = API.request<TemplateUpdateRequest, TemplateUpdateResponse>(
				`/api/v2/courses/${courseId}/templates/${page.params.id}`,
				{
					method: 'PUT',
					body: form.fields
				}
			);
		}
		return request.then((res) => setResult(res));
	}
</script>

<div class="m-8">
	{#await data.template}
		<Pageloader></Pageloader>
	{:then staticResourceData}
		<div class="flex flex-row justify-between">
			<h1 class="mb-8 text-2xl">
				Template management:
				<b>
					{staticResourceData?.data?.title ?? 'New template'}
				</b>
			</h1>
		</div>
		<Form.Root bind:form onsubmit={handleSubmit} isCreating={data.creating}>
			<div class="flex flex-col gap-4 p-2 my-4">
				<div class="grid grid-cols-12 gap-4">
					<Form.TextInput
						title={m.template_title()}
						name="title"
						id="title"
						type="text"
						class="col-span-12 sm:col-span-6"
						bind:value={form.fields.title}
						error={form.errors.title}
					></Form.TextInput>
					<Form.Checkbox
						title={m.template_mixblocks()}
						name="mixBlocks"
						id="mixBlocks"
						class="col-span-6 sm:col-span-3"
						bind:value={form.fields.mixBlocks}
						error={form.errors.mixBlocks}
					></Form.Checkbox>
					<Form.Checkbox
						title={m.template_mixeverything()}
						name="mixEverything"
						id="mixEverything"
						class="col-span-6 sm:col-span-3"
						bind:value={form.fields.mixEverything}
						error={form.errors.mixEverything}
					></Form.Checkbox>
					<Form.TextArea
						title={m.template_description()}
						name="description"
						id="description"
						class="col-span-12 sm:col-span-12"
						bind:value={form.fields.description}
						error={form.errors.description}
					></Form.TextArea>
				</div>
			</div>

			<Blocks bind:form></Blocks>
		</Form.Root>
	{/await}
</div>
