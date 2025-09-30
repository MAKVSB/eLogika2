<script lang="ts">
	import { Toaster as Sonner, type ToasterProps as SonnerProps } from 'svelte-sonner';
	import { mode } from 'mode-watcher';

	let { ...restProps }: SonnerProps = $props();
</script>

<div>
	<Sonner
		closeButton
		duration={7000}
		richColors
		position="top-center"
		expand
		theme={mode.current}
		class="toaster group"
		style="--normal-bg: var(--popover); --normal-text: var(--popover-foreground); --normal-border: var(--border);"
		{...restProps}
	/>
</div>

