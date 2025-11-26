import { renderComponent, SortButton, type ColDef } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import DataTableByUser from '$lib/components/ui/data-table/data-table-by-user.svelte';
import type { TestListItemDTO } from '$lib/api_types';
import { FilterTypeEnum, type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
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
		type: FilterTypeEnum.SELECT,
		accessorKey: 'termId',
		values: [], //Fill from page
		placeholder: m.filter_term(),
		emptyValue: m.no_filter()
	}
];

export const columns: ColDef<TestListItemDTO>[] = [
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
		columnName: m.test_name(),
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.test_name(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'group',
		columnName: m.test_group(),
		header: m.test_group()
	},
	{
		accessorKey: 'termId',
		columnName: m.test_term(),
		header: m.test_term(),
		cell: ({ row }) => {
			return row.original.term;
		}
	},
	{
		accessorKey: 'createdBy',
		columnName: m.test_createdby(),
		header: m.test_createdby(),
		cell: ({ row }) => {
			return renderComponent(DataTableByUser, {
				user: row.original.createdBy,
				time: row.original.createdAt
			});
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
