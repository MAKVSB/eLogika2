<script lang="ts">
	import SidebarIcon from '@lucide/svelte/icons/sidebar';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import Logo from '$lib/images/logo.svg';
	import ModeToggle from './mode-toggle.svelte';
	import LangToggle from './lang-toggle.svelte';

	const sidebar = Sidebar.useSidebar();
</script>

<header class="bg-background sticky top-0 z-50 flex w-full items-center border-b">
	<div class="flex h-(--header-height) w-full items-center gap-2 px-4">
		<Button class="size-8" variant="ghost" size="icon" onclick={sidebar.toggle}>
			<SidebarIcon />
		</Button>
		<Separator orientation="vertical" class="mr-2 h-4" />
		<a href="/" class="ms-2 flex md:me-24">
			<img src={Logo} alt="eLogika logo" class="me-3 h-8" />
			<span class="self-center text-xl font-semibold whitespace-nowrap sm:text-2xl dark:text-white">
				eLogika
			</span>
		</a>
		<div class="flex-1"></div>
		<LangToggle></LangToggle>
		<ModeToggle></ModeToggle>
	</div>
</header>
