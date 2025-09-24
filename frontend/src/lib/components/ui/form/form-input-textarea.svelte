<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { Input } from '../input';
	import { Label } from '../label';
	import type { WithElementRef } from 'bits-ui';
	import type { ErrorObject } from './types';
	import Textarea from '../textarea/textarea.svelte';

	let {
		ref = $bindable(null),
		name,
		title,
		id,
		placeholder = '',
		required = false,
		value = $bindable(''),
		error = '',
		class: className,
		...restProps
	}: {
		title?: string;
		name: string;
		id: string;
		placeholder?: string;
		required?: boolean;
		value: string;
		error: string | ErrorObject;
	} & WithElementRef<HTMLAttributes<HTMLElement>> = $props();
</script>

<div class="flex flex-col gap-2 {className}">
	{#if title}
		<Label for={id}>{title} {required ? '*' : ''}</Label>
	{/if}
	<Textarea {id} bind:value {placeholder} {required} {...restProps}></Textarea>
	{#if error}
		<p class="text-sm text-red-500">{error}</p>
	{/if}
</div>
