<script lang="ts">
	import { goto } from '$app/navigation';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import GlobalState from '$lib/shared.svelte';
	import { base } from '$app/paths';

	let loaded = $state(false);

	$effect(() => {
		if (GlobalState.loggedUser?.courses[0]) {
			console.log("Transfering 15")
			goto(base+"/app/" + GlobalState.loggedUser?.courses[0].id)
		} else {
			loaded = true;
		}
	});
</script>

{#if !loaded}
	<Pageloader></Pageloader>
{:else}
	<div class="flex flex-col gap-4 m-8">
		<h1 class="text-2xl">
			There is no course available for you now.
		</h1>
		<h2 class="text-xl">
            If you think this is a mistake, please contact tutor or garant of the expected course
		</h2>
	</div>
{/if}
