<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { m } from '$lib/paraglide/messages.js';

	let { id }: { id: number | string } = $props();
</script>

<div class="flex justify-between">
	<Button variant="ghost" class="relative" href="/app/{id}/tutor/course">
		<span>{m.edit()}</span>
	</Button>

	<Button variant="ghost" class="relative" href="/app/{id}/tutor/course/clone">
		<span>{m.clone()}</span>
	</Button>
</div>
