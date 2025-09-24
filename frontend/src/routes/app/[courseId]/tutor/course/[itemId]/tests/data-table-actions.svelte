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
		href="/app/{page.params.courseId}/tutor/course/{page.params.itemId}/tests/{id}/instances"
	>
		<span>{m.test_participants()}</span>
	</Button>
	<Button
		variant="ghost"
		class="relative"
		onclick={() =>
			handleActionClick('print', {
				instances: false
			})}
	>
		<span>Print test</span>
	</Button>
	<!-- <Button
		variant="ghost"
		class="relative"
		onclick={() => handleActionClick('print', { instances: true })}
	>
		<span>Print instances</span>
	</Button> -->
	<Button variant="destructive" class="relative" onclick={() => handleActionClick('delete')}>
		<span>{m.delete()}</span>
	</Button>
</div>
