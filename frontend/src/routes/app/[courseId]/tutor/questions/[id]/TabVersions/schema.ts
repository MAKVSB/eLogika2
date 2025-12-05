import type { InitialTableState } from '@tanstack/table-core';
import { renderComponent, SortButton, type ColDef } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import type { QuestionVersionDTO } from '$lib/api_types';
import { type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import DataTableDate from '$lib/components/ui/data-table/data-table-date.svelte';
import DataTableByUser from '$lib/components/ui/data-table/data-table-by-user.svelte';
import { DataTableActionMode } from '$lib/components/ui/data-table/data-table-component.svelte';

export const initialState: InitialTableState = {};

export const filters: Filter[] = [];
export const columns: ColDef<QuestionVersionDTO>[] = [
	{
		accessorKey: 'row_index',
		header: 'ID',
		columnName: 'ID',
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
		accessorKey: 'title',
		columnName: m.question_title(),
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.question_title(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'createdBy',
		header: m.created_by(),
		columnName: m.created_by(),
		cell: ({ row }) => {
			return renderComponent(DataTableByUser, {
				user: row.original.createdBy,
				time: row.original.createdAt
			});
		}
	},
	{
		accessorKey: 'modifiedAt',
		header: m.question_modified_at(),
		columnName: m.question_modified_at(),
		cell: ({ row }) => {
			return renderComponent(DataTableDate, {
				dateTime: row.original.updatedAt
			});
		}
	},
	{
		header: m.actions(),
		columnName: m.actions(),
		cell: ({ row, column }) => {
			return renderComponent(DataTableActions, {
				id: row.original.id,
				isArchive: row.original.isArchiveVersion,
				meta: column.columnDef.meta
			});
		},
		id: 'actions'
	}
];

export const tableConfig = {
	columns,
	filters,
	initialState,
	paginationMode: DataTableActionMode.DISABLED,
	sortingMode: DataTableActionMode.FRONTEND,
	filterMode: DataTableActionMode.FRONTEND,
	columnHiding: false
};
