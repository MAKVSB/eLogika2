<script lang="ts" module>
	export enum MenuTreeType {
		CATEGORY,
		FLEX_GROW
	}

	export type SidebarFlex = {
		type: MenuTreeType.FLEX_GROW;
	};

	export interface SidebarCategoryItem {
		title: string;
		url: string;
		icon: Component;
		isActive?: boolean;
		requiredRoles?: string[];
	}

	export interface SidebarCategory {
		type: MenuTreeType.CATEGORY;
		name: string;
		items: SidebarCategoryItem[];
		requiredRoles?: string[];
	}
</script>

<script lang="ts">
	import type { Component } from 'svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import GlobalState from '$lib/shared.svelte';

	let { name, items, requiredRoles }: SidebarCategory = $props();
</script>

{#if !requiredRoles || requiredRoles?.length == 0 || (GlobalState.activeCourse && requiredRoles?.includes(GlobalState.activeRole ?? "ERR"))}
	<Sidebar.Group>
		<Sidebar.GroupLabel>{name}</Sidebar.GroupLabel>
		<Sidebar.Menu>
			{#each items as item (item.title)}
				{#if !item.requiredRoles || item.requiredRoles?.length == 0 || (GlobalState.activeCourse && item.requiredRoles?.includes(GlobalState.activeRole ?? "ERR"))}
					<Sidebar.MenuItem>
						<Sidebar.MenuButton>
							{#snippet child({ props })}
								<a href={item.url} {...props}>
									<item.icon />
									<span>{item.title}</span>
								</a>
							{/snippet}
						</Sidebar.MenuButton>
					</Sidebar.MenuItem>
				{/if}
			{/each}
		</Sidebar.Menu>
	</Sidebar.Group>
{/if}
