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
		placeholder = '',
		required = false,
		value = $bindable(new Date().getFullYear()),
		error = '',
		class: className,
		...restProps
	}: {
		title?: string;
		name: string;
		id: string;
		type?: string;
		placeholder?: string;
		required?: boolean;
		value: number;
		error: string | ErrorObject;
	} & WithElementRef<HTMLAttributes<HTMLElement>> = $props();

	let nextYear = $derived(Number(value) + 1);
</script>

<div class="flex flex-col gap-2 {className}">
	{#if title}
		<Label for={id}>{title} {required ? '*' : ''}</Label>
	{/if}
	<div class="flex items-center gap-2">
		<Input {id} bind:value {placeholder} {required} type="number" {...restProps} />
		/
		<Input {id} bind:value={nextYear} type="number" {...restProps} disabled />
	</div>
	{#if error}
		<p class="text-sm text-red-500">{error}</p>
	{/if}
</div>
