<script lang="ts">
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';

	let {
		checked,
		meta,
		id
	}: {
		checked: boolean;
		meta?: any;
		id?: number;
	} = $props();

	function onclick() {
		if (meta && 'clickEventHandler' in meta && id) {
			meta.clickEventHandler(id);
		}
	}
</script>

<div class="flex items-center">
	<div class="flex flex-col">
		<Checkbox
			{checked}
			onclick={() => onclick()}
			disabled={!(meta && 'clickEventHandler' in meta && id)}
		/>
	</div>
</div>
