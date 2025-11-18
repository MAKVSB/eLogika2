<script lang="ts">
	import { type SupportTicketCommentDTO } from '$lib/api_types';
	import TiptapRenderer from '$lib/components/tiptap-renderer/tiptap-renderer.svelte';
	import { m } from '$lib/paraglide/messages';
	import { displayUserName } from '$lib/utils';

	let {
		comment
	}: {
		comment: SupportTicketCommentDTO;
	} = $props();
</script>

{#if comment.content}
	<div class="flex flex-col gap-4 p-4 border-2">
		{m.ticket_comment_sentby({
			userString: displayUserName(comment.createdBy),
			createdAt:new Date(comment.createdAt).toLocaleString('cs', {
			dateStyle: 'short',
			timeStyle: 'short'
		})
		})}:
		<TiptapRenderer jsonContent={comment.content}></TiptapRenderer>
	</div>
{/if}
