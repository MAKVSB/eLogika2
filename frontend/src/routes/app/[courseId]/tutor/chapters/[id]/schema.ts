import { renderComponent, type ColDef } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import DataTableOrder from './data-table-order.svelte';
import DataTableCheck from '$lib/components/ui/data-table/data-table-check.svelte';
import type { Filter } from '$lib/components/ui/data-table/filter';
import type { ChapterListItemDTO } from '$lib/api_types';
import { m } from '$lib/paraglide/messages';
import type { InitialTableState } from '@tanstack/table-core';
import { DataTableActionMode } from '$lib/components/ui/data-table/data-table-component.svelte';

export const searchParam = 'subchapters';

export const initialState: InitialTableState = {
	pagination: {
		pageIndex: 0,
		pageSize: 25
	}
};

export const filters: Filter[] = [];

export const columns: ColDef<ChapterListItemDTO>[] = [
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
		accessorKey: 'order',
		header: m.chapter_order(),
		columnName: m.chapter_order(),
		cell: ({ row, column, table }) => {
			const rows = table.getRowModel().rows;
			const index = rows.findIndex((r) => r.id === row.id);
			return renderComponent(DataTableOrder, {
				id: row.original.id,
				order: row.original.order,
				meta: column.columnDef.meta,
				isFirst: index === 0,
				isLast: index === rows.length - 1
			});
		},
		maxSize: 1
	},
	{
		accessorKey: 'name',
		columnName: m.chapter_name(),
		header: m.chapter_name()
	},
	{
		accessorKey: 'visible',
		columnName: m.chapter_visible(),
		header: m.chapter_visible(),
		cell: ({ row }) => {
			return renderComponent(DataTableCheck, {
				checked: row.original.visible
			});
		}
	},
	{
		header: m.actions(),
		columnName: m.actions(),
		cell: ({ row }) => {
			return renderComponent(DataTableActions, { id: row.original.id });
		},
		enableHiding: false,
		id: 'actions'
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
