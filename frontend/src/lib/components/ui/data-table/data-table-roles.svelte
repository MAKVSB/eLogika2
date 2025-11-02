<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { m } from '$lib/paraglide/messages';

	let {
		id,
		meta,
		roles,
		showButtons = false,
	}: {
		id: number | string;
		meta: any;
		roles: string[];
		showButtons?: boolean;
	} = $props();

	function handleActionClick(event: string, params?: any) {
		if ('clickEventHandler' in meta) {
			meta.clickEventHandler(event, id, params);
		}
	}
</script>

<div class="flex gap-2">
	{#if showButtons}
		{#each roles as role}
			<Button variant="destructive" onclick={() => handleActionClick('remove_role', { role })}
				>{m.course_user_role_enum({ value: role })}</Button
			>
		{/each}
	{:else}
		{roles.map((role) => m.course_user_role_enum({ value: role })).join(', ')}
	{/if}
</div>
