<script lang="ts">
	import {
		type TestInstanceDTO
	} from '$lib/api_types.js';
	import * as Tabs from '$lib/components/ui/tabs';
	import InstanceEditor from './InstanceEditor.svelte';
	import InstanceEvents from './InstanceEvents.svelte';

	let { data } = $props();

	let loading = $state(true);
	let instanceData: TestInstanceDTO | undefined = $state();

	$effect(() => {
		data.test
			.then((res) => {
				instanceData = res.instanceData;
				instanceData.points = instanceData.points- instanceData.bonusPoints; 
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
		<Tabs.Root value="instanceEditor">
			<Tabs.List>
				<Tabs.Trigger value="instanceEditor">Test editor</Tabs.Trigger>
				<Tabs.Trigger value="events">Events</Tabs.Trigger>
			</Tabs.List>
			<Tabs.Content value="instanceEditor" class="p-4">
				<InstanceEditor {instanceData} editable={true}></InstanceEditor>
			</Tabs.Content>
			<Tabs.Content value="events" class="p-4">
				<InstanceEvents></InstanceEvents>
			</Tabs.Content>
			<Tabs.List direction="up">
				<Tabs.Trigger value="instanceEditor" direction="up">Test editor</Tabs.Trigger>
				<Tabs.Trigger value="events" direction="up">Events</Tabs.Trigger>
			</Tabs.List>
		</Tabs.Root>
	{/if}
</div>
