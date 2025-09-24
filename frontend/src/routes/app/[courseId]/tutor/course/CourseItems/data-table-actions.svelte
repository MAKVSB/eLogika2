<script lang="ts">
	import { page } from '$app/state';
	import { CourseItemTypeEnum } from '$lib/api_types';
	import { Button } from '$lib/components/ui/button/index.js';
	import { m } from '$lib/paraglide/messages.js';

	let {
		id,
		meta,
		type,
		editable
	}: {
		id: number | string;
		meta: any;
		type: CourseItemTypeEnum;
		editable: boolean;
	} = $props();

	function handleActionClick(event: string, params?: any) {
		if ('clickEventHandler' in meta) {
			meta.clickEventHandler(event, id, params);
		}
	}
</script>

<div class="flex justify-between">
	{#if editable}
		<Button variant="ghost" class="relative" href="/app/{page.params.courseId}/tutor/course/{id}">
			<span>{m.edit()}</span>
		</Button>
	{:else}
		<Button variant="ghost" class="relative" href="/app/{page.params.courseId}/tutor/course/{id}">
			<span>{m.view()}</span>
		</Button>
	{/if}

	<!-- {#if editable}
		<Button variant="ghost" class="relative" href="/app/{page.params.courseId}/tutor/course/{id}">
			<span>Výsledky</span>
		</Button>
	{:else}
		<div></div>
	{/if} -->

	{#if type == CourseItemTypeEnum.TEST}
		<Button
			variant="ghost"
			class="relative"
			href="/app/{page.params.courseId}/tutor/course/{id}/tests"
		>
			<span>{m.courseitem_test_generated()}</span>
		</Button>
	{:else if type == CourseItemTypeEnum.ACTIVITY}
		<Button
			variant="ghost"
			class="relative"
			href="/app/{page.params.courseId}/tutor/course/{id}/activities"
		>
			<span>{m.courseitem_activity_submissions()}</span>
		</Button>
	{:else}
		<div></div>
	{/if}
	{#if editable}
		<Button variant="destructive" class="relative" onclick={() => handleActionClick('delete')}>
			<span>{m.delete()}</span>
		</Button>
	{:else}
		<div></div>
	{/if}
</div>
