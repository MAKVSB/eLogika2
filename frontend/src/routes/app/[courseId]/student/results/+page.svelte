<script lang="ts">
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import * as Table from '$lib/components/ui/table/index.js';
	import { m } from '$lib/paraglide/messages';
	import CourseItemDisplay from './CourseItemDisplay.svelte';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { getLocale } from '$lib/paraglide/runtime';
	import Button from '$lib/components/ui/button/button.svelte';
	import { page } from '$app/state';
	import { base } from '$app/paths';

	let { data } = $props();
</script>

{#snippet TableHeader()}
	<Table.Row>
		<Table.Head>{m.course_item_name()}</Table.Head>
		<Table.Head>{m.course_item_type()}</Table.Head>
		<Table.Head>{m.testinstance_points()}</Table.Head>
		<Table.Head>{m.course_item_mandatory()}</Table.Head>
		<Table.Head>{m.course_item_passed()}</Table.Head>
		<Table.Head>{m.course_item_tries()}</Table.Head>
	</Table.Row>
{/snippet}

{#snippet TableHeaderResults()}
	<Table.Row>
		<Table.Head>{m.courseitem_result_name()}</Table.Head>
		<Table.Head>{m.courseitem_result_group()}</Table.Head>
		<Table.Head>{m.term_name()}</Table.Head>
		<Table.Head>{m.testinstance_points()}</Table.Head>
		<!-- <Table.Head>Attempt</Table.Head> -->
		<Table.Head>{m.updatedby()}</Table.Head>
		<Table.Head>{m.result_selected()}</Table.Head>
		<Table.Head>{m.actions()}</Table.Head>
	</Table.Row>
{/snippet}

<div class="flex flex-col gap-8 m-8">
	{#await data.courseItems}
		<Pageloader></Pageloader>
	{:then staticResourceData}
		<div class="flex flex-col gap-4">
			<h1 class="text-2xl">{m.menu_student_results()}:</h1>
			<Table.Root>
				<Table.Header class="font-medium border-t bg-muted/50">
					{@render TableHeader()}
				</Table.Header>
				<Table.Body>
					{#each staticResourceData.items as item}
						<CourseItemDisplay {item}></CourseItemDisplay>
					{/each}
				</Table.Body>
				<Table.Footer>
					{@render TableHeader()}
					<Table.Row>
						<Table.Cell>{m.points_total()}</Table.Cell>
						<Table.Cell></Table.Cell>
						<Table.Cell
							>{staticResourceData.items.reduce((accumulator, current) => {
								return accumulator + Math.min(current.points, current.pointsMax);
							}, 0)}</Table.Cell
						>
						<Table.Cell></Table.Cell>
						<Table.Cell></Table.Cell>
						<Table.Cell></Table.Cell>
					</Table.Row>
				</Table.Footer>
			</Table.Root>
			<p>* : Result is not labeled as final. Waiting for teacher to assign points</p>
		</div>

		<div class="flex flex-col gap-4">
			<h1 class="text-2xl">{m.results_testactivities()}:</h1>
			<Table.Root>
				<Table.Header class="font-medium border-t bg-muted/50">
					{@render TableHeaderResults()}
				</Table.Header>
				<Table.Body>
					{#each staticResourceData.results as result, i}
						<Table.TableRow>
							<Table.Cell>{result.courseItemName}</Table.Cell>
							<Table.Cell>{result.courseItemGroupName}</Table.Cell>
							<Table.Cell>
								<Tooltip.Provider>
									<Tooltip.Root>
										<Tooltip.Trigger>{result.termName}</Tooltip.Trigger>
										<Tooltip.Content class="grid grid-cols-2">
											<p>Active from:</p>
											<p>
												{new Date(result.termActiveFrom).toLocaleString(getLocale())}
											</p>
											<p>Active to:</p>
											<p>
												{new Date(result.termActiveTo).toLocaleString(getLocale())}
											</p>
											<p>Instance started at:</p>
											<p>
												{new Date(result.instanceStartTime).toLocaleString(getLocale())}
											</p>
										</Tooltip.Content>
									</Tooltip.Root>
								</Tooltip.Provider>
							</Table.Cell>
							<Table.Cell
								>{result.points}
								{#if !result.final}*{/if}
							</Table.Cell>
							<!--<Table.Cell>TODO attempt</Table.Cell>-->
							<Table.Cell>
								<Tooltip.Provider>
									<Tooltip.Root>
										<Tooltip.Trigger>{result.updatedBy}</Tooltip.Trigger>
										<Tooltip.Content>
											<p>
												{new Date(result.updatedAt).toLocaleString(getLocale())}
											</p>
										</Tooltip.Content>
									</Tooltip.Root>
								</Tooltip.Provider>
							</Table.Cell>
							<Table.Cell>{m.yes_no({ value: String(result.selected) })}</Table.Cell>
							<Table.Cell>
								{#if result.activityInstanceId}
									<Button
										variant="default"
										href="{base}/app/{page.params
											.courseId}/student/activities/{result.activityInstanceId}"
									>
										{m.view()}
									</Button>
								{:else}
									<!-- TODO  test: {result.testInstanceId} -->
								{/if}
							</Table.Cell>
						</Table.TableRow>
					{/each}
				</Table.Body>
				<Table.Footer>
					{#if staticResourceData.results.length > 10}
						{@render TableHeaderResults()}
					{/if}
				</Table.Footer>
			</Table.Root>
			<p>* : Result is not labeled as final. Waiting for teacher to assign points</p>
		</div>
	{/await}
</div>
