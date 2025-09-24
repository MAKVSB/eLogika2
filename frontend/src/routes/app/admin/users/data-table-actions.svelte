<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { m } from '$lib/paraglide/messages';

	let { id }: { id: number | string } = $props();
</script>

<div class="flex justify-between">
	<Button variant="ghost" class="relative" href="/app/admin/users/{id}">
		<span>{m.edit()}</span>
	</Button>
</div>
