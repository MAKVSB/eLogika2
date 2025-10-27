<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { m } from '$lib/paraglide/messages';

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
	<Button variant="destructive" class="relative" onclick={() => handleActionClick('revoke')}>
		<span>{m.user_token_revoke()}</span>
	</Button>
</div>
