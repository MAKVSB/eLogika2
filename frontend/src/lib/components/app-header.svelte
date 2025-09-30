<script lang="ts">
	import SidebarIcon from '@lucide/svelte/icons/sidebar';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import Logo from '$lib/images/logo.svg';
	import ModeToggle from './mode-toggle.svelte';
	import LangToggle from './lang-toggle.svelte';
	import { base } from '$app/paths';

	const sidebar = Sidebar.useSidebar();
</script>

<header class="sticky top-0 z-50 flex items-center w-full border-b bg-background">
	<div class="flex h-(--header-height) w-full items-center gap-2 px-4">
		<Button class="size-8" variant="ghost" size="icon" onclick={sidebar.toggle}>
			<SidebarIcon />
		</Button>
		<Separator orientation="vertical" class="h-4 mr-2" />
		<a href="/{base}" class="flex ms-2 md:me-24">
			<img src={Logo} alt="eLogika logo" class="h-8 me-3" />
			<span class="self-center text-xl font-semibold whitespace-nowrap sm:text-2xl dark:text-white">
				eLogika
			</span>
		</a>
		<div class="flex-1"></div>
		<LangToggle></LangToggle>
		<ModeToggle></ModeToggle>
	</div>
</header>
