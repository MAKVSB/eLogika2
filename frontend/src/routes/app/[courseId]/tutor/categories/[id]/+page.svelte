<script lang="ts">
	import { onMount } from 'svelte';
	import { API, ApiError } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import * as Form from '$lib/components/ui/form';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import type {
		CategoryGetByIdResponse,
		CategoryInsertRequest,
		CategoryInsertResponse,
		CategoryUpdateRequest,
		CategoryUpdateResponse,
		ChapterListRequest,
		ChapterListResponse,
		StepDTO
	} from '$lib/api_types';
	import { Button } from '$lib/components/ui/button';
	import { m } from '$lib/paraglide/messages';
	import Label from '$lib/components/ui/label/label.svelte';
	import type { ErrorObject } from '$lib/components/ui/form/types';
	import type { SelectOptions } from '$lib/components/ui/form/form-input-select-single.svelte';
	import { CategoryInsertRequestSchema } from '$lib/schemas';

	let courseId = $derived<string | null>(page.params.courseId);
	let { data } = $props();

	$effect(() => {
		if (data.category) {
			data.category.then((data) => setResult(data)).catch(() => {});
		} else {
			form.fields = defaultFormData;
		}
	});

	let changeCounter = $state(0);
	let chapters: SelectOptions = $state([]);

	const defaultFormData = {
		name: '',
		steps: [] as StepDTO[],
		chapterId: 0,
		version: 0
	};

	let form = $state(Form.createForm(CategoryInsertRequestSchema, defaultFormData));

	onMount(async () => {
		API.request<ChapterListRequest, ChapterListResponse>(`/api/v2/courses/${courseId}/chapters`)
			.then((res) => {
				chapters = res.items.map((c) => {
					return {
						value: c.id,
						display: c.name
					};
				});
			})
			.catch(() => {});
	});

	function setResult(
		res: CategoryGetByIdResponse | CategoryInsertResponse | CategoryUpdateResponse
	) {
		form.fields = res.data;
		changeCounter += 1;
		console.log('Transfering 19');
		goto(String(res.data.id), {
			replaceState: true
		});
	}

	async function handleSubmit(): Promise<any> {
		let request;
		if (data.creating) {
			request = API.request<CategoryInsertRequest, CategoryInsertResponse>(
				`/api/v2/courses/${courseId}/categories`,
				{
					method: 'POST',
					body: form.fields
				}
			);
		} else {
			request = API.request<CategoryUpdateRequest, CategoryUpdateResponse>(
				`/api/v2/courses/${courseId}/categories/${page.params.id}`,
				{
					method: 'PUT',
					body: form.fields
				}
			);
		}

		return request.then((res) => {
			setResult(res);
			if (data.creating) {
				toast.success('Created succesfully');
			} else {
				toast.success('Saved succesfully');
			}
		});
	}

	function deleteAnswerAtIndex(index: number) {
		form.fields.steps.splice(index, 1);
		changeCounter += 1;
	}

	function addAnswer() {
		form.fields.steps.push({
			id: 0,
			name: '',
			difficulty: 1,
			deleted: false
		});
		changeCounter += 1;
	}
</script>

<div class="m-8">
	{#await data.category}
		<Pageloader></Pageloader>
	{:then staticResourceData}
		<div class="flex flex-row justify-between">
			<h1 class="mb-8 text-2xl">
				Category management:
				<b>
					{staticResourceData?.data.name ?? 'New category'}
				</b>
			</h1>
		</div>
		<Form.Root bind:form onsubmit={handleSubmit} isCreating={data.creating}>
			<div class="flex flex-col gap-4 p-2">
				<div class="grid grid-cols-12 gap-2">
					<Form.TextInput
						title={m.category_name()}
						name="name"
						id="name"
						type="text"
						class="col-span-12 sm:col-span-6"
						bind:value={form.fields.name}
						error={form.errors.name}
					></Form.TextInput>
					<Form.SingleSelect
						title={m.category_chapter()}
						name="chapter"
						id="chapter"
						class="col-span-12 sm:col-span-6"
						bind:options={chapters}
						bind:value={form.fields.chapterId}
						error={form.errors.chapter}
					></Form.SingleSelect>
				</div>
				<div class="flex flex-col gap-2">
					<div class="flex items-center justify-between">
						<h1 class="text-xl">{m.category_steps()}:</h1>
						<Button variant="outline" onclick={addAnswer}>{m.step_add()}</Button>
					</div>

					<div class="relative flex flex-col gap-4 p-4 border rounded-md">
						{#key changeCounter}
							{#each form.fields.steps.filter((ans) => !ans.deleted) as answer, index}
								{@const stepsErrors = (form.errors.steps as ErrorObject) ?? {}}
								{@const stepErrors = (stepsErrors[String(index)] as ErrorObject) ?? ''}

								<div class="grid grid-cols-12 gap-4">
									<Form.TextInput
										title="Name"
										name="name"
										id="name"
										type="text"
										class="col-span-12 sm:col-span-5"
										bind:value={answer.name}
										error={stepErrors.name}
									></Form.TextInput>
									<Form.TextInput
										title="Difficulty"
										name="difficulty"
										id="difficulty"
										type="number"
										class="col-span-6 sm:col-span-5"
										bind:value={answer.difficulty}
										error={stepErrors.difficulty}
									></Form.TextInput>
									<div class="flex flex-col col-span-6 gap-2 sm:col-span-2">
										<Label>&nbsp;</Label>
										<Button
											variant="destructive"
											onclick={() => {
												if (answer.id !== 0) {
													answer.deleted = true;
												} else {
													deleteAnswerAtIndex(index);
												}
											}}>{m.delete()}</Button
										>
									</div>
								</div>
							{/each}
						{/key}
					</div>
				</div>
			</div>
		</Form.Root>
	{/await}
</div>
