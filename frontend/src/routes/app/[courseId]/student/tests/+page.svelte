<script lang="ts">
	import { API, ApiError } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';

	import * as Table from '$lib/components/ui/table/index.js';
	import DataTableDateRange from '../../../admin/courses/[courseId]/[itemId]/Terms/data-table-date-range.svelte';
	import {
		TestInstanceStateEnum,
		type ListAvailableTestsResponse,
		type StudentTestDTO,
		type TestInstancePrepareRequest,
		type TestInstancePrepareResponse
	} from '$lib/api_types';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages';
	import { invalidate } from '$app/navigation';

	let courseId = $derived<string>(page.params.courseId);
	let isLoading = $state(true);

	let { data } = $props();
	let lastData: ListAvailableTestsResponse = $state({
		instances: [],
		items: []
	});

	$effect(() => {
		data.tests
			.then((res) => {
				lastData = res;
				isLoading = false;
			})
			.catch(() => {});
	});

	const startInstance = async (item: StudentTestDTO) => {
		await API.request<TestInstancePrepareRequest, TestInstancePrepareResponse>(
			`/api/v2/courses/${courseId}/tests/prepare`,
			{
				method: 'POST',
				body: {
					termId: item.termId,
					courseItemId: item.courseItemId
				}
			}
		)
			.then((res) => {
				invalidate((url) => {
					return url.href.endsWith('/tests/available');
				});
				openInstance(res.instanceId);
			})
			.catch(() => {});
	};

	const openInstance = async (instanceId: number) => {
		window.open(
			`/test/${instanceId}`,
			'_blank',
			'toolbar=no,menubar=no,scrollbars=no,resizable=yes,width=screen.width,height=screen.height'
		);
	};
</script>

{#snippet TableHeader()}
	<Table.Row>
		<Table.Head>Name</Table.Head>
		<Table.Head>Activity</Table.Head>
		<Table.Head>Attempts left</Table.Head>
		<Table.Head>Active (from/to)</Table.Head>
		<Table.Head>{m.actions()}</Table.Head>
	</Table.Row>
{/snippet}
{#snippet TableHeaderRunning()}
	<Table.Row>
		<Table.Head>Name</Table.Head>
		<Table.Head>Activity</Table.Head>
		<Table.Head>State</Table.Head>

		<Table.Head>Start until</Table.Head>
		<Table.Head>Started at</Table.Head>
		<Table.Head>Ends at</Table.Head>
		<Table.Head>{m.actions()}</Table.Head>
	</Table.Row>
{/snippet}
<div class="flex flex-col gap-8 m-8">
	{#if isLoading}
		<Pageloader></Pageloader>
	{:else}
		<div class="flex flex-row justify-between">
			<h1 class="text-2xl">Testy k vypracování:</h1>
		</div>

		<div class="flex flex-col gap-4">
			<h2 class="text-xl">Running tests</h2>
			<Table.Root>
				<Table.Header class="font-medium border-t bg-muted/50">
					{@render TableHeaderRunning()}
				</Table.Header>
				<Table.Body>
					{@const filtered = lastData.instances}
					{#if filtered.length === 0}
						<Table.Row>
							<Table.Cell colspan={7}>No items found</Table.Cell>
						</Table.Row>
					{:else}
						{#each filtered as item}
							<Table.Row class={item.state == TestInstanceStateEnum.ACTIVE ? 'bg-green-100' : ';'}>
								<Table.Cell>{item.termName}</Table.Cell>
								<Table.Cell>{item.courseItemName}</Table.Cell>
								<Table.Cell>{m.test_instance_state_enum({ value: item.state })}</Table.Cell>
								<Table.Cell>
									{new Date(item.startUntil).toLocaleString('cs', {
										dateStyle: 'short',
										timeStyle: 'short'
									})}
								</Table.Cell>
								<Table.Cell>
									{#if new Date(item.startedAt).getTime() > 0}
										{new Date(item.startedAt).toLocaleString('cs', {
											dateStyle: 'short',
											timeStyle: 'short'
										})}
									{/if}
								</Table.Cell>
								<Table.Cell>
									{#if new Date(item.endsAt).getTime() > 0}
										{new Date(item.endsAt).toLocaleString('cs', {
											dateStyle: 'short',
											timeStyle: 'short'
										})}
									{/if}
								</Table.Cell>
								<Table.Cell>
									<Button onclick={() => openInstance(item.id)}>Continue test</Button>
								</Table.Cell>
							</Table.Row>
						{/each}
					{/if}
				</Table.Body>
				{#if lastData.items.length > 10}
					<Table.Footer>
						{@render TableHeaderRunning()}
					</Table.Footer>
				{/if}
			</Table.Root>
		</div>

		<div class="flex flex-col gap-4">
			<h2 class="text-xl">Tests to start</h2>
			<Table.Root>
				<Table.Header class="font-medium border-t bg-muted/50">
					{@render TableHeader()}
				</Table.Header>
				<Table.Body>
					{@const filtered = lastData.items}
					{#if filtered.length === 0}
						<Table.Row>
							<Table.Cell colspan={7}>No items found</Table.Cell>
						</Table.Row>
					{:else}
						{#each filtered as item}
							<Table.Row class={item.canStart ? '' : 'opacity-50'}>
								<Table.Cell>{item.termName}</Table.Cell>
								<Table.Cell>{item.courseItemName}</Table.Cell>
								<Table.Cell>{item.triesLeft}</Table.Cell>
								<Table.Cell>
									<DataTableDateRange start={item.activeFrom} end={item.activeTo} showTime={true}
									></DataTableDateRange>
								</Table.Cell>
								<Table.Cell>
									<Button disabled={!item.canStart} onclick={() => startInstance(item)}
										>Start test</Button
									>
								</Table.Cell>
							</Table.Row>
						{/each}
					{/if}
				</Table.Body>
				{#if lastData.items.length > 10}
					<Table.Footer>
						{@render TableHeader()}
					</Table.Footer>
				{/if}
			</Table.Root>
		</div>
	{/if}
</div>
