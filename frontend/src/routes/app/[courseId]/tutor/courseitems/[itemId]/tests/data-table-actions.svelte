<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button/index.js';
	import { m } from '$lib/paraglide/messages.js';
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
	<Button
		variant="outline"
		class="relative"
		href="{base}/app/{page.params.courseId}/tutor/courseitems/{page.params
			.itemId}/tests/{id}/instances"
	>
		<span>{m.test_participants()}</span>
	</Button>
	<Button variant="outline" class="relative" onclick={() => handleActionClick('reevaluate')}>
		<span>{m.test_reevaluate({ type: 'single' })}</span>
	</Button>
	<Button variant="outline" class="relative" onclick={() => handleActionClick('print')}>
		<span>{m.print()}</span>
	</Button>
	<Button variant="destructive" class="relative" onclick={() => handleActionClick('delete')}>
		<span>{m.delete()}</span>
	</Button>
</div>
