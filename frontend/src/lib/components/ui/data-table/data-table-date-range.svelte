<script lang="ts">
	import type { StringDate } from '$lib/api_types';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
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
		<Tooltip.Trigger class="w-full text-left">
			<p> {showLabel ? "From" : ""}
				{new Date(start).toLocaleString('cs', {
					dateStyle: 'short',
					timeStyle: showTime ? 'short' : undefined
				})}
			</p>
			<p> {showLabel ? "To" : ""}
				{new Date(end).toLocaleString('cs', {
					dateStyle: 'short',
					timeStyle: showTime ? 'short' : undefined
				})}
			</p>
		</Tooltip.Trigger>
		<Tooltip.Content class="grid grid-cols-2">
			<p>From</p>
			<p>
				{new Date(start).toLocaleString(getLocale())}
			</p>
			<p>To</p>
			<p> 
				{new Date(end).toLocaleString(getLocale())}
			</p>
		</Tooltip.Content>
	</Tooltip.Root>
</Tooltip.Provider>
