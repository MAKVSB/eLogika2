<script lang="ts">
	import * as Popover from '$lib/components/ui/popover';
	import { DateRangeField } from 'bits-ui';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { RangeCalendar } from '$lib/components/ui/range-calendar';
	import { cn, type WithoutChildrenOrChild } from '$lib/utils';
	import CalendarIcon from '@lucide/svelte/icons/calendar-days';

	let {
		value = $bindable(),
		locale = 'en',
		granularity = 'day',
		class: className,
		hideTimeZone = true,
		onValueChange: onValueChange2,
		disabled
	}: WithoutChildrenOrChild<DateRangeField.RootProps> = $props();

	const onValueChange = (e: any) => {
		console.log("onValueChange")
		if (onValueChange2) {
			onValueChange2(e)
		}
	}
</script>

<DateRangeField.Root
	class="flex justify-between w-full gap-1 group"
	bind:value
	{granularity}
	{locale}
	{hideTimeZone}
	{onValueChange}
	{disabled}
	hourCycle={24}
>
	<div
		class={cn(
			'selection:bg-primary dark:bg-input/30 selection:text-primary-foreground border-input ring-offset-background placeholder:text-muted-foreground flex sm:h-9 w-full min-w-0 rounded-md border bg-transparent px-3 py-2 text-sm font-medium shadow-xs transition-[color,box-shadow] outline-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm',
			'focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]',
			'aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive',
			'items-center flex-col sm:flex-row',
			disabled ? "cursor-not-allowed opacity-50" : "",
			className,
		)}
	>
		{#each ['start', 'end'] as const as type (type)}
			<DateRangeField.Input {type}>
				{#snippet children({ segments })}
					{#each segments as { part, value }, i (part + i)}
						<div class="inline-block select-none">
							{#if part === 'literal'}
								<DateRangeField.Segment {part} class="text-muted-foreground">
									{value}
								</DateRangeField.Segment>
							{:else}
								<DateRangeField.Segment
									{part}
									class="rounded-5px hover:bg-muted focus:bg-muted focus:text-foreground aria-[valuetext=Empty]:text-muted-foreground px-1 py-1 focus-visible:ring-0! focus-visible:ring-offset-0!"
								>
									{value}
								</DateRangeField.Segment>
							{/if}
						</div>
					{/each}
				{/snippet}
			</DateRangeField.Input>
			{#if type === 'start'}
				<div aria-hidden="true" class="px-1 text-muted-foreground">–⁠⁠⁠⁠⁠</div>
			{/if}
		{/each}
		<div class="flex-1"></div>
		<Popover.Root>
			<Popover.Trigger class={buttonVariants({ variant: 'ghost', size: 'sm' })}
				><CalendarIcon /></Popover.Trigger
			>
			<Popover.Content class="w-64 p-0 m-0">
				<RangeCalendar
					bind:value
					{locale}
					{granularity}
					hideTimeZone
					captionLayout="dropdown"
					{disabled}
					{onValueChange}
				/>
			</Popover.Content>
		</Popover.Root>
	</div>
</DateRangeField.Root>
