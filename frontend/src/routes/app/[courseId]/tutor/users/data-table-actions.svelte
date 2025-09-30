<script lang="ts">
	import { CourseUserRoleEnum } from '$lib/api_types';
	import { Button } from '$lib/components/ui/button/index.js';
	import { m } from '$lib/paraglide/messages';
	import GlobalState from '$lib/shared.svelte';

	let {
		id,
		meta
	}: {
		id: number | string;
		meta: any;
	} = $props();

	function remove() {
		if ('clickEventHandler' in meta) {
			meta.clickEventHandler(id);
		}
	}
</script>

<div class="flex justify-between gap-2">
	{#if GlobalState.activeRole == CourseUserRoleEnum.ADMIN}
		<Button variant="outline" onclick={() => remove()}>{m.remove()}</Button>
	{/if}
</div>
