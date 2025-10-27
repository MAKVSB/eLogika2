<script lang="ts">
	import { Tabs as TabsPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		direction = "down",
		...restProps
	}: TabsPrimitive.ListProps & {
		direction?: "up" | "down"
	}= $props();
</script>

<TabsPrimitive.List
	bind:ref
	data-slot="tabs-list"
	class={cn(
		'w-full',
		'text-muted-foreground flex h-9 items-center justify-center',
		direction == "down" ? "border-b" : "border-t",
		className
	)}
	{...restProps}
/>
