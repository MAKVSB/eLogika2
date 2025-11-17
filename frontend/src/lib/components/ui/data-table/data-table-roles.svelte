<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { m } from '$lib/paraglide/messages';

	let {
		id,
		meta,
		roles,
		showButtons = false
	}: {
		id: number | string;
		meta: any;
		roles: string[];
		showButtons?: boolean;
	} = $props();

	let sortedRoles = $derived(
		[...roles].sort((a, b) => {
			const orderA = roleOrder[a];
			const orderB = roleOrder[b];

			// If both roles are found, compare them numerically.
			return orderA - orderB;
		})
	);

	const roleOrder: Record<string, number> = {
		STUDENT: 1,
		TUTOR: 2,
		GARANT: 3
	};

	function handleActionClick(event: string, params?: any) {
		if ('clickEventHandler' in meta) {
			meta.clickEventHandler(event, id, params);
		}
	}
</script>

<div class="flex gap-2">
	{#if showButtons}
		{#each sortedRoles as role}
			<Button variant="destructive" onclick={() => handleActionClick('remove_role', { role })}
				>{m.course_user_role_enum({ value: role })}</Button
			>
		{/each}
	{:else}
		{roles.map((role) => m.course_user_role_enum({ value: role })).join(', ')}
	{/if}
</div>
