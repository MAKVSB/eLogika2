import { renderComponent, SortButton, type ColDef } from '$lib/components/ui/data-table/index.js';
import { type CourseUserDTO } from '$lib/api_types';
import { type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import DataTableActions from './data-table-actions.svelte';
import type { InitialTableState } from '@tanstack/table-core';

export const searchParam = 'createInstanceSearch';

export const initialState: InitialTableState = {
	pagination: {
		pageIndex: 0,
		pageSize: 25
	}
};

export const filters: Filter[] = [];

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
