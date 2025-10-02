<script lang="ts">
	import type { QuestionCreatedByDTO } from '$lib/api_types';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { getLocale } from '$lib/paraglide/runtime';

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
			{createdBy.degreeBefore}
			{createdBy.firstName}
			{createdBy.familyName}
			{createdBy.degreeAfter}
		</Tooltip.Trigger>
		<Tooltip.Content class="grid grid-cols-2">
			<p>
				{createdBy.degreeBefore}
				{createdBy.firstName}
				{createdBy.familyName}
				{createdBy.degreeAfter}
			</p>
			({new Date(createdAt).toLocaleString(getLocale())})
		</Tooltip.Content>
	</Tooltip.Root>
</Tooltip.Provider>
