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
		url?: string;
		fn?: () => void;
		icon: Component;
		isActive?: boolean;
		requiredRoles?: string[];
	}

	export interface SidebarCategory {
		type: MenuTreeType.CATEGORY;
		name: string;
		items: SidebarCategoryItem[];
		requiredRoles?: string[];
		noCourse?: boolean;
	}
</script>

<script lang="ts">
	import type { Component } from 'svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import GlobalState from '$lib/shared.svelte';

	const { setOpenMobile } = Sidebar.useSidebar();

	let { name, items, requiredRoles, noCourse }: SidebarCategory = $props();

	const canShow = (requiredRoles?: string[]): boolean => {
		if (!requiredRoles || requiredRoles?.length == 0) {
			return true;
		}

		if (!GlobalState.activeRole) {
			return false;
		}
		if (GlobalState.activeCourse || noCourse) {
			return requiredRoles?.includes(GlobalState.activeRole) ? true : false;
		}
		return false;
	};
</script>

{#if canShow(requiredRoles)}
	<Sidebar.Group>
		<Sidebar.GroupLabel>{name}</Sidebar.GroupLabel>
		<Sidebar.Menu>
			{#each items as item (item.title)}
				{#if canShow(item.requiredRoles)}
					<Sidebar.MenuItem>
						<Sidebar.MenuButton>
							{#snippet child({ props })}
								<a
									href={item.url}
									onclick={() => {
										if (item.fn) {
											item.fn();
										}
										setOpenMobile(false);
									}}
									{...props}
								>
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
