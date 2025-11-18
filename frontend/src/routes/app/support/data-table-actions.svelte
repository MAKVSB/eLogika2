<script lang="ts">
	import { page } from '$app/state';
	import type { QuestionCheckedByDTO } from '$lib/api_types';
	import { Button } from '$lib/components/ui/button/index.js';
	import GlobalState from '$lib/shared.svelte';
	import { m } from '$lib/paraglide/messages';
	import { base } from '$app/paths';

	let {
		id,
		meta
	}: {
		id: number | string;
		meta: any;
	} = $props();

	function handleActionClick(event: string, params?: any) {
		if ('clickEventHandler' in meta) {
			meta.clickEventHandler(event, id, params);
		}
	}
</script>

<div class="flex justify-between gap-2">
	<Button variant="outline" class="relative" href="{base}/app/support/{id}">
		<span>{m.view()}</span>
	</Button>
	<Button variant="destructive" class="relative" onclick={() => handleActionClick('delete')}>
		<span>{m.delete()}</span>
	</Button>
</div>
