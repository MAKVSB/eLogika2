<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { Input } from '../input';
	import { Label } from '../label';
	import type { WithElementRef } from 'bits-ui';
	import type { ErrorObject } from './types';

	let {
		ref = $bindable(null),
		name,
		title,
		id,
		type = 'text',
		placeholder = '',
		required = false,
		value = $bindable(''),
		error = '',
		disabled = false,
		class: className,
		...restProps
	}: {
		title?: string;
		name: string;
		id: string;
		type?: string;
		placeholder?: string;
		required?: boolean;
		value: string | number | null;
		error: string | ErrorObject;
		disabled?: boolean;
	} & WithElementRef<HTMLAttributes<HTMLElement>> = $props();
</script>

<div class="flex flex-col gap-2 {className}">
	{#if title}
		<Label for={id}>{title} {required ? '*' : ''}</Label>
	{/if}
	<Input {id} bind:value {placeholder} {required} {type} {...restProps} {disabled}/>
	{#if error}
		<p class="text-sm text-red-500">{error}</p>
	{/if}
</div>
