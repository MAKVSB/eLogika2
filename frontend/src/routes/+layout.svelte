<script lang="ts">
	import '../app.css';

	import GlobalState from '$lib/shared.svelte';
	import { goto } from '$app/navigation';
	import { Toaster } from '$lib/components/ui/sonner';
	import { ModeWatcher } from 'mode-watcher';
	import { page } from '$app/state';
	import { base } from '$app/paths';

	let { children } = $props();

	$effect(() => {
		if (GlobalState.loggedUser === null && !page.url.pathname.startsWith(base+"/login")) {
			console.log("Transfering 12")
			goto(base+'/login');
		}
	});
	// TODO: Opravdu tu nemůže být loading image ?
</script>

<!-- {#if GlobalState.loggedUser == null}
	<PageLoader></PageLoader>
{:else} -->
<ModeWatcher />
<div class="[--header-height:calc(--spacing(14))]">
	{@render children()}
</div>
<Toaster richColors />
<!-- {/if} -->
