<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { API, ApiError } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';

	import * as Table from '$lib/components/ui/table/index.js';
	import DateRange from '$lib/components/date-range.svelte';
	import type { StudentTermDTO, StudentTermsListResponse, TermsJoinResponse } from '$lib/api_types';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages';
	import { invalidate } from '$app/navigation';
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import Label from '$lib/components/ui/label/label.svelte';

	let courseId = $derived<string>(page.params.courseId);
	let isLoading = $state(true);

	let { data } = $props();
	let lastData: StudentTermsListResponse = $state({
		items: [],
		itemsCount: 0
	});

	let autoUpdate = $state(false);
	let autoUpdateInterval: number | undefined = $state();

	$effect(() => {
		if (autoUpdate && !autoUpdateInterval) {
			autoUpdateInterval = setInterval(() => {
				invalidate((url) => {
					return url.href.endsWith('/terms');
				});
			}, 1000);
		} else if (!autoUpdate) {
			clearInterval(autoUpdateInterval);
			autoUpdateInterval = undefined;
		}
	});

	$effect(() => {
		data.terms
			.then((res) => {
				lastData = res;
				isLoading = false;
			})
			.catch(() => {});
	});

	const getItemClass = (item: StudentTermDTO) => {
		if (item.joined) {
			if (item.canJoin || item.canLeave) {
				return 'bg-green-200 hover:bg-green-100 dark:bg-green-800';
			} else {
				return 'bg-green-200 hover:bg-green-100 opacity-50 dark:bg-green-900';
			}
		} else {
			if (item.canJoin || item.canLeave) {
				return '';
			} else {
				return 'opacity-50';
			}
		}
	};

	const joinTerm = async (term: StudentTermDTO) => {
		if (term.willSignOut) {
			if (!confirm(m.signout_warning())) {
				return;
			}
		}

		await API.request<any, TermsJoinResponse>(
			`/api/v2/courses/${courseId}/items/${term.courseItemId}/terms/${term.id}/students`,
			{
				method: 'POST',
				body: {}
			}
		)
			.then((res) => {
				invalidate((url) => {
					return url.href.endsWith('/terms');
				});
			})
			.catch(() => {});
	};

	const leaveTerm = async (term: StudentTermDTO) => {
		await API.request<any, TermsJoinResponse>(
			`/api/v2/courses/${courseId}/items/${term.courseItemId}/terms/${term.id}/students`,
			{
				method: 'DELETE',
				body: {}
			}
		)
			.then((res) => {
				invalidate((url) => {
					return url.href.endsWith('/terms');
				});
			})
			.catch(() => {});
	};

	onDestroy(() => {
		if (autoUpdateInterval) {
			clearInterval(autoUpdateInterval);
		}
	});
</script>

{#snippet TableHeader()}
	<Table.Row>
		<Table.Head>{m.term_name()}</Table.Head>
		<Table.Head>{m.term_activityname()}</Table.Head>
		<Table.Head>{m.term_active()}</Table.Head>
		<Table.Head>{m.term_signin()}</Table.Head>
		<Table.Head>{m.term_signout()}</Table.Head>
		<Table.Head>{m.term_classroom()}</Table.Head>
		<Table.Head>{m.term_students()}</Table.Head>
		<Table.Head>{m.actions()}</Table.Head>
	</Table.Row>
{/snippet}

<div class="flex flex-col gap-8 m-8">
	{#if isLoading}
		<Pageloader></Pageloader>
	{:else}
		<div class="flex flex-row justify-between">
			<h1 class="mb-8 text-2xl">Termíny:</h1>
		</div>

		<div class="flex flex-col gap-4">
			<div class="flex justify-end gap-4">
				<Checkbox bind:checked={autoUpdate} name="periodicUpdates" id="periodicUpdates"></Checkbox>
				<Label for="periodicUpdates">{m.toggle_periodic_update()}</Label>
			</div>

			<Table.Root>
				<Table.Header class="font-medium border-t bg-muted/50">
					{@render TableHeader()}
				</Table.Header>
				<Table.Body>
					{#if lastData.items.length === 0}
						<Table.Row>
							<Table.Cell colspan={7}>No items found</Table.Cell>
						</Table.Row>
					{:else}
						{#each lastData.items as item}
							<Table.Row class={getItemClass(item)}>
								<Table.Cell>{item.name}</Table.Cell>
								<Table.Cell>{item.courseItemName}</Table.Cell>
								<Table.Cell>
									<DateRange start={item.activeFrom} end={item.activeTo} showTime={true}
									></DateRange>
								</Table.Cell>
								<Table.Cell>
									<DateRange start={item.signInFrom} end={item.signInTo} showTime={true}
									></DateRange>
								</Table.Cell>
								<Table.Cell>
									<DateRange start={item.signOutFrom} end={item.signOutTo} showTime={true}
									></DateRange>
								</Table.Cell>
								<Table.Cell>{item.classroom}</Table.Cell>
								<Table.Cell>
									{item.studentsJoined}/{item.studentsMax}
								</Table.Cell>
								<Table.Cell>
									{#if item.joined}
										<Button disabled={!item.canLeave} onclick={() => leaveTerm(item)}
											>{m.term_leave()}</Button
										>
									{:else}
										<Button disabled={!item.canJoin} onclick={() => joinTerm(item)}
											>{m.term_join()}</Button
										>
									{/if}
								</Table.Cell>
							</Table.Row>
						{/each}
					{/if}
				</Table.Body>
				{#if lastData.itemsCount > 10}
					<Table.Footer>
						{@render TableHeader()}
					</Table.Footer>
				{/if}
			</Table.Root>
		</div>
	{/if}
</div>
