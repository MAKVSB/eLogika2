<script lang="ts">
	import {
		type CategoryGetByIdResponse,
		type CategoryDTO,
		type CategoryListRequest,
		type CategoryListResponse,
		type ChapterListRequest,
		type ChapterListResponse
	} from '$lib/api_types';
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import * as Form from '$lib/components/ui/form';
	import Label from '$lib/components/ui/label/label.svelte';
	import { m } from '$lib/paraglide/messages';
	import { API } from '$lib/services/api.svelte';
	import { onMount } from 'svelte';

	let {
		form = $bindable(),
		courseId
	}: {
		form: any;
		courseId: string;
	} = $props();

	let changeCounter = $state(0);

	let chapters: Form.SelectOptions = $state([]);
	let chapterLoading = $state(false);

	let categories: Form.SelectOptions = $state([]);
	let categoryLoading = $state(false);

	let category: CategoryDTO | undefined = $state();

	export async function loadCategory(courseId: string, categoryId: number) {
		categoryLoading = true;
		if (categoryId) {
			await API.request<null, CategoryGetByIdResponse>(
				`/api/v2/courses/${courseId}/categories/${categoryId}`
			)
				.then((res) => {
					category = res.data;
				})
				.catch(() => {});
		} else {
			category = undefined;
		}
		categoryLoading = false;
	}

	export async function loadCategories(courseId: string, chapterId?: number) {
		chapterLoading = true;
		if (chapterId) {
			await API.request<CategoryListRequest, CategoryListResponse>(
				`/api/v2/courses/${courseId}/chapters/${chapterId}/categories`
			)
				.then((res) => {
					categories = res.items.map((v) => {
						return {
							display: v.name,
							value: v.id
						};
					});
				})
				.catch(() => {});

			chapterLoading = false;
		} else {
			categories = [];
			chapterLoading = false;
		}
	}

	export async function loadChapters(courseId: string) {
		await API.request<ChapterListRequest, ChapterListResponse>(
			`/api/v2/courses/${courseId}/chapters`
		)
			.then((res) => {
				chapters = res.items.map((v) => {
					return {
						display: v.name,
						value: v.id
					};
				});
			})
			.catch(() => {});
	}

	async function chapterChanged() {
		await loadCategories(courseId, form.fields.chapterId);
		form.fields.categoryId = null;
		form.fields.steps = [];
		category = undefined;
	}

	async function categoryChanged() {
		await loadCategory(courseId, form.fields.categoryId);
		form.fields.steps = [];
	}

	onMount(async () => {
		await Promise.all([
			loadChapters(courseId),
			loadCategories(courseId, form.fields.chapterId),
			loadCategory(courseId, form.fields.categoryId)
		]);
	});

	function toggleStep(stepId: number) {
		if (form.fields.steps.includes(stepId)) {
			form.fields.steps = form.fields.steps.filter((s: number) => s != stepId);
		} else {
			form.fields.steps.push(stepId);
		}
		changeCounter += 1;
	}
</script>

<div class="flex flex-col gap-4 p-2 my-4">
	<div class="grid grid-cols-12 gap-4">
		<Form.SingleSelect
			title={m.question_grouping_chapter()}
			name="chapter"
			id="chapter"
			placeholder="Select chapter"
			class="col-span-12 sm:col-span-6"
			bind:value={form.fields.chapterId}
			options={chapters}
			onchange={chapterChanged}
			error={form.errors.chapterId}
		></Form.SingleSelect>
		<Form.SingleSelect
			title={m.question_grouping_category()}
			name="category"
			id="category"
			placeholder="Select category"
			class="col-span-12 sm:col-span-6"
			nullable={true}
			bind:value={form.fields.categoryId}
			options={categories}
			onchange={categoryChanged}
			error={form.errors.categoryId}
			loading={chapterLoading}
		></Form.SingleSelect>

		<h1 class="col-span-12 text-xl">{m.question_grouping_steps()}</h1>
		{#if category?.steps}
			<ul class="flex flex-col w-full col-span-12 gap-2">
				{#each category.steps as step}
					{#key changeCounter}
						<li class="flex gap-4">
							<Checkbox
								id={'step' + step.id}
								checked={form.fields.steps.includes(step.id)}
								onclick={() => toggleStep(step.id)}
							></Checkbox>
							<Label for={'step' + step.id}>
								{step.name} ({step.difficulty})
							</Label>
						</li>
					{/key}
				{/each}
			</ul>
		{:else}
			<p class="col-span-12">{m.question_grouping_steps_missing()}</p>
		{/if}
	</div>
</div>
