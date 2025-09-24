<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button/index.js';
	import { m } from '$lib/paraglide/messages.js';

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

<div class="flex justify-between">
	<Button
		variant="ghost"
		class="relative"
		href="/app/{page.params.courseId}/tutor/course/{page.params.itemId}/tests/{page.params
			.testId}/instances/{id}"
	>
		<span>{m.edit()}</span>
	</Button>
	<Button variant="ghost" class="relative" onclick={() => handleActionClick('print')}>
		<span>{m.print()}</span>
	</Button>
	<Button variant="destructive" class="relative" onclick={() => handleActionClick('delete')}>
		<span>{m.delete()}</span>
	</Button>
</div>
