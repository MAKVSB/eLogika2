<script lang="ts">
	import ChevronsUpDownIcon from '@lucide/svelte/icons/chevrons-up-down';
	import LogOutIcon from '@lucide/svelte/icons/log-out';
	import SettingsIcon from '@lucide/svelte/icons/settings';

	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { LoggedUserDTO, LogoutResponse } from '$lib/api_types';
	import { API } from '$lib/services/api.svelte';

	import GlobalState from '$lib/shared.svelte';
	import { m } from '$lib/paraglide/messages';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';

	let {
		user = $bindable(undefined)
	}: {
		user: LoggedUserDTO | undefined;
	} = $props();

	const sidebar = Sidebar.useSidebar();

	async function logout(all: boolean = false) {
		let url = '/api/v2/auth/logout';
		if (all) {
			url += '/all';
		}

		await API.request<null, LogoutResponse>(url, {
			method: 'POST',
			credentials: 'include'
		});
		GlobalState.loggedUser = null;
		console.log('Transfering 10');
		goto(base + '/');
	}
</script>

{#if user}
	<Sidebar.Menu>
		<Sidebar.MenuItem>
			<DropdownMenu.Root>
				<DropdownMenu.Trigger>
					{#snippet child({ props })}
						<Sidebar.MenuButton
							size="lg"
							class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
							{...props}
						>
							<Avatar.Root class="rounded-lg size-8">
								<Avatar.Fallback class="rounded-lg">
									{user.firstName.charAt(0)}{user.familyName.charAt(0)}
								</Avatar.Fallback>
							</Avatar.Root>
							<div class="grid flex-1 text-sm leading-tight text-left">
								<span class="font-medium truncate">{user.firstName} {user.familyName}</span>
								<span class="text-xs truncate">{user.email}</span>
							</div>
							<ChevronsUpDownIcon class="ml-auto size-4" />
						</Sidebar.MenuButton>
					{/snippet}
				</DropdownMenu.Trigger>
				<DropdownMenu.Content
					class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
					side={sidebar.isMobile ? 'bottom' : 'right'}
					align="end"
					sideOffset={4}
				>
					<DropdownMenu.Label class="p-0 font-normal">
						<div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
							<Avatar.Root class="rounded-lg size-8">
								<Avatar.Fallback class="rounded-lg">
									{user.firstName.charAt(0)}{user.familyName.charAt(0)}
								</Avatar.Fallback>
							</Avatar.Root>
							<div class="grid flex-1 text-sm leading-tight text-left">
								<span class="font-medium truncate">{user.firstName} {user.familyName}</span>
								<span class="text-xs truncate">{user.email}</span>
							</div>
						</div>
					</DropdownMenu.Label>
					<DropdownMenu.Separator />
					<DropdownMenu.Item onclick={() => goto(base+"/app/user")}>
						<SettingsIcon />
						Profile
					</DropdownMenu.Item>
					<DropdownMenu.Item onclick={() => logout(false)}>
						<LogOutIcon />
						{m.logout()}
					</DropdownMenu.Item>
					<DropdownMenu.Item onclick={() => logout(true)}>
						<LogOutIcon />
						{m.logout_all()}
					</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</Sidebar.MenuItem>
	</Sidebar.Menu>
{/if}
