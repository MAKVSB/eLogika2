<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { Label } from '../label';
	import type { WithElementRef } from 'bits-ui';
	import Checkbox from '../checkbox/checkbox.svelte';
	import { cn } from '$lib/utils';
	import type { ErrorObject } from './types';

	let {
		ref = $bindable(null),
		name,
		title,
		disabled,
		id,
		class: className,
		innerClass,
		type = 'text',
		required = false,
		value = $bindable(false),
		error = $bindable(''),
		wide = false,
		...restProps
	}: {
		title?: string;
		name: string;
		id: string;
		type?: string;
		class?: string;
		innerClass?: string;
		required?: boolean;
		disabled?: boolean;
		value: boolean;
		error: string | ErrorObject;
		wide?: boolean;
	} & WithElementRef<HTMLAttributes<HTMLElement>> = $props();
</script>

<div class={cn('flex flex-col gap-2', className, wide ? 'w-full' : '')}>
	{#if title}
		<Label for={id}>{title} {required ? '*' : ''}</Label>
	{/if}
	<Checkbox
		class={cn('h-9 w-full rounded-md', wide ? 'w-full' : 'w-9', innerClass)}
		{id}
		bind:checked={value}
		{required}
		{disabled}
		{...restProps}
	/>
	{#if error}
		<p class="text-sm text-red-500">{error}</p>
	{/if}
</div>
