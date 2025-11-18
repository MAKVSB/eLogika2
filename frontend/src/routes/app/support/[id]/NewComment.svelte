<script lang="ts">
	import { API } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import * as Form from '$lib/components/ui/form';
	import { toast } from 'svelte-sonner';
	import {
		type SupportTicketCommentInsertRequest,
		type SupportTicketCommentInsertResponse,
	} from '$lib/api_types';
	import { TipTapDefaultContent } from '$lib/constants';
	import { m } from '$lib/paraglide/messages';
	import { SupportTicketCommentInsertRequestSchema } from '$lib/schemas';
	import { invalidate } from '$app/navigation';

	const defaultFormData: SupportTicketCommentInsertRequest = {
		content: TipTapDefaultContent,
	};
	let form = $state(Form.createForm(SupportTicketCommentInsertRequestSchema, defaultFormData));

	async function handleSubmit(e: any) {
		let request = API.request<SupportTicketCommentInsertRequest, SupportTicketCommentInsertResponse>(
				`/api/v2/support/${page.params.id}/comment`,
				{
					method: 'PUT',
					body: form.fields
				}
			);

		await request
			.then(() => {
				invalidate((url) => {
					return url.href.endsWith(`/api/v2/support/${page.params.id}`);
				});
				toast.success('Saved');
			})
			.catch(() => {});
	}
</script>

<Form.Root bind:form onsubmit={handleSubmit} isCreating>
	<div class="grid grid-cols-12 gap-4">
		<Form.Tiptap
			title={m.support_comment_new_content()}
			name="content"
			id="content"
			class="col-span-12"
			bind:value={form.fields.content}
			error={form.errors.content}
			enableFileUpload
			enableFileLink
		></Form.Tiptap>
	</div>
</Form.Root>
