<script lang="ts">
	import GlobalState from '$lib/shared.svelte';
	import { goto } from '$app/navigation';

	$effect(() => {
		if (GlobalState.loggedUser !== null && GlobalState.loggedUser !== undefined) {
			console.log("Transfering 14")
			goto('/app');
		}
	});

	let { children } = $props();
</script>

{@render children()}
