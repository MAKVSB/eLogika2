import { renderComponent, SortButton, type ColDef } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import DataTableCreatedBy from '$lib/components/ui/data-table/data-table-created-by.svelte';
import type { SupportTicketListItemDTO } from '$lib/api_types';
import { FilterTypeEnum, type Filter } from '$lib/components/ui/data-table/filter';
import DataTableCheck from '$lib/components/ui/data-table/data-table-check.svelte';
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
		type: FilterTypeEnum.STRING,
		accessorKey: 'name',
		placeholder: m.filter_title()
	},
	{
		type: FilterTypeEnum.SELECT,
		accessorKey: 'solved',
		values: [
			{
				value: '1',
				display: m.yes_no({ value: 'true' })
			},
			{
				value: '0',
				display: m.yes_no({ value: 'false' })
			}
		],
		emptyValue: 'No filter',
		placeholder: m.filter_checkedstate()
	}
];

export const columns: ColDef<SupportTicketListItemDTO>[] = [
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
		columnName: m.ticket_title(),
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.ticket_title(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'createdBy',
		columnName: m.ticket_created_by(),
		header: m.ticket_created_by(),
		cell: ({ row }) => {
			return renderComponent(DataTableCreatedBy, {
				createdBy: row.original.createdBy,
				createdAt: row.original.createdAt
			});
		}
	},
	{
		accessorKey: 'solved',
		columnName: m.ticket_solved(),
		header: m.ticket_solved(),
		cell: ({ row, column }) => {
			return renderComponent(DataTableCheck, {
				checked: row.original.solved
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
