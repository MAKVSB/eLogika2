<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { API, ApiError } from '$lib/services/api.svelte';
	import * as Form from '$lib/components/ui/form';
	import { toast } from 'svelte-sonner';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import { Label } from '$lib/components/ui/label';
	import { Button } from '$lib/components/ui/button';
	import { TipTapDefaultContent } from '$lib/constants';
	import { m } from '$lib/paraglide/messages.js';
	import type {
		ChapterMoveResponse,
		MoveDirectionEnum,
		ChapterGetByIdResponse,
		ChapterInsertRequest,
		ChapterInsertResponse,
		ChapterUpdateResponse,
		ChapterDTO
	} from '$lib/api_types';
	import CategoryList from './CategoryList/CategoryList.svelte';
	import { ChapterInsertRequestSchema } from '$lib/schemas';
	import { base } from '$app/paths';

	let courseId = $derived<string | null>(page.params.courseId);
	let { data } = $props();

	$effect(() => {
		if (data.course) {
			data.course.then((data) => setResult(data)).catch(() => {});
		} else {
			form.fields = defaultFormData;
		}
	});

	const defaultFormData: ChapterDTO = {
		id: 0,
		// svelte-ignore state_referenced_locally
		courseId: Number(courseId),
		version: 0,
		name: '',
		content: TipTapDefaultContent,
		visible: true,
		order: 0,
		childs: []
	};
	let form = $state(Form.createForm(ChapterInsertRequestSchema, defaultFormData));

	const actionColumn = columns.find((c) => c.uniqueId == 'actions');
	if (actionColumn) {
		actionColumn.meta = {
			...(columns[0].meta ?? {}),
			changeEventHandler: (id: number, direction: MoveDirectionEnum) => {
				if (!courseId) return;
				API.request<null, ChapterMoveResponse>(
					`/api/v2/courses/${courseId}/chapters/${id}/move/${direction}`,
					{
						method: 'PATCH'
					}
				)
					.then((res) => {
						form.fields.childs = res.childs.sort((a, b) => a.order - b.order);
					})
					.catch(() => {});

				return true;
			}
		};
	}

	function setResult(res: ChapterGetByIdResponse | ChapterInsertResponse | ChapterUpdateResponse) {
		res.data = {
			...res.data,
			childs: res.data.childs.sort((a, b) => a.order - b.order)
		};
		form.fields = res.data;
		console.log("Transfering 20")
		goto(String(res.data.id), {
			replaceState: true
		});
	}

	async function handleSubmit(): Promise<any> {
		let request;
		if (data.creating) {
			const parentId = page.url.searchParams.get('parentId');
			if (!parentId) {
				toast.error('Failed to detect parent chapter');
				return;
			}

			request = API.request<ChapterInsertRequest, ChapterInsertResponse>(
				`/api/v2/courses/${page.params.courseId}/chapters/${parentId}`,
				{
					method: 'POST',
					body: {
						...form.fields
					}
				}
			);
		} else {
			request = API.request<ChapterInsertRequest, ChapterInsertResponse>(
				`/api/v2/courses/${page.params.courseId}/chapters/${page.params.id}`,
				{
					method: 'PUT',
					body: {
						...form.fields
					}
				}
			);
		}

		return request
			.then((res) => {
				setResult(res);
			});
	}
</script>

<div class="m-8">
	{#await data.course}
		<Pageloader></Pageloader>
	{:then staticResourceData}
		<div class="flex flex-row justify-between">
			<h1 class="mb-8 text-2xl">
				Chapter management:
				<b>
					{staticResourceData?.data?.name ?? 'New subchapter'}
				</b>
			</h1>
		</div>
		{#if !data.creating}
			<div class="flex justify-between gap-2">
				<Button
					disabled={!staticResourceData?.data?.parentId}
					href={String(staticResourceData?.data?.parentId)}
				>
					{m.chapter_parent()}
				</Button>
				<Button href="{base}/app/{courseId}/tutor/chapters/0?parentId={form.fields.id}"
					>{m.chapter_subchapter_add()}</Button
				>
			</div>
			<div class="flex flex-col gap-2 p-2">
				<Label class="flex justify-between">
					{m.chapter_subchapters()}
					<p class="text-xs">({m.table_changes_instant()})</p>
				</Label>
				<DataTable data={form.fields.childs ?? []} {columns} {filters} paginationEnabled={false} queryParam='search'/>
			</div>
		{/if}
		<Form.Root bind:form onsubmit={handleSubmit} isCreating={data.creating}>
			<div class="flex flex-col gap-4 p-2">
				<div class="grid grid-cols-4 gap-4">
					<Form.TextInput
						title={m.chapter_name()}
						name="name"
						id="name"
						type="text"
						class="col-span-4 sm:col-span-3"
						bind:value={form.fields.name}
						error={form.errors.name}
					></Form.TextInput>
					<Form.Checkbox
						title={m.chapter_visible()}
						name="visible"
						id="visible"
						wide={true}
						class="col-span-4 sm:col-span-1"
						bind:value={form.fields.visible}
						error={form.errors.visible}
					></Form.Checkbox>
					{#key form.fields}
						<Form.Tiptap
							title={m.chapter_content()}
							name="content"
							id="content"
							class="col-span-4"
							bind:value={form.fields.content}
							error={form.errors.content}
							enableFileUpload
							enableFileLink
						></Form.Tiptap>
					{/key}
				</div>
			</div>
		</Form.Root>
		{#if !data.creating}
			<CategoryList></CategoryList>
		{/if}
	{/await}
</div>
