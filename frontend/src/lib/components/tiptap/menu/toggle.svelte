<script lang="ts">
	import { Toggle } from '$lib/components/ui/toggle';
	import { TiptapMenuItemType } from './menutype';
	import type { Component } from 'svelte';

	export type TipTapMenuItemToggle = {
		type: TiptapMenuItemType.TOGGLE;
		tooltip: string;
		icon?: Component<any>;
		command: () => void;
		active?: boolean;
		title?: string;
		title_class?: string;
		disabled?: boolean;
		visible?: boolean;
	};

	let {
		icon,
		command,
		active = false,
		title,
		title_class = '',
		disabled = false,
		tooltip,
		visible = true
	}: TipTapMenuItemToggle = $props();
</script>

{#if visible}
	<Toggle bind:pressed={active} onclick={command} class={title_class} {disabled}>
		{#if icon}
			{@const Icon = icon}
			<Icon />
		{/if}
		{#if title}
			{title}
		{/if}
	</Toggle>
{/if}
