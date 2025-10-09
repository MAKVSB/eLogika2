import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, SortButton } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import { Checkbox } from '$lib/components/ui/checkbox/index.js';
import { type TermDTO } from '$lib/api_types';
import { type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import DateRange from '$lib/components/date-range.svelte';

export const filters: Filter[] = [];

// studentsMax: number;

export const columns: (ColumnDef<TermDTO> & { uniqueId?: string })[] = [
	{
		accessorKey: 'id',
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
		accessorKey: 'name',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.term_name(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'activeFrom',
		header: m.term_active(),
		cell: ({ row }) => {
			return renderComponent(DateRange, {
				start: row.original.activeFrom,
				end: row.original.activeTo
			});
		}
	},
	{
		accessorKey: 'signInFrom',
		header: m.term_signin(),
		cell: ({ row }) => {
			return renderComponent(DateRange, {
				start: row.original.signInFrom,
				end: row.original.signInTo
			});
		}
	},
	{
		accessorKey: 'signOutFrom',
		header: m.term_signout(),
		cell: ({ row }) => {
			return renderComponent(DateRange, {
				start: row.original.signOutFrom,
				end: row.original.signOutTo
			});
		}
	},
	{
		accessorKey: 'studentsMax',
		header: m.term_maximumstudents()
	},
	{
		accessorKey: 'studentsJoined',
		header: m.term_signed_students()
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
