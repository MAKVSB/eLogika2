<script lang="ts">
	import type { QuestionCheckedByDTO } from '$lib/api_types';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { m } from '$lib/paraglide/messages';
	import { getLocale } from '$lib/paraglide/runtime';

	let { users }: { users: QuestionCheckedByDTO[] } = $props();
	const showItems = 4;
</script>

<Tooltip.Provider>
	<Tooltip.Root>
		<Tooltip.Trigger class="w-full">
			<table>
				<tbody>
					{#each users.slice(0, showItems) as user}
						<tr>
							<td class="pr-2">
								{user.firstName}
								{user.familyName}
							</td>
							<td>
								{new Date(user.checkedAt).toLocaleDateString(getLocale())}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>

			{#if users.length > showItems}
				{m.and_num_more({
					number: users.length - showItems
				})}
			{/if}
		</Tooltip.Trigger>
		<Tooltip.Content class="grid grid-cols-2">
			{#each users as user}
				<p>
					{user.firstName}
					{user.familyName}
				</p>
				({new Date(user.checkedAt).toLocaleString(getLocale())})
			{/each}
		</Tooltip.Content>
	</Tooltip.Root>
</Tooltip.Provider>
