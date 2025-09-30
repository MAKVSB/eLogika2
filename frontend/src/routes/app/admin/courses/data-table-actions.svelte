<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { m } from '$lib/paraglide/messages.js';
	import { base } from '$app/paths';

	let { id }: { id: number | string } = $props();
</script>

<div class="flex justify-between gap-2">
	<Button variant="outline" class="relative" href="{base}/app/{id}/tutor/course">
		<span>{m.edit()}</span>
	</Button>

	<Button variant="outline" class="relative" href="{base}/app/{id}/tutor/course/clone">
		<span>{m.clone()}</span>
	</Button>
</div>
