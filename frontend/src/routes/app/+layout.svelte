<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import AppHeader from '$lib/components/app-header.svelte';
	import AppSidebar from '$lib/components/app-sidebar.svelte';
	import { API } from '$lib/services/api.svelte';
	import { onMount } from 'svelte';
	import type { UserCourseListResponse } from '$lib/api_types';
	import { toast } from 'svelte-sonner';
	import GlobalState from '$lib/shared.svelte';

	let loaded = $state(false);

	const loadData = () => {
		if (GlobalState.availableCourses.length != 0) {
			loaded = true;
		}
		API.request<null, UserCourseListResponse>(`/api/v2/courses/user`, {})
			.then((res) => {
				GlobalState.availableCourses = res.items;
				loaded = true;
			})
			.catch((err) => {
				toast.error('Failed to load available courses');
			});
	};

	onMount(() => {
		loadData();
	});

	let { children } = $props();
</script>

{#if loaded}
	{#key GlobalState.reloadCounter}
		<div class="[--header-height:calc(--spacing(14))]">
			<Sidebar.Provider class="flex flex-col">
				<AppHeader />
				<div class="flex flex-1">
					<AppSidebar />
					<Sidebar.Inset>
						<div class="flex flex-col flex-1 w-full mx-auto contain-inline-size">
							{@render children()}
						</div>
					</Sidebar.Inset>
				</div>
			</Sidebar.Provider>
		</div>
	{/key}
{/if}
