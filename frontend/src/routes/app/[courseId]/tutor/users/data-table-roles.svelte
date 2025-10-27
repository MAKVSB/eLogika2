<script lang="ts">
	import { CourseUserRoleEnum } from '$lib/api_types';
	import { Button } from '$lib/components/ui/button/index.js';
	import { m } from '$lib/paraglide/messages';
	import GlobalState from '$lib/shared.svelte';

	let {
		id,
		meta,
		roles
	}: {
		id: number | string;
		meta: any;
		roles: string[];
	} = $props();

	function handleActionClick(event: string, params?: any) {
		if ('clickEventHandler' in meta) {
			meta.clickEventHandler(event, id, params);
		}
	}
</script>

<div class="flex gap-2">
	{#if GlobalState.activeRole == CourseUserRoleEnum.ADMIN}
		{#each roles as role}
			<Button variant="destructive" onclick={() => handleActionClick('remove_role', { role })}
				>{m.course_user_role_enum({ value: role })}</Button
			>
		{/each}
	{:else}
		{roles.map((role) => m.course_user_role_enum({ value: role })).join(', ')}
	{/if}
</div>
