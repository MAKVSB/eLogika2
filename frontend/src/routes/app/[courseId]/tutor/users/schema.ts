import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, SortButton } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import { CourseUserRoleEnum, type CourseUserDTO } from '$lib/api_types';
import { FilterTypeEnum, type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import { enumToOptions } from '$lib/utils';

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
		emptyValue: 'All',
		placeholder: m.filter_role()
	}
];

export const columns: (ColumnDef<CourseUserDTO> & { uniqueId?: string })[] = [
	{
		accessorKey: 'row_index',
		header: 'ID',
		cell: ({ row, table }) => {
			return (
				table.getState().pagination.pageIndex * table.getState().pagination.pageSize + row.index + 1
			);
		},
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
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.user_username(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'firstName',
		header: m.user_first_name()
	},
	{
		accessorKey: 'familyName',
		header: m.user_family_name()
	},
	{
		accessorKey: 'degreeBefore',
		header: m.user_degree_before()
	},
	{
		accessorKey: 'degreeAfter',
		header: m.user_degree_after()
	},
	{
		accessorKey: 'email',
		header: m.user_email()
	},
	{
		accessorKey: 'roles',
		header: m.user_roles(),
		cell: ({ row }) => {
			return row.original.roles.map((role) => m.course_user_role_enum({ value: role })).join(', ');
		}
	},
	{
		accessorKey: 'studyForm',
		header: m.classes_studyform(),
		cell: ({ row }) => {
			return m.study_form_enum({ value: String(row.original.studyForm) });
		}
	},
	{
		header: m.actions(),
		cell: ({ row, column }) => {
			return renderComponent(DataTableActions, {
				id: row.original.id,
				meta: column.columnDef.meta
			});
		},
		uniqueId: 'actions'
	}
];
