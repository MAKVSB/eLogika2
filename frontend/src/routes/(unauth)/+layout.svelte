<script lang="ts">
	import GlobalState from '$lib/shared.svelte';
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';

	$effect(() => {
		if (GlobalState.loggedUser !== null && GlobalState.loggedUser !== undefined) {
			console.log("Transfering 14")
			goto(base+'/app');
		}
	});

	let { children } = $props();
</script>

{@render children()}
