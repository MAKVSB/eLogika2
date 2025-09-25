<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button/index.js';
	import { m } from '$lib/paraglide/messages';
	import { base } from '$app/paths';

	let { id }: { id: number | string } = $props();
</script>

<div class="flex justify-between">
	<Button variant="ghost" class="relative" href="{base}/app/{page.params.courseId}/tutor/chapters/{id}">
		<span>{m.edit()}</span>
	</Button>
</div>
