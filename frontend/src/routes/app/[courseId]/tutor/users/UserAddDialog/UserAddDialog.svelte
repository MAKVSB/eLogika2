<script lang="ts">
	import {
		CourseUserRoleEnum,
		StudyFormEnum,
		type AddCourseUserRequest,
		type AddCourseUserResponse,
		type CourseUserDTO,
		type ListCourseUsersResponse
	} from '$lib/api_types';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { API, encodeJsonToBase64Url } from '$lib/services/api.svelte';
	import { onMount } from 'svelte';
	import { columns, filters } from './schema';
	import { DataTable } from '$lib/components/ui/data-table';
	import type {
		ColumnFiltersState,
		InitialTableState,
		SortingState,
		TableState
	} from '@tanstack/table-core';
	import Loader from '$lib/components/ui/loader/loader.svelte';
	import { invalidate } from '$app/navigation';
	import * as Form from '$lib/components/ui/form';
	import { enumToOptions } from '$lib/utils';
	import { m } from '$lib/paraglide/messages';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import { DataTableSearchParams } from '$lib/api_types_static';

	let {
		defaultRole,
		endpoint
	}: {
		defaultRole?: string;
		endpoint: string;
	} = $props();

	let newUserRole: CourseUserRoleEnum = $state(CourseUserRoleEnum.STUDENT);
	let newUserStudyForm: StudyFormEnum | null = $state(StudyFormEnum.FULLTIME);

	let search = $state('');

	let rowItems: CourseUserDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});
	let loading = $state(true);

	const actionsColumn = columns.find((c) => c.uniqueId == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (id: number) => {
				await API.request<AddCourseUserRequest, AddCourseUserResponse>(endpoint, {
					method: 'POST',
					body: {
						userId: id,
						role: newUserRole,
						...(newUserStudyForm ? { studyForm: newUserStudyForm } : {})
					}
				})
					.then((res) => {
						invalidate((url) => {
							return url.pathname.endsWith(`courses/${page.params.courseId}/users`);
						});
						toast.success(m.toast_course_user_added({role: m.course_user_role_enum({value: newUserRole})}))
					})
					.catch(() => {});

				return true;
			}
		};
	}

	const loadData = () => {
		API.request<null, ListCourseUsersResponse>(
			`/api/v2/users`,
			{
				searchParams: {
					...(search ? { search: search } : {})
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

	function refetch(state: TableState) {
		search = DataTableSearchParams.fromDataTable(state).toURL();
		loadData();
	}

	onMount(() => {
		search = encodeJsonToBase64Url(initialState);
		loadData();
		loading = false;
	});
</script>

<Dialog.Content class="max-h-full w-300 overflow-scroll sm:max-h-[90%] sm:max-w-[90%]">
	<Dialog.Header>
		<Dialog.Title>{m.course_user_add_modal_settings()}</Dialog.Title>
	</Dialog.Header>
	<div class="grid grid-cols-2 gap-4">
		<Form.SingleSelect
			title={m.course_user_add_modal_settings_role()}
			name="newUserRole"
			id="newUserRole"
			bind:value={newUserRole}
			required={true}
			error=""
			nullable={false}
			options={enumToOptions(CourseUserRoleEnum, m.course_user_role_enum)}
		></Form.SingleSelect>
		<Form.SingleSelect
			title={m.course_user_add_modal_settings_studyform()}
			name="newUserStudyForm"
			id="newUserStudyForm"
			bind:value={newUserStudyForm}
			required={newUserRole == CourseUserRoleEnum.STUDENT}
			error=""
			nullable={false}
			options={enumToOptions(StudyFormEnum, m.study_form_enum)}
			disabled={newUserRole != CourseUserRoleEnum.STUDENT}
		></Form.SingleSelect>
	</div>
	<Dialog.Header>
		<Dialog.Title>{m.course_user_add_modal_userselect()}</Dialog.Title>
	</Dialog.Header>
	{#if !loading}
		<DataTable
			data={rowItems}
			{columns}
			{filters}
			{refetch}
			{rowCount}
			queryParam="usersearch"
			{initialState}
		/>
	{:else}
		<Loader></Loader>
	{/if}
</Dialog.Content>
