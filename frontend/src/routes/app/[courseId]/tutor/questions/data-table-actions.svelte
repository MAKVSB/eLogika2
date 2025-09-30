<script lang="ts">
	import { page } from '$app/state';
	import type { QuestionCheckedByDTO } from '$lib/api_types';
	import { Button } from '$lib/components/ui/button/index.js';
	import GlobalState from '$lib/shared.svelte';
	import { m } from '$lib/paraglide/messages';
	import { base } from '$app/paths';

	let {
		id,
		checkedBy,
		meta
	}: {
		id: number | string;
		checkedBy: QuestionCheckedByDTO[];
		meta: any;
	} = $props();

	function handleActionClick(event: string, params?: any) {
		if ('clickEventHandler' in meta) {
			meta.clickEventHandler(event, id, params);
		}
	}
</script>

<div class="flex justify-between gap-2">
	<Button variant="outline" class="relative" href="{base}/app/{page.params.courseId}/tutor/questions/{id}">
		<span>{m.edit()}</span>
	</Button>

	{#if checkedBy.find((usr) => usr.id == GlobalState.loggedUser?.id)}
		<Button variant="outline" onclick={() => handleActionClick('uncheck')}
			>{m.question_check_action_uncheck_list()}</Button
		>
	{:else}
		<Button variant="outline" onclick={() => handleActionClick('check')}>
			{m.question_check_action_check_list()}
		</Button>
	{/if}

	<Button variant="outline" class="relative" onclick={() => handleActionClick('print')}>
		<span>{m.print()}</span>
	</Button>

	<Button variant="destructive" class="relative" onclick={() => handleActionClick('delete')}>
		<span>{m.delete()}</span>
	</Button>
</div>
