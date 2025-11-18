<script lang="ts">
	import { API } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import * as Form from '$lib/components/ui/form';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import {
		type SupportTicketInsertRequest,
		type SupportTicketInsertResponse,
		type SupportTicketGetByIdResponse,
		type SupportTicketUpdateResponse,
		type SupportTicketUpdateRequest
	} from '$lib/api_types';
	import { TipTapDefaultContent } from '$lib/constants';
	import { m } from '$lib/paraglide/messages';
	import { SupportTicketInsertRequestSchema } from '$lib/schemas';
	import Comment from './Comment.svelte';
	import NewComment from './NewComment.svelte';

	let { data } = $props();

	let isLoaded = $state(false);

	$effect(() => {
		if (data.ticket) {
			data.ticket
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

	const defaultFormData: SupportTicketInsertRequest = {
		name: '',
		content: TipTapDefaultContent,
		solved: false,
		url: page.url.searchParams.get('origin')
			? decodeURIComponent(page.url.searchParams.get('origin')!)
			: ''
	};
	let form = $state(Form.createForm(SupportTicketInsertRequestSchema, defaultFormData));

	function setResult(
		res: SupportTicketGetByIdResponse | SupportTicketInsertResponse | SupportTicketUpdateResponse
	) {
		form.fields = res.data;
		isLoaded = true;
		goto(String(res.data.id), {
			replaceState: true
		});
	}

	async function handleSubmit(e: any) {
		let request;
		if (data.creating) {
			request = API.request<SupportTicketInsertRequest, SupportTicketInsertResponse>(
				`/api/v2/support`,
				{
					method: 'POST',
					body: form.fields
				}
			);
		} else {
			request = API.request<SupportTicketUpdateRequest, SupportTicketUpdateResponse>(
				`/api/v2/support/${page.params.id}`,
				{
					method: 'PUT',
					body: form.fields
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
</script>

<div class="m-8">
	{#await data.ticket}
		<Pageloader></Pageloader>
	{:then staticResourceData}
		{#if isLoaded}
			<div class="flex flex-row justify-between">
				<h1 class="mb-8 text-2xl">
					Support tickets:
					<b>
						{staticResourceData?.data?.name ?? 'New ticket'}
					</b>
				</h1>
			</div>
			<Form.Root
				bind:form
				onsubmit={handleSubmit}
				isCreating={data.creating}
				hideDefaultbutton={!data.creating && !staticResourceData?.data.editable}
			>
				<div class="grid grid-cols-12 gap-4">
					<Form.TextInput
						title={m.ticket_title()}
						name="name"
						id="name"
						type="text"
						class="col-span-12 sm:col-span-8"
						bind:value={form.fields.name}
						error={form.errors.name}
						disabled={!data.creating && !staticResourceData?.data.editable}
					></Form.TextInput>
					{#if !data.creating}
					<Form.Checkbox
						title={m.ticket_solved()}
						name="solved"
						id="solved"
						class="col-span-12 sm:col-span-4"
						wide={true}
						bind:value={form.fields.solved}
						error={form.errors.solved}
						disabled={!data.creating && !staticResourceData?.data.editable}
					></Form.Checkbox>
					{/if}
					<Form.TextInput
						title={m.ticket_url()}
						name="url"
						id="url"
						type="text"
						class="col-span-12"
						bind:value={form.fields.url}
						error={form.errors.url}
						disabled={!data.creating && !staticResourceData?.data.editable}
					></Form.TextInput>
					<Form.Tiptap
						title={m.ticket_content()}
						name="content"
						id="content"
						class="col-span-12"
						bind:value={form.fields.content}
						error={form.errors.content}
						disabled={!data.creating && !staticResourceData?.data.editable}
						enableFileUpload
						enableFileLink
					></Form.Tiptap>
				</div>
			</Form.Root>
		{/if}

		{#if staticResourceData}
			<div class="flex flex-col gap-4 mt-8">
				<h2 class="text-lg">{m.support_comments()}:</h2>
				{#each staticResourceData.data.comments as comment}
					<Comment {comment}></Comment>
				{/each}
				<NewComment></NewComment>
			</div>
		{/if}
	{/await}
</div>
