<script lang="ts">
	import { page } from '$app/state';
	import {
		CourseUserRoleEnum,
		TestInstanceFormEnum,
		type CourseUserDTO,
		type ListCourseUsersResponse,
		type TestGeneratorRequest,
		type TestGeneratorResponse
	} from '$lib/api_types';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { API, ApiError, encodeJsonToBase64Url } from '$lib/services/api.svelte';
	import { onMount } from 'svelte';
	import { columns, filters } from './schema';
	import { DataTable } from '$lib/components/ui/data-table';
	import type { InitialTableState, RowSelectionState } from '@tanstack/table-core';
	import Loader from '$lib/components/ui/loader/loader.svelte';
	import { Label } from '$lib/components/ui/label/index.js';
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import { Input } from '$lib/components/ui/input';
	import { toast } from 'svelte-sonner';
	import { invalidateAll } from '$app/navigation';
	import * as Form from '$lib/components/ui/form';
	import { enumToOptions } from '$lib/utils';

	let {
		termId,
		openState = $bindable()
	}: {
		termId: number;
		openState: boolean;
	} = $props();

	let generateForSigned = $state(false);
	let numberToGenerate = $state(0);
	let usersAll = $state(false);
	let users: number[] = $state([]);

	let rowItems: CourseUserDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});
	let loading = $state(true);

	let instanceForm = $state(TestInstanceFormEnum.OFFLINE);

	const loadData = () => {
		initialState.columnFilters = [
			{
				id: 'role',
				value: CourseUserRoleEnum.STUDENT
			}
		];
		initialState.pagination = {
			pageIndex: 0,
			pageSize: 10000
		};
		const search = encodeJsonToBase64Url(initialState);

		API.request<null, ListCourseUsersResponse>(
			`/api/v2/courses/${page.params.courseId}/items/${page.params.itemId}/terms/${termId}/students`,
			{
				searchParams: {
					...(search ? { search: search } : {})
				}
			},
			fetch
		)
			.then((res) => {
				rowItems = res.items;
				rowCount = res.items.length;
			})
			.catch(() => {});
	};

	function selection(state: RowSelectionState, all: boolean) {
		if (all) {
			usersAll = true;
			users = [];
		} else {
			usersAll = false;
			users = Object.entries(state).map((s) => {
				return Number(s[0]);
			});
		}
	}

	const generateAndClose = () => {
		API.request<TestGeneratorRequest, TestGeneratorResponse>(
			`/api/v2/courses/${page.params.courseId}/tests/${page.params.itemId}/${termId}/generate`,
			{
				method: 'POST',
				body: {
					variants: numberToGenerate,
					usersAll: usersAll,
					usersIds: users,
					form: instanceForm
				}
			},
			fetch
		)
			.then(() => {
				invalidateAll();
				openState = false;
			})
			.catch(() => {});
	};

	onMount(() => {
		loadData();
		loading = false;
	});
</script>

<Dialog.Content class="max-h-full w-300 overflow-scroll sm:max-h-[90%] sm:max-w-[90%]">
	<Dialog.Header>
		<Dialog.Title>Generate tests</Dialog.Title>
	</Dialog.Header>

	<Form.SingleSelect
		title="Instance form"
		name="instanceForm"
		id="instandeForm"
		bind:value={instanceForm}
		options={enumToOptions(TestInstanceFormEnum)}
		error=""
	></Form.SingleSelect>

	<div class="flex gap-2">
		<Checkbox
			class="rounded-md h-9 w-9"
			name="generateSigned"
			id="generateSigned"
			bind:checked={generateForSigned}
		/>
		<Label for="generateSigned">Generate for signed students</Label>
	</div>
	{#if generateForSigned}
		{#if !loading}
			<DataTable
				data={rowItems}
				{columns}
				{filters}
				{selection}
				{rowCount}
				paginationEnabled={false}
				{initialState}
			/>
		{:else}
			<Loader></Loader>
		{/if}
	{:else}
		<div class="flex flex-col gap-2">
			<Label for="generateNumber">Number of tests to generate</Label>
			<Input
				id="generateNumber"
				bind:value={numberToGenerate}
				placeholder={'Number of tests to generate'}
				required
				type="number"
			/>
		</div>
	{/if}

	<Button onclick={() => generateAndClose()}>Generate tests</Button>
</Dialog.Content>
