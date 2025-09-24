<script lang="ts">
	import '../app.css';
	import GlobalState from '$lib/shared.svelte';
	import { goto } from '$app/navigation';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';

	$effect(() => {
		console.log("Transfering 13")
		if (GlobalState.loggedUser !== null) {
			goto('/app');
		} else {
			goto('/login');
		}
	});
</script>

<Pageloader></Pageloader>
