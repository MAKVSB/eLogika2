<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { Label } from '../label';
	import type { WithElementRef } from 'bits-ui';
	import Tiptap from '../../tiptap/Tiptap.svelte';
	import type { ErrorObject } from './types';

	// TODO cleanup
	let {
		ref = $bindable(null),
		name,
		title,
		id,
		required = false,
		value = $bindable(),
		error = '',
		class: className,
		disabled = false,
		...restProps
	}: {
		title?: string;
		name: string;
		id: string;
		value: any; // TODO replace any with correct type for tiptap
		required?: boolean;
		error: string | ErrorObject;
		disabled?: boolean;
		enableFileUpload?: boolean;
		enableFileLink?: boolean;
		enabledExtensions?: string[];
	} & WithElementRef<HTMLAttributes<HTMLElement>> = $props();
</script>

<div class="flex flex-col gap-2 {className}">
	{#if title}
		<Label for={id}>{title}</Label>
	{/if}
	<Tiptap {id} {name} {required} bind:value {...restProps} {disabled}></Tiptap>
	{#if error}
		<p class="text-sm text-red-500">{error}</p>
	{/if}
</div>
