<script lang="ts">
	import { CourseItemTypeEnum, type StudentCourseItemDTO } from '$lib/api_types';
	import * as Table from '$lib/components/ui/table';
	import CourseItemDisplay from './CourseItemDisplay.svelte';
	import ChevronsRightIcon from '@lucide/svelte/icons/chevron-right';
	import ChevronsDownIcon from '@lucide/svelte/icons/chevron-down';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { m } from '$lib/paraglide/messages';

	let {
		item,
		tabCount = 1,
		open = $bindable(false)
	}: {
		item: StudentCourseItemDTO;
		tabCount?: number;
		open?: boolean;
	} = $props();
</script>

<Table.Row
	onclick={() => {
		open = !open;
	}}
	class={open ? 'border-b-0' : ''}
>
	<Table.Cell
		class="align-top w-min"
		rowspan={open && item.type != CourseItemTypeEnum.GROUP ? 1 : 1}
	>
		<div class="flex">
			{#each Array(tabCount - 1) as _, i}
				<div class="w-10"></div>
			{/each}
			{#if item.childs}
				<div class="w-10">
					{#if open}
						<ChevronsDownIcon />
					{:else}
						<ChevronsRightIcon />
					{/if}
				</div>
			{:else}
				<div class="w-10"></div>
			{/if}
			{item.name}
		</div>
	</Table.Cell>
	<Table.Cell>{m.course_item_type_enum({ value: item.type })}</Table.Cell>
	<Table.Cell>
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger class={item.passed ? 'text-green-500' : 'text-red-500'}
					>{item.points > item.pointsMax ? item.pointsMax : item.points}</Tooltip.Trigger
				>
				<Tooltip.Content>
					<p>
						{m.points_min()}: {item.pointsMin}
						<br />
						{m.points_max()}: {item.pointsMax}
						<br />
						{m.points_real()}: {item.points}
					</p>
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	</Table.Cell>
	<Table.Cell>{m.yes_no({ value: String(item.mandatory) })}</Table.Cell>
	<Table.Cell>{m.yes_no({ value: String(item.passed) })}</Table.Cell>
	<Table.Cell>{m.yes_no({ value: String(item.includeInResults) })}</Table.Cell>
	<Table.Cell>{item.maxAttempts}</Table.Cell>
</Table.Row>
{#if open}
	{#if item.type == CourseItemTypeEnum.GROUP}
		{#each item.childs as innerItem}
			<CourseItemDisplay item={innerItem} tabCount={tabCount + 1}></CourseItemDisplay>
		{/each}
	{/if}
{/if}
