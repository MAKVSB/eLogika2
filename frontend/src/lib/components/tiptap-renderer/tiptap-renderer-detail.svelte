<script lang="ts">
	import ChevronsRightIcon from '@lucide/svelte/icons/chevron-right';
	import ChevronsDownIcon from '@lucide/svelte/icons/chevron-down';
	import * as Collapsible from '$lib/components/ui/collapsible/index.js';
	import type { JSONContent } from '@tiptap/core';
	import TiptapRenderer from './tiptap-renderer.svelte';
	import Button from '$lib/components/ui/button/button.svelte';

	let {
		jsonContent
	}: {
		jsonContent: JSONContent;
	} = $props();

    let open = $state(false)

    let summary = $derived(jsonContent.content?.find((i) => i.type == "detailsSummary"))
    let content = $derived(jsonContent.content?.find((i) => i.type == "detailsContent"))
</script>

<Collapsible.Root class="gap-1 rounded-[0.5rem] border p-2 m-2" bind:open>
    <Collapsible.Trigger>
        <Button class="" variant="outline">
            {#if open}
                <ChevronsDownIcon />
            {:else}
                <ChevronsRightIcon />
            {/if}
        </Button>
        {#each summary?.content ?? [] as innerContent}
            <TiptapRenderer jsonContent={innerContent}></TiptapRenderer>
        {/each}
    </Collapsible.Trigger>
	<Collapsible.Content class="pt-4 ml-4">
        {#each content?.content ?? [] as innerContent}
            <TiptapRenderer jsonContent={innerContent}></TiptapRenderer>
        {/each}
	</Collapsible.Content>
</Collapsible.Root>