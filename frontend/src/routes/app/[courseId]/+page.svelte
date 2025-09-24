<script lang="ts">
	import { API, ApiError } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import {
		type CourseInsertRequest,
		type CourseInsertResponse,
		SemesterEnum,
		type CourseDTO,
		type CourseUpdateRequest,
		type CourseUpdateResponse,
		StudyFormEnum,
		type CourseGetByIdResponse
	} from '$lib/api_types';
	import * as Form from '$lib/components/ui/form';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import { TipTapDefaultContent } from '$lib/constants';
	import { enumToOptions } from '$lib/utils';
	import { m } from '$lib/paraglide/messages';
	import { CourseInsertRequestSchema } from '$lib/schemas';
	import TiptapRenderer from '$lib/components/tiptap-renderer/tiptap-renderer.svelte';

	let { data } = $props();

	$effect(() => {
		if (data.course) {
			data.course.catch((err) => {
				console.error(err);
				toast.error('Failed to load course info');
			});
		}
	});
</script>

<div class="flex flex-col gap-8 m-8">
	{#await data.course}
		<Pageloader></Pageloader>
	{:then staticResourceData}
		{#if staticResourceData && staticResourceData.data}
			<div class="flex flex-col justify-between gap-4">
				<h1 class="text-2xl">
					Course: <b>{staticResourceData.data.name} ({staticResourceData.data.shortname})</b>
				</h1>
				<h2 class="text-xl">
					Academic year: <b>{staticResourceData.data.year} - {staticResourceData.data.year + 1}</b>
				</h2>
				<h2 class="text-xl">
					Semester: <b>{m.semester_enum({ value: staticResourceData.data.semester })}</b>
				</h2>
			</div>
			<div>
				{#if staticResourceData.data.content}
					<TiptapRenderer jsonContent={staticResourceData.data.content}></TiptapRenderer>
				{:else}
					No content to display
				{/if}
			</div>
		{/if}
	{/await}
</div>
