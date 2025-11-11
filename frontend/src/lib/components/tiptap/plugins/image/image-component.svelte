<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import { Label } from '$lib/components/ui/label';
	import { API } from '$lib/services/api.svelte';
	import { resize, type ResizeDetail } from '@svelte-put/resize';
	import type { NodeViewProps } from '@tiptap/core';
	import { toast } from 'svelte-sonner';
	import { NodeViewWrapper } from 'svelte-tiptap';

	let { node, updateAttributes, selected, deleteNode, editor }: NodeViewProps = $props();

	const updateSrc = (e: Event & { currentTarget?: EventTarget & HTMLInputElement }) => {
		updateAttributes({ src: e.currentTarget.value });
	};

	async function downloadFile() {
		try {
			let url = node.attrs.src as string;
			if (node.attrs.mode === 'storage') {
				url = import.meta.env.VITE_API_URL + '/api/v2/files/' + node.attrs.src;
			}

			const response = await API.request<any, Blob>(url, {
				method: 'GET'
			});

			// Create download link
			const a = document.createElement('a');
			const objectUrl = URL.createObjectURL(response);
			a.href = objectUrl;
			a.download = node.attrs.originalFilename;
			document.body.appendChild(a);
			a.click();
			// Cleanup
			URL.revokeObjectURL(objectUrl);
			a.remove();
		} catch (err) {
			toast.error('Download failed');
			console.error(err);
		}
	}

	function onResized(e: CustomEvent<ResizeDetail>) {
		updateAttributes({ width: e.detail.entry.borderBoxSize[0].inlineSize });
	}
</script>

<NodeViewWrapper class="my-4 w-full {selected && editor.options.editable ? 'border' : ''}">
	{#if selected && editor.options.editable}
		{#if (node.attrs.type as string).startsWith('image/')}
			{#if node.attrs.src}
				<div
					class="my-4 w-fit resize-x overflow-auto"
					use:resize
					onresized={onResized}
					style="width: {node.attrs.width}px"
				>
					<img
						src={node.attrs.mode === 'storage'
							? import.meta.env.VITE_API_URL + '/api/v2/files/' + node.attrs.src
							: node.attrs.src}
						alt={node.attrs.alt}
						class="my-4 w-full"
					/>
				</div>
			{:else}
				<span class="border border-red-500">No image</span>
			{/if}
			<div class="grid grid-cols-2 gap-4">
				<div class="col-span-2">
					<Label>Image source</Label>
					<Input value={node.attrs.src} disabled={node.attrs.mode != 'url'}></Input>
				</div>
				<Button variant="default" onclick={downloadFile}>Download file</Button>
				{#if editor.options.editable}
					<Button variant="destructive" onclick={deleteNode}>Delete file</Button>
				{/if}
			</div>
		{:else}
			<div class="flex items-center gap-4 rounded-2xl border p-4">
				Attachment: {node.attrs.originalFilename}
				<Button variant="default" onclick={downloadFile}>Download file</Button>
				<Button variant="destructive" onclick={deleteNode}>Delete file</Button>
			</div>
		{/if}
	{:else if (node.attrs.type as string).startsWith('image/')}
		{#if node.attrs.src}
			<div class="relative h-max w-max">
				<img
					class="border-accent border"
					src={node.attrs.mode === 'storage'
						? import.meta.env.VITE_API_URL + '/api/v2/files/' + node.attrs.src
						: node.attrs.src}
					alt={node.attrs.alt}
					style="width: {node.attrs.width}px"
					height="auto"
				/>
			</div>
		{:else}
			<span class="border border-red-500">No image</span>
		{/if}
	{:else}
		<div class="flex items-center gap-4 rounded-2xl border p-4">
			Attachment: {node.attrs.originalFilename}
			<Button variant="default" onclick={downloadFile}>Download file</Button>
		</div>
	{/if}
</NodeViewWrapper>
