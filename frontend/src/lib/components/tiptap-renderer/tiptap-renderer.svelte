<script lang="ts">
	import type { JSONContent } from '@tiptap/core';
	import TiptapRenderer from './tiptap-renderer.svelte';
	import { cn } from '$lib/utils';
	import TiptapRendererDetail from './tiptap-renderer-detail.svelte';
	import katex from 'katex';
	import CodeBlock from '../ui/code-block/code-block.svelte';

	let {
		jsonContent
	}: {
		jsonContent: JSONContent;
	} = $props();

	function isVisiblyEmpty(str?: string) {
		if (!str) return true;

		const cleaned = str.replace(/[\s\u200B-\u200D\uFEFF]/g, '');
		return cleaned.length === 0;
	}
</script>

{#if jsonContent.type == 'text'}
	{#if jsonContent.marks && jsonContent.marks.length > 0}
		{@const mark = jsonContent.marks[0]}
		{@const rest = {
			...jsonContent,
			marks: jsonContent.marks.slice(1)
		}}
		{#if mark.type == 'highlight'}
			<span class="bg-amber-200">
				<TiptapRenderer jsonContent={rest}></TiptapRenderer>
			</span>
		{:else if mark.type == 'subscript'}
			<sub>
				<TiptapRenderer jsonContent={rest}></TiptapRenderer>
			</sub>
		{:else if mark.type == 'superscript'}
			<sup>
				<TiptapRenderer jsonContent={rest}></TiptapRenderer>
			</sup>
		{:else if mark.type == 'strike'}
			<span class="line-through">
				<TiptapRenderer jsonContent={rest}></TiptapRenderer>
			</span>
		{:else if mark.type == 'bold'}
			<span class="font-bold">
				<TiptapRenderer jsonContent={rest}></TiptapRenderer>
			</span>
		{:else if mark.type == 'underline'}
			<span class="underline">
				<TiptapRenderer jsonContent={rest}></TiptapRenderer>
			</span>
		{:else if mark.type == 'italic'}
			<span class="italic">
				<TiptapRenderer jsonContent={rest}></TiptapRenderer>
			</span>
		{:else if mark.type == 'link'}
			<a
				href={mark.attrs?.href}
				target={mark.attrs?.target}
				rel={mark.attrs?.rel}
				class="underline"
			>
				<TiptapRenderer jsonContent={rest}></TiptapRenderer>
			</a>
		{:else}
			<span class="text-red-500">{JSON.stringify(mark)}</span>
		{/if}
	{:else if isVisiblyEmpty(jsonContent.text)}
		<br />
	{:else}
		{jsonContent.text}
	{/if}
{:else if jsonContent.type == 'heading'}
	<h1
		class={cn(
			jsonContent.attrs?.level == 1 ? 'text-2xl' : '',
			jsonContent.attrs?.level == 2 ? 'text-xl' : '',
			jsonContent.attrs?.level == 3 ? 'text-lg' : '',
			jsonContent.attrs?.textAlign == 'right' ? 'text-right' : '',
			jsonContent.attrs?.textAlign == 'left' ? 'text-left' : '',
			jsonContent.attrs?.textAlign == 'center' ? 'text-center' : '',
			jsonContent.attrs?.textAlign == 'justify' ? 'text-justify' : ''
		)}
	>
		{#each jsonContent.content ?? [] as innerContent}
			<TiptapRenderer jsonContent={innerContent}></TiptapRenderer>
		{/each}
	</h1>
{:else if jsonContent.type == 'paragraph'}
	<p
		class={cn(
			jsonContent.attrs?.textAlign == 'right' ? 'text-right' : '',
			jsonContent.attrs?.textAlign == 'left' ? 'text-left' : '',
			jsonContent.attrs?.textAlign == 'center' ? 'text-center' : '',
			jsonContent.attrs?.textAlign == 'justify' ? 'text-justify' : ''
		)}
	>
		{#each jsonContent.content ?? [] as innerContent}
			<TiptapRenderer jsonContent={innerContent}></TiptapRenderer>
		{:else}
			<br />
		{/each}
	</p>
{:else if jsonContent.type == 'custom-image'}
	<img
		src={jsonContent.attrs?.mode === 'storage'
			? import.meta.env.VITE_API_URL + '/api/v2/files/' + jsonContent.attrs?.src
			: jsonContent.attrs?.src}
		alt={jsonContent.attrs?.alt}
		style="width: {jsonContent.attrs?.width}px"
		height="auto"
		class="my-4"
	/>
{:else if jsonContent.type == 'table'}
	<div class="w-full overflow-scroll">
		<table class="border">
			{#each jsonContent.content ?? [] as innerContent}
				<TiptapRenderer jsonContent={innerContent}></TiptapRenderer>
			{/each}
		</table>
	</div>
{:else if jsonContent.type == 'tableRow'}
	<tr class="border">
		{#each jsonContent.content ?? [] as innerContent}
			<TiptapRenderer jsonContent={innerContent}></TiptapRenderer>
		{/each}
	</tr>
{:else if jsonContent.type == 'tableHeader'}
	<th
		class="p-2 border"
		colspan={jsonContent.attrs?.colspan}
		rowspan={jsonContent.attrs?.rowspan}
		style="width:{jsonContent.attrs?.colwidth}px"
	>
		{#each jsonContent.content ?? [] as innerContent}
			<TiptapRenderer jsonContent={innerContent}></TiptapRenderer>
		{/each}
	</th>
{:else if jsonContent.type == 'tableCell'}
	<td
		class="p-2 border"
		colspan={jsonContent.attrs?.colspan}
		rowspan={jsonContent.attrs?.rowspan}
		style="width:{jsonContent.attrs?.colwidth}px"
	>
		{#each jsonContent.content ?? [] as innerContent}
			<TiptapRenderer jsonContent={innerContent}></TiptapRenderer>
		{/each}
	</td>
{:else if jsonContent.type == 'orderedList'}
	<ol start={jsonContent.attrs?.start ?? 1}>
		{#each jsonContent.content ?? [] as innerContent}
			<TiptapRenderer jsonContent={innerContent}></TiptapRenderer>
		{/each}
	</ol>
{:else if jsonContent.type == 'bulletList'}
	<ul>
		{#each jsonContent.content ?? [] as innerContent}
			<TiptapRenderer jsonContent={innerContent}></TiptapRenderer>
		{/each}
	</ul>
{:else if jsonContent.type == 'listItem'}
	<li>
		{#each jsonContent.content ?? [] as innerContent}
			<TiptapRenderer jsonContent={innerContent}></TiptapRenderer>
		{/each}
	</li>
{:else if jsonContent.type == 'inlineMath'}
	{#if jsonContent.attrs && jsonContent.attrs.latex}
		{@html katex.renderToString(jsonContent.attrs?.latex)}
	{/if}
{:else if jsonContent.type == 'blockMath'}
	<div class="w-full text-center">
		{#if jsonContent.attrs && jsonContent.attrs.latex}
			{@html katex.renderToString(jsonContent.attrs?.latex)}
		{/if}
	</div>
{:else if jsonContent.type == 'details'}
	<TiptapRendererDetail {jsonContent}></TiptapRendererDetail>
{:else if jsonContent.type == 'codeBlock'}
	{#if jsonContent.attrs?.filename}
		<img
			src={import.meta.env.VITE_API_URL + '/api/v2/files/' + jsonContent.attrs?.filename}
			alt={jsonContent.attrs?.alt}
			style="width: {jsonContent.attrs?.width}px"
			height="auto"
			class="my-4 bg-white"
		/>
	{:else if jsonContent.content?.length == 1}
		<CodeBlock code={jsonContent.content[0].text ?? ''}></CodeBlock>
	{/if}
{:else if jsonContent.type == 'doc'}
	{#each jsonContent.content ?? [] as innerContent}
		<TiptapRenderer jsonContent={innerContent}></TiptapRenderer>
	{/each}
{:else}
	<p class="text-red-500">{JSON.stringify(jsonContent)}</p>
{/if}
