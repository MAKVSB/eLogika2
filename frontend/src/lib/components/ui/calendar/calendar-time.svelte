<script lang="ts">
	import { TimeField } from 'bits-ui';
	import type { WithoutChildrenOrChild } from '$lib/utils';

	let {
		ref = $bindable(),
		value = $bindable(),
		placeholder = $bindable(),
		locale = 'en',
		hideTimeZone = true,
		disabled = false,
		readonly = false,
		required = false,
		onValueChange
	}: WithoutChildrenOrChild<TimeField.RootProps> = $props();
</script>

<TimeField.Root
	class="group flex w-full max-w-[320px] flex-col gap-1.5"
	{locale}
	{hideTimeZone}
	bind:value
	{disabled}
	{readonly}
	{required}
	{onValueChange}
>
	<div
		class="h-input rounded-input border-border-input text-foreground focus-within:border-border-input-hover focus-within:shadow-date-field-focus hover:border-border-input-hover group-data-invalid:border-destructive mt-3 flex w-full items-center justify-around border-t px-2 py-3 text-sm tracking-[0.01em] select-none"
	>
		<TimeField.Input>
			{#snippet children({ segments })}
				{#each segments as { part, value }, i (part + i)}
					{#if part != 'literal' || value != ' '}
						<div class="inline-block select-none">
							<TimeField.Segment
								{part}
								class="rounded-5px hover:bg-muted focus:bg-muted focus:text-foreground aria-[valuetext=Empty]:text-muted-foreground px-1 py-1 focus-visible:ring-0! focus-visible:ring-offset-0!"
							>
								{value}
							</TimeField.Segment>
						</div>
					{/if}
				{/each}
			{/snippet}
		</TimeField.Input>
	</div>
</TimeField.Root>
