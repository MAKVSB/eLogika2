<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import GlobalState from '$lib/shared.svelte';
	import { toast } from 'svelte-sonner';
	import { base } from '$app/paths';

	let { children } = $props();

	let valid = $derived.by(() => {
		const courseId = Number(page.params.courseId);
		if (!isNaN(courseId)) {
			const found = GlobalState.loggedUser?.courses.find((course) => course.id == courseId);
			if (found) {
				return true;
			}
		}
		toast.error('Missing permissions to access course');
		console.log("Transfering 16")
		goto(base+'/app/');
		return false;
	});
</script>

{#if valid}
	{@render children()}
{/if}
