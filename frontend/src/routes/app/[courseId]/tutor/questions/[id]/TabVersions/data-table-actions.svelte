<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button/index.js';
	import { m } from '$lib/paraglide/messages';
	import { base } from '$app/paths';

	let {
		id,
		isArchive,
		meta
	}: {
		id: number | string;
		isArchive: boolean;
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
		href="{base}/app/{page.params.courseId}/tutor/questions/{id}"
	>
		<span>{m.edit()}</span>
	</Button>

	<Button
		variant={page.params.id == id ? 'default' : 'outline'}
		class="relative"
		onclick={() => handleActionClick('selectversion')}
		disabled={!isArchive}
	>
		<span
			>{m.question_version_select()}
			{#if page.params.id == id}
				({m.current()})
			{/if}
		</span>
	</Button>
</div>
