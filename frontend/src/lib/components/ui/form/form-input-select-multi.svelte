<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { Label } from '../label';
	import type { WithElementRef } from 'bits-ui';
	import * as Select from '../select';
	import type { ErrorObject } from './types';

	let {
		ref = $bindable(null),
		name,
		title,
		id,
		placeholder = '',
		required = false,
		value = $bindable([]),
		options = $bindable([]),
		error = $bindable(''),
		class: className,
		...restProps
	}: {
		title?: string;
		name: string;
		id: string;
		value: string[];
		placeholder?: string;
		required?: boolean;
		options: string[];
		error: string | ErrorObject;
	} & WithElementRef<HTMLAttributes<HTMLElement>> = $props();
</script>

<div class="flex flex-col gap-2 {className}">
	{#if title}
		<Label for={id}>{title} {required ? '*' : ''}</Label>
	{/if}
	<Select.Root type="multiple" bind:value {required} {...restProps}>
		<Select.Trigger {id}>
			{value ?? placeholder}
		</Select.Trigger>
		<Select.Content>
			<!-- TODO -->
			<!-- {#each options as option}
				<Select.Item
					value={option}
					onclick={() => {
						value = option;
					}}
					>{option}
				</Select.Item>
			{/each} -->
		</Select.Content>
	</Select.Root>
	{#if error}
		<p class="text-sm text-red-500">{error}</p>
	{/if}
</div>
