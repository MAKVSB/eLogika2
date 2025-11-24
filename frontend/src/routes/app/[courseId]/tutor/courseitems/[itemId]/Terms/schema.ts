import { renderComponent, SortButton, type ColDef } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import { type TermDTO } from '$lib/api_types';
import { type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import DataTableDateRange from '$lib/components/ui/data-table/data-table-date-range.svelte';
import type { InitialTableState } from '@tanstack/table-core';

export const searchParam = 'termsearch';

export const initialState: InitialTableState = {
	pagination: {
		pageIndex: 0,
		pageSize: 50
	}
};

export const filters: Filter[] = [];

export const columns: ColDef<TermDTO>[] = [
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
		accessorKey: 'name',
		columnName: m.term_name(),
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.term_name(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'activeFrom',
		columnName: m.term_active(),
		header: m.term_active(),
		cell: ({ row }) => {
			return renderComponent(DataTableDateRange, {
				start: row.original.activeFrom,
				end: row.original.activeTo
			});
		}
	},
	{
		accessorKey: 'signInFrom',
		columnName: m.term_signin(),
		header: m.term_signin(),
		cell: ({ row }) => {
			return renderComponent(DataTableDateRange, {
				start: row.original.signInFrom,
				end: row.original.signInTo
			});
		}
	},
	{
		accessorKey: 'signOutFrom',
		columnName: m.term_signout(),
		header: m.term_signout(),
		cell: ({ row }) => {
			return renderComponent(DataTableDateRange, {
				start: row.original.signOutFrom,
				end: row.original.signOutTo
			});
		}
	},
	{
		accessorKey: 'studentsMax',
		columnName: m.term_maximumstudents(),
		header: m.term_maximumstudents()
	},
	{
		accessorKey: 'studentsJoined',
		columnName: m.term_signed_students(),
		header: m.term_signed_students()
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
