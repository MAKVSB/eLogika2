<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Toggle } from '$lib/components/ui/toggle';
	import { TiptapMenuItemType } from './menutype';
	import type { Component } from 'svelte';

	export type TipTapMenuItemDropdown = {
		active?: boolean;
		type: TiptapMenuItemType.DROPDOWN;
		items: TipTapMenuItemDropdownItem[];
		default: TipTapMenuItemDropdownDefaultItem;
		show_title?: boolean;
		disabled?: boolean;
		tooltip: string;
	};

	export type TipTapMenuItemDropdownItem = {
		type: TiptapMenuItemType.DROPDOWN_ITEM;
		icon?: Component<any>;
		command: () => void;
		active?: boolean;
		title?: string;
		title_class?: string;
		disabled?: boolean;
	};

	export type TipTapMenuItemDropdownDefaultItem = {
		type: TiptapMenuItemType.DROPDOWN_ITEM;
		icon: Component<any>;
		title?: string;
		title_class?: string;
		disabled?: boolean;
	};

	let {
		active = false,
		items,
		default: def,
		disabled = false,
		show_title = false,
		tooltip
	}: TipTapMenuItemDropdown = $props();
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger disabled={disabled || (!active && def.disabled)}>
		{#if active}
			<Toggle bind:pressed={active}>
				{@const toggleItems = items.filter((item) => item.active)}
				{@const toggleItem = toggleItems.length > 0 ? toggleItems[0] : def}

				{#if toggleItem.icon}
					{@const Icon = toggleItem.icon}
					<Icon />
				{/if}
				{#if show_title && toggleItem.title}
					<span class={toggleItem.title_class}>{toggleItem.title}</span>
				{/if}
			</Toggle>
		{:else}
			<Toggle disabled={def.disabled}>
				{#if def.icon}
					{@const Icon = def.icon}
					<Icon />
				{/if}
				{#if show_title && def.title}
					<span class={def.title_class}>{def.title}</span>
				{/if}
			</Toggle>
		{/if}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content>
		<DropdownMenu.Group>
			{#each items as item}
				<DropdownMenu.Item
					variant={item.active ? 'highlighted' : 'default'}
					onclick={item.command}
					disabled={item.disabled}
				>
					{#if item.icon}
						{@const Icon = item.icon}
						<Icon />
					{/if}
					{#if item.title}
						<span class={item.title_class}>{item.title}</span>
					{/if}
				</DropdownMenu.Item>
			{/each}
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>
