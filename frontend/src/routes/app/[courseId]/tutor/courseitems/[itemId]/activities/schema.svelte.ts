import { renderComponent, SortButton, type ColDef } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import type { ActivityListItemDTO } from '$lib/api_types';
import { type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import { displayUserName } from '$lib/utils';
import type { InitialTableState } from '@tanstack/table-core';

export const searchParam = 'search';

export const initialState: InitialTableState = {
	pagination: {
		pageIndex: 0,
		pageSize: 25
	}
};

export const filters: Filter[] = [];

export const columns: ColDef<ActivityListItemDTO>[] = [
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
		accessorKey: 'participant',
		columnName: m.created_by(),
		header: m.created_by(),
		cell: ({ row }) => {
			return `${displayUserName(row.original.participant)} (${row.original.participant.username})`;
		}
	},
	{
		accessorKey: 'termId',
		columnName: m.activity_term(),
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.activity_term(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
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

export const tableConfig = $state({
	columns,
	filters,
	initialState,
	searchParam
});
