<script lang="ts">
	import { page } from '$app/state';
	import {
		TestInstanceFormEnum,
		type JoinedStudentDTO,
		type ListJoinedStudentsResponse,
		type TestGeneratorRequest,
		type TestGeneratorResponse
	} from '$lib/api_types';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { API } from '$lib/services/api.svelte';
	import { columns, filters } from './schema';
	import { DataTable } from '$lib/components/ui/data-table';
	import type { InitialTableState, RowSelectionState, TableState } from '@tanstack/table-core';
	import { Label } from '$lib/components/ui/label/index.js';
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import { Input } from '$lib/components/ui/input';
	import { invalidateAll } from '$app/navigation';
	import * as Form from '$lib/components/ui/form';
	import { enumToOptions } from '$lib/utils';
	import { m } from '$lib/paraglide/messages';

	$effect(() => {
		searchCache = page.url.searchParams.get('generateTestSearch')
	})

	function refetch(_: any, str: string) {
		searchCache = str
	}

	$effect(() => {
		loadData(searchCache, skipUsersWithInstance)
	})
	
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

	let rowItems: JoinedStudentDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({
		pagination: {
			pageIndex: 0,
			pageSize: 10000
		}
	});
	let searchCache: string | null = $state("")

	let instanceForm = $state(TestInstanceFormEnum.OFFLINE);
	let skipUsersWithInstance = $state(true);


	const loadData = (search: string | null, skipUsersWithInstance: boolean) => {
		API.request<null, ListJoinedStudentsResponse>(
			`/api/v2/courses/${page.params.courseId}/items/${page.params.itemId}/terms/${termId}/students`,
			{
				searchParams: {
					...(search ? { search: search } : {}),
					...(skipUsersWithInstance ? { skipUsersWithInstance: String(skipUsersWithInstance) } : {})
				}
			},
			fetch
		)
			.then((res) => {
				rowItems = res.items;
				rowCount = res.itemsCount;
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
					usersAll: usersAll ?? users.length == 0,
					usersIds: users,
					form: instanceForm,
					skipUsersWithInstance: skipUsersWithInstance
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
</script>

<Dialog.Content class="max-h-full w-300 overflow-scroll sm:max-h-[90%] sm:max-w-[90%]">
	<Dialog.Header>
		<Dialog.Title>{m.tests_generate()}</Dialog.Title>
	</Dialog.Header>

	<Form.SingleSelect
		title="Instance form"
		name="instanceForm"
		id="instandeForm"
		bind:value={instanceForm}
		options={enumToOptions(TestInstanceFormEnum, m.test_instance_form_enum)}
		error=""
	></Form.SingleSelect>

	<div class="flex gap-2">
		<Checkbox
			class="rounded-md h-9 w-9"
			name="generateSigned"
			id="generateSigned"
			bind:checked={generateForSigned}
		/>
		<Label for="generateSigned">{m.test_generate_signed()}</Label>
	</div>
	{#if generateForSigned}
		<div class="flex gap-2">
			<Checkbox
				class="rounded-md h-9 w-9"
				name="skipUsersWithInstance"
				id="skipUsersWithInstance"
				bind:checked={skipUsersWithInstance}
			/>
			<Label for="skipUsersWithInstance">{m.test_generate_skip()}</Label>
		</div>

		<DataTable
			data={rowItems}
			{columns}
			{filters}
			{selection}
			{rowCount}
			paginationEnabled={false}
			{initialState}
			queryParam="generateTestSearch"
			{refetch}
			replaceState
		/>
	{:else}
		<div class="flex flex-col gap-2">
			<Label for="generateNumber">{m.test_generate_number()}</Label>
			<Input
				id="generateNumber"
				bind:value={numberToGenerate}
				placeholder={m.test_generate_number()}
				required
				type="number"
			/>
		</div>
	{/if}

	<Button onclick={() => generateAndClose()}>{m.tests_generate()}</Button>
</Dialog.Content>
