<script lang="ts">
	import type { QuestionCreatedByDTO, StringDate } from '$lib/api_types';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { getLocale } from '$lib/paraglide/runtime';
	import { displayUserName } from '$lib/utils';

	let {
		createdBy,
		createdAt
	}: {
		createdBy: QuestionCreatedByDTO;
		createdAt: StringDate;
	} = $props();
</script>

<Tooltip.Provider>
	<Tooltip.Root>
		<Tooltip.Trigger>
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
