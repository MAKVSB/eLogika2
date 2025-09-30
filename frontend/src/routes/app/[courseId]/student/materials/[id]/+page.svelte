<script lang="ts">
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import TiptapRenderer from '$lib/components/tiptap-renderer/tiptap-renderer.svelte';
	import { m } from '$lib/paraglide/messages.js';

	let { data } = $props()
</script>

<div class="m-8">
	{#await data.chapter}
		<Pageloader></Pageloader>
	{:then chapter}
		{#if chapter.data.visible}
			<div class="flex flex-row justify-between">
				<h1 class="mb-8 text-2xl">Study materials
					<b>
						{chapter.data.name}
					</b>
				</h1>
			</div>
			{#if chapter.data.childs}
				<h1 class="text-2xl">Sub-chapters</h1>
				<ul>
					{#each chapter.data.childs as subchapter}
						<li>
							<a href={String(subchapter.id)}>{subchapter.name}</a>
						</li>
					{/each}
				</ul>
			{/if}

			<h1 class="text-2xl">Chapter content</h1>
			<TiptapRenderer jsonContent={chapter.data.content}></TiptapRenderer>
		{:else}
				{m.study_materials_invisible()}
		{/if}
	{/await}
</div>
