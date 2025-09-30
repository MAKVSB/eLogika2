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

	function add() {
		if ('clickEventHandler' in meta) {
			meta.clickEventHandler(id);
		}
	}
</script>

<div class="flex justify-between gap-2">
	<Button variant="outline" onclick={() => add()}>{m.course_users_add()}</Button>
</div>
