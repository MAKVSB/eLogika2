<script lang="ts">
	import type { TestInstanceDTO } from '$lib/api_types.js';
	import InstanceEditor from './InstanceEditor.svelte';

	let { data } = $props();

	let loading = $state(true);
	let instanceData: TestInstanceDTO | undefined = $state();

	$effect(() => {
		data.test
			.then((res) => {
				instanceData = res.instanceData;
				instanceData.points = instanceData.points - instanceData.bonusPoints;
				loading = false;
			})
			.catch(() => {});
	});
</script>

<div class="m-8">
	<div class="flex flex-row justify-between">
		<h1 class="mb-8 text-2xl">Test instance edit</h1>
	</div>
	{#if !loading && instanceData}
		<InstanceEditor {instanceData} editable={false}></InstanceEditor>
	{/if}
</div>
