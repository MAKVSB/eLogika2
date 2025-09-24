<script lang="ts">
	import { Label } from '../label';
	import type { WithElementRef } from 'bits-ui';
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';
	import type { ErrorObject } from './types';

	let {
		ref = $bindable(null),
		title,
		id,
		class: className,
		required = false,
		error = $bindable(''),
		children
	}: WithElementRef<{
		title?: string;
		id: string;
		class?: string;
		required?: boolean;
		error: string | ErrorObject;
		children: Snippet;
	}> = $props();
</script>

<div class={cn('flex w-full flex-col gap-2', className)}>
	{#if title}
		<Label for={id}>{title} {required ? '*' : ''}</Label>
	{/if}

	{@render children?.()}

	{#if error}
		<p class="text-sm text-red-500">{error}</p>
	{/if}
</div>
