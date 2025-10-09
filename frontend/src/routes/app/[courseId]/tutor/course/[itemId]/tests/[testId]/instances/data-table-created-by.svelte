<script lang="ts">
	import type { QuestionCreatedByDTO } from '$lib/api_types';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { getLocale } from '$lib/paraglide/runtime';
	import { displayUserName } from '$lib/utils';

	let {
		createdBy,
		createdAt
	}: {
		createdBy: QuestionCreatedByDTO;
		createdAt: string;
	} = $props();
</script>

<Tooltip.Provider>
	<Tooltip.Root>
		<Tooltip.Trigger class="w-full">
			{displayUserName(createdBy)}
		</Tooltip.Trigger>
		<Tooltip.Content class="grid grid-cols-2">
			<p>
				{displayUserName(createdBy)}
			</p>
			({new Date(createdAt).toLocaleString(getLocale())})
		</Tooltip.Content>
	</Tooltip.Root>
</Tooltip.Provider>
