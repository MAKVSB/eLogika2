<script lang="ts">
	import '../app.css';
	import GlobalState from '$lib/shared.svelte';
	import { goto } from '$app/navigation';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import { base } from '$app/paths';

	$effect(() => {
		console.log("Transfering 13")
		if (GlobalState.loggedUser !== null) {
			goto(base+'/app');
		} else {
			goto(base+'/login');
		}
	});
</script>

<Pageloader></Pageloader>
