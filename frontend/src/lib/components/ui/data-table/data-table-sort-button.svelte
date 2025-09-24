<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import ArrowUpDownIcon from '@lucide/svelte/icons/arrow-up-down';
	import ArrowUpIcon from '@lucide/svelte/icons/arrow-up';
	import ArrowDownIcon from '@lucide/svelte/icons/arrow-down';
	import { Button } from '$lib/components/ui/button/index.js';

	let {
		variant = 'ghost',
		name,
		sorted,
		...restProps
	}: ComponentProps<typeof Button> & {
		name: string;
		sorted: false | 'asc' | 'desc';
	} = $props();
</script>

<Button {variant} {...restProps}>
	{name}
	{#if sorted === 'desc'}
		<ArrowUpIcon class="ml-2" />
	{:else if sorted === 'asc'}
		<ArrowDownIcon class="ml-2" />
	{:else}
		<ArrowUpDownIcon class="ml-2" />
	{/if}
</Button>
