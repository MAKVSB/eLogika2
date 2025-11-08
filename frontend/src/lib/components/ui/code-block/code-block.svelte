<script lang="ts">
	import hljs from 'highlight.js';
	import { onMount } from 'svelte';
	import { mode } from 'mode-watcher';

	let {
		code,
		languageNames = ['latex']
	}: {
		code: string;
		languageNames?: string[];
	} = $props();

	let highlighted = $state('');
	let lang = $state('');

	onMount(() => {
		const { value, language } = hljs.highlightAuto(code, languageNames);
		highlighted = value;
		if (language) {
			lang = language;
		} else {
			lang = '';
		}
	});
</script>

<svelte:head>
	{#if mode.current == 'dark'}
		<style>
			/* Comment */
			.hljs-comment,
			.hljs-quote {
				color: #d4d0ab;
			}
			/* Red */
			.hljs-variable,
			.hljs-template-variable,
			.hljs-tag,
			.hljs-name,
			.hljs-selector-id,
			.hljs-selector-class,
			.hljs-regexp,
			.hljs-deletion {
				color: #ffa07a;
			}
			/* Orange */
			.hljs-number,
			.hljs-built_in,
			.hljs-literal,
			.hljs-type,
			.hljs-params,
			.hljs-meta,
			.hljs-link {
				color: #f5ab35;
			}
			/* Yellow */
			.hljs-attribute {
				color: #ffd700;
			}
			/* Green */
			.hljs-string,
			.hljs-symbol,
			.hljs-bullet,
			.hljs-addition {
				color: #abe338;
			}
			/* Blue */
			.hljs-title,
			.hljs-section {
				color: #00e0e0;
			}
			/* Purple */
			.hljs-keyword,
			.hljs-selector-tag {
				color: #dcc6e0;
			}
		</style>
	{:else}
		<style>
			/* Comment */
			.hljs-comment,
			.hljs-quote {
				color: #696969;
			}
			/* Red */
			.hljs-variable,
			.hljs-template-variable,
			.hljs-tag,
			.hljs-name,
			.hljs-selector-id,
			.hljs-selector-class,
			.hljs-regexp,
			.hljs-deletion {
				color: #d91e18;
			}
			/* Orange */
			.hljs-number,
			.hljs-built_in,
			.hljs-literal,
			.hljs-type,
			.hljs-params,
			.hljs-meta,
			.hljs-link {
				color: #aa5d00;
			}
			/* Yellow */
			.hljs-attribute {
				color: #aa5d00;
			}
			/* Green */
			.hljs-string,
			.hljs-symbol,
			.hljs-bullet,
			.hljs-addition {
				color: #008000;
			}
			/* Blue */
			.hljs-title,
			.hljs-section {
				color: #007faa;
			}
			/* Purple */
			.hljs-keyword,
			.hljs-selector-tag {
				color: #7928a1;
			}
		</style>
	{/if}
	<style>
		pre code.hljs {
			display: block;
			overflow-x: auto;
			padding: 1em;
		}
		.hljs-emphasis {
			font-style: italic;
		}
		.hljs-strong {
			font-weight: bold;
		}
		@media screen and (-ms-high-contrast: active) {
			.hljs-addition,
			.hljs-attribute,
			.hljs-built_in,
			.hljs-bullet,
			.hljs-comment,
			.hljs-link,
			.hljs-literal,
			.hljs-meta,
			.hljs-number,
			.hljs-params,
			.hljs-string,
			.hljs-symbol,
			.hljs-type,
			.hljs-quote {
				color: highlight;
			}
			.hljs-keyword,
			.hljs-selector-tag {
				font-weight: bold;
			}
		}
	</style>
</svelte:head>

<div
	class="flex flex-col w-full border not-prose border-border bg-card text-card-foreground overflow-clip rounded-xl"
>
	<div class="w-full overflow-x-auto text-[13px] [&>pre]:px-4 [&>pre]:py-4">
		<pre><code
				>{#if highlighted}{@html highlighted}{:else}{code}{/if}</code
			></pre>
	</div>
</div>
