<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { TiptapMenuItemType } from './menutype';
	import type { Component } from 'svelte';

	export type TipTapMenuItemButton = {
		type: TiptapMenuItemType.BUTTON;
		icon?: Component<any>;
		command: () => void;
		title?: string;
		title_class?: string;
		disabled?: boolean;
		tooltip: string;
	};

	let {
		icon,
		command,
		title,
		title_class = '',
		disabled = false,
		tooltip
	}: TipTapMenuItemButton = $props();
</script>

<Button variant="ghost" {disabled} onclick={command} class={title_class}>
	{#if icon}
		{@const Icon = icon}
		<Icon />
	{/if}
	{#if title}
		{title}
	{/if}
</Button>
