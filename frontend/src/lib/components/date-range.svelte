<script lang="ts">
	import type { StringDate } from '$lib/api_types';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { m } from '$lib/paraglide/messages';
	import { getLocale } from '$lib/paraglide/runtime';

	let {
		start,
		end,
		showTime = true,
		showLabel = false
	}: {
		start: StringDate;
		end: StringDate;
		showTime?: boolean;
		showLabel?: boolean;
	} = $props();
</script>

<Tooltip.Provider>
	<Tooltip.Root>
		<Tooltip.Trigger class="text-left">
			<p>
				{#if showLabel}
					{m.date_range_from()}
				{/if}
				{new Date(start).toLocaleString('cs', {
					dateStyle: 'short',
					timeStyle: showTime ? 'short' : undefined
				})}
			</p>
			<p>
				{#if showLabel}
					{m.date_range_to()}
				{/if}
				{new Date(end).toLocaleString('cs', {
					dateStyle: 'short',
					timeStyle: showTime ? 'short' : undefined
				})}
			</p>
		</Tooltip.Trigger>
		<Tooltip.Content class="grid grid-cols-2 gap-4">
			<p>
				{#if showLabel}
					{m.date_range_from()}
				{/if}
				{new Date(start).toLocaleString(getLocale())}
			</p>
			<p>
				{#if showLabel}
					{m.date_range_to()}
				{/if}
				{new Date(end).toLocaleString(getLocale())}
			</p>
		</Tooltip.Content>
	</Tooltip.Root>
</Tooltip.Provider>
