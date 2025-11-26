import { renderComponent, SortButton, type ColDef } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import { CourseUserRoleEnum, type CourseUserDTO } from '$lib/api_types';
import { FilterTypeEnum, type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import { enumToOptions } from '$lib/utils';
import DataTableRoles from '$lib/components/ui/data-table/data-table-roles.svelte';
import GlobalState from '$lib/shared.svelte';
import type { InitialTableState } from '@tanstack/table-core';

export const searchParam = 'search';

export const initialState: InitialTableState = {
	pagination: {
		pageIndex: 0,
		pageSize: 25
	}
};

export const filters: Filter[] = [
	{
		type: FilterTypeEnum.STRING,
		accessorKey: 'username',
		placeholder: m.filter_username()
	},
	{
		type: FilterTypeEnum.STRING,
		accessorKey: 'firstName',
		placeholder: m.filter_firstname()
	},
	{
		type: FilterTypeEnum.STRING,
		accessorKey: 'familyName',
		placeholder: m.filter_lastname()
	},
	{
		type: FilterTypeEnum.SELECT,
		accessorKey: 'roles',
		values: enumToOptions(CourseUserRoleEnum, m.course_user_role_enum),
		emptyValue: m.no_filter(),
		placeholder: m.filter_role()
	}
];

export const columns: ColDef<CourseUserDTO>[] = [
	{
		accessorKey: 'row_index',
		header: 'ID',
		columnName: 'ID',
		cell: ({ row, table }) => {
			return (
				table.getState().pagination.pageIndex * table.getState().pagination.pageSize + row.index + 1
			);
		},
		enableHiding: false,
		size: 0
	},
	// {
	// 	id: 'select',
	// 	header: ({ table }) =>
	// 		renderComponent(Checkbox, {
	// 			checked: table.getIsAllPageRowsSelected(),
	// 			indeterminate: table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected(),
	// 			onCheckedChange: (value: boolean) => table.toggleAllPageRowsSelected(!!value),
	// 			'aria-label': m.select_all()
	// 		}),
	// 	cell: ({ row }) =>
	// 		renderComponent(Checkbox, {
	// 			checked: row.getIsSelected(),
	// 			onCheckedChange: (value: boolean) => row.toggleSelected(!!value),
	// 			'aria-label': m.select_row()
	// 		}),
	// 	enableSorting: false,
	// 	enableHiding: false
	// },

	{
		accessorKey: 'username',
		columnName: m.user_username(),
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.user_username(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'familyName',
		columnName: m.user_family_name(),
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.user_family_name(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'firstName',
		columnName: m.user_first_name(),
		header: m.user_first_name()
	},
	{
		accessorKey: 'degreeBefore',
		columnName: m.user_degree_before(),
		header: m.user_degree_before()
	},
	{
		accessorKey: 'degreeAfter',
		columnName: m.user_degree_after(),
		header: m.user_degree_after()
	},
	{
		accessorKey: 'email',
		columnName: m.user_email(),
		header: m.user_email()
	},
	{
		accessorKey: 'roles',
		columnName: m.user_roles(),
		header: m.user_roles(),
		cell: ({ row, column }) => {
			return renderComponent(DataTableRoles, {
				id: row.original.id,
				meta: column.columnDef.meta,
				roles: row.original.roles,
				showButtons: GlobalState.activeRole == CourseUserRoleEnum.ADMIN
			});
		},
		id: 'roles'
	},
	{
		accessorKey: 'studyForm',
		columnName: m.classes_studyform(),
		header: m.classes_studyform(),
		cell: ({ row }) => {
			return m.study_form_enum({ value: String(row.original.studyForm) });
		}
	},
	{
		header: m.actions(),
		columnName: m.actions(),
		cell: ({ row, column }) => {
			return renderComponent(DataTableActions, {
				id: row.original.id,
				meta: column.columnDef.meta
			});
		},
		enableHiding: false,
		id: 'actions'
	}
];

export const tableConfig = {
	columns,
	filters,
	initialState,
	searchParam
};
