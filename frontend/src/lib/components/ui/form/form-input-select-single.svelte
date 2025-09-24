<script lang="ts" module>
	export type SelectValue = {
		value: string | undefined;
		display: string | undefined;
	};

	export type SelectOption = {
		value: string | number;
		display: string;
		enabled?: boolean;
		tooltip?: string
	};

	export type SelectOptions = SelectOption[];
</script>

<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { Label } from '../label';
	import type { WithElementRef } from 'bits-ui';
	import * as Select from '../select';
	import type { ErrorObject } from './types';
	import Loader from '../loader/loader.svelte';

	let {
		ref = $bindable(null),
		name,
		title,
		id,
		placeholder = '',
		required = false,
		value = $bindable(),
		options = $bindable([]),
		error = '',
		class: className,
		nullable,
		nullableText = 'Nothing',
		loading = $bindable(false),
		onchange,
		disabled = false,
		...restProps
	}: {
		title?: string;
		name: string;
		id: string;
		value: string | number | undefined | null;
		placeholder?: string;
		required?: boolean;
		options: SelectOptions;
		error: string | ErrorObject;
		nullable?: boolean;
		nullableText?: string;
		loading?: boolean;
		disabled?: boolean;
	} & WithElementRef<HTMLAttributes<HTMLElement>> = $props();

	let triggerDisplay: string | undefined = $derived.by(() => {
		if (value) {
			const found = options.find((v) => {
				return value === v.value;
			});
			if (found) {
				return found.display;
			} else {
				return String(value);
			}
		} else {
			return undefined;
		}
	});
</script>

<div class="flex flex-col gap-2 {className}">
	{#if title}
		<Label for={id}>{title} {required ? '*' : ''}</Label>
	{/if}
	<Select.Root type="single" bind:value={triggerDisplay} {required} {...restProps} {disabled}>
		<Select.Trigger {id} class="w-full overflow-hidden">
			{#if loading}
				<Loader class="mx-auto"></Loader>
			{:else}
				{triggerDisplay !== '' ? triggerDisplay : placeholder}
			{/if}
		</Select.Trigger>
		<Select.Content>
			{#if nullable}
				<Select.Item
					value="none"
					onclick={(e) => {
						value = null;
						onchange?.(e);
					}}
				>
					<span class="text-gray-500">{nullableText}</span>
				</Select.Item>
			{/if}
			{#each options as option}
				<Select.Item
					value={option.display}
					onclick={(e) => {
						value = option.value;
						onchange?.(e);
					}}
					disabled={option.enabled == false}
				>
					{#if option.tooltip}
					{option.display} ({option.tooltip})
					{:else}
						{option.display}
					{/if}
				</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>
	{#if error}
		<p class="text-sm text-red-500">{error}</p>
	{/if}
</div>
