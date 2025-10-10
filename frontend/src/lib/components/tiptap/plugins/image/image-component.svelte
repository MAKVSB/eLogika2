<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import { Label } from '$lib/components/ui/label';
	import { API } from '$lib/services/api.svelte';
	import type { NodeViewProps } from '@tiptap/core';
	import { toast } from 'svelte-sonner';
	import { NodeViewWrapper } from 'svelte-tiptap';

	let { node, updateAttributes, selected, deleteNode, editor }: NodeViewProps = $props();

	const updateSrc = (e: Event & { currentTarget?: EventTarget & HTMLInputElement }) => {
		updateAttributes({ src: e.currentTarget.value });
	};

	let width = $state(300); // initial width in px
	let isResizing = false;
	let imageEl: HTMLElement | undefined = $state();

	function startResize(event: MouseEvent) {
		isResizing = true;
		event.preventDefault();
	}

	function doResize(event: MouseEvent) {
		if (isResizing && imageEl) {
			console.log(event);
			width = event.clientX - imageEl.getBoundingClientRect().left;
			updateAttributes({ width: width > 100 ? width : 100 });
		}
	}

	function stopResize() {
		isResizing = false;
	}

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
</script>

<svelte:window onmouseup={stopResize} onmousemove={doResize} />

<NodeViewWrapper class="my-4 w-full {selected && editor.options.editable ? 'border' : ''}">
	{#if selected && editor.options.editable}
		{#if (node.attrs.type as string).startsWith('image/')}
			{#if node.attrs.src}
				<div class="relative h-max w-max" bind:this={imageEl}>
					<img
						class="border border-accent"
						src={node.attrs.mode === 'storage'
							? import.meta.env.VITE_API_URL + '/api/v2/files/' + node.attrs.src
							: node.attrs.src}
						alt={node.attrs.alt}
						style="width: {node.attrs.width}px"
						height="auto"
					/>
					<div
						class="absolute top-0 right-0 z-10 w-3 h-full bg-accent cursor-ew-resize"
						onmousedown={startResize}
						role="button"
						tabindex="0"
					></div>
				</div>
			{:else}
				<span class="border border-red-500">No image</span>
			{/if}
			{node.attrs.src}
			<div class="flex flex-col gap-4">
				{#if node.attrs.mode == 'url'}
					<div>
						<Label>Image source</Label>
						<Input value={node.attrs.src} onchange={updateSrc}></Input>
					</div>
				{/if}
				<div>
					<Button variant="default" onclick={downloadFile}>Download file</Button>
					{#if editor.options.editable}
						<Button variant="destructive" onclick={deleteNode}>Delete file</Button>
					{/if}
				</div>
			</div>
		{:else}
			<div class="flex items-center gap-4 p-4 border rounded-2xl">
				Attachment: {node.attrs.originalFilename}
				<Button variant="default" onclick={downloadFile}>Download file</Button>
				<Button variant="destructive" onclick={deleteNode}>Delete file</Button>
			</div>
		{/if}
	{:else}
		{#if (node.attrs.type as string).startsWith("image/")}
			{#if node.attrs.src}
				<div class="relative h-max w-max">
					<img
						class="border border-accent"
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
			<div class="flex items-center gap-4 p-4 border rounded-2xl">
				Attachment: {node.attrs.originalFilename}
				<Button variant="default" onclick={downloadFile}>Download file</Button>
			</div>
		{/if}
	{/if}
</NodeViewWrapper>
