<script lang="ts">
	import type { StringDate } from '$lib/api_types';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { getLocale } from '$lib/paraglide/runtime';

	let {
		dateTime,
		showTime = true
	}: {
		dateTime?: StringDate;
		showTime?: boolean;
	} = $props();
</script>

{#if dateTime}
	<Tooltip.Provider>
		<Tooltip.Root>
			<Tooltip.Trigger class="w-full text-left">
				<p>
					{new Date(dateTime).toLocaleString('cs', {
						dateStyle: 'short',
						timeStyle: showTime ? 'short' : undefined
					})}
				</p>
			</Tooltip.Trigger>
			<Tooltip.Content class="grid grid-cols-2">
				<p>
					{new Date(dateTime).toLocaleString(getLocale())}
				</p>
			</Tooltip.Content>
		</Tooltip.Root>
	</Tooltip.Provider>
{/if}
