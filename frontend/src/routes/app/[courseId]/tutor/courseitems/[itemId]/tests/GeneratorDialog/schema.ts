import { renderComponent, SortButton, type ColDef } from '$lib/components/ui/data-table/index.js';
import { type JoinedStudentDTO } from '$lib/api_types';
import { type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import { Checkbox } from '$lib/components/ui/checkbox';
import DataTableDate from '$lib/components/ui/data-table/data-table-date.svelte';
import type { InitialTableState } from '@tanstack/table-core';
import { DataTableActionMode } from '$lib/components/ui/data-table/data-table-component.svelte';

export const searchParam = 'generateTestSearch';

export const initialState: InitialTableState = {};

export const filters: Filter[] = [];

export const columns: ColDef<JoinedStudentDTO>[] = [
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
		id: 'select',
		columnName: 'select',
		header: ({ table }) =>
			renderComponent(Checkbox, {
				checked: table.getIsAllPageRowsSelected(),
				indeterminate: table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected(),
				onCheckedChange: (value: boolean) => table.toggleAllPageRowsSelected(!!value),
				'aria-label': m.select_all()
			}),
		cell: ({ row }) =>
			renderComponent(Checkbox, {
				checked: row.getIsSelected(),
				onCheckedChange: (value: boolean) => row.toggleSelected(!!value),
				'aria-label': m.select_row()
			}),
		enableSorting: false,
		enableHiding: false
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
		accessorKey: 'createdAt',
		columnName: m.course_item_term_user_signedinat(),
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.course_item_term_user_signedinat(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			}),
		cell: ({ row }) => {
			return renderComponent(DataTableDate, {
				dateTime: row.original.createdAt
			});
		}
	}
];

export const tableConfig = {
	columns,
	filters,
	initialState,
	searchParam,
	paginationMode: DataTableActionMode.DISABLED,
	sortingMode: DataTableActionMode.FRONTEND,
	filterMode: DataTableActionMode.FRONTEND
};
