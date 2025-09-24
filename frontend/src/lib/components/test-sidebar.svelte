<script lang="ts" module>
	import GlobalState from '$lib/shared.svelte';
</script>

<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import NavUser from './nav-user.svelte';
	import type { TestInstanceDTO } from '$lib/api_types';
	import { cn } from '$lib/utils';

	let {
		ref = $bindable(null),
		questionIndex = $bindable(),
		instanceData = $bindable(),
		...restProps
	}: ComponentProps<typeof Sidebar.Root> & {
		questionIndex: number;
		instanceData: TestInstanceDTO;
	} = $props();
</script>

<Sidebar.Root class="top-(--header-height) h-[calc(100svh-var(--header-height))]!" {...restProps}>
	<Sidebar.Content>
		<Sidebar.Group>
			<Sidebar.GroupLabel>Questions</Sidebar.GroupLabel>
			<Sidebar.Menu class="grid grid-cols-4">
				{#each instanceData.questions ?? [] as question}
					<button
						class={cn(
							"flex items-center justify-center border aspect-square",
							questionIndex == question.order - 1 ? "bg-green-100 dark:bg-green-900" : ""
						)}
						onclick={() => {
							questionIndex = question.order - 1;
						}}
					>
						{question.order}
					</button>
				{/each}
			</Sidebar.Menu>
		</Sidebar.Group>
		<div class="flex-grow"></div>
	</Sidebar.Content>
	<Sidebar.Footer>
		{#if GlobalState.loggedUser != null}
			<NavUser bind:user={GlobalState.loggedUser} />
		{/if}
	</Sidebar.Footer>
</Sidebar.Root>
