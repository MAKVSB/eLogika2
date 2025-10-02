<script lang="ts">
	import type { CourseItemResultsDTO } from '$lib/api_types';
	import * as Table from '$lib/components/ui/table/index.js';
	import { m } from '$lib/paraglide/messages';
	import * as Collapsible from '$lib/components/ui/collapsible/index.js';

	import ChevronsRightIcon from '@lucide/svelte/icons/chevron-right';
	import ChevronsDownIcon from '@lucide/svelte/icons/chevron-down';
	import { Button } from '$lib/components/ui/button';
	import ResultsTable from './ResultsTable.svelte';

	let {
		userData
	}: {
		userData: CourseItemResultsDTO;
	} = $props();

	let open = $state(false);
</script>

<Table.Row>
	<Table.Cell>
		<Collapsible.Root bind:open>
			<Collapsible.Trigger>
				<Button class="" variant="outline">
					{#if open}
						<ChevronsDownIcon />
					{:else}
						<ChevronsRightIcon />
					{/if}
				</Button>
			</Collapsible.Trigger>
		</Collapsible.Root>
	</Table.Cell>
	<Table.Cell>{userData.username}</Table.Cell>
	<Table.Cell>
		{userData.degreeBefore}
		{userData.firstName}
		{userData.familyName}
		{userData.degreeAfter}
	</Table.Cell>
	<Table.Cell></Table.Cell>
	<Table.Cell>
		{userData.points}
	</Table.Cell>
	<Table.Cell>{m.yes_no({ value: String(userData.passed) })}</Table.Cell>
</Table.Row>
{#if open}
	<Table.Row>
		<Table.Cell></Table.Cell>
		<Table.Cell colspan={6}>
			<ResultsTable studentId={userData.id} results={userData.results}></ResultsTable>
		</Table.Cell>
	</Table.Row>
{/if}
