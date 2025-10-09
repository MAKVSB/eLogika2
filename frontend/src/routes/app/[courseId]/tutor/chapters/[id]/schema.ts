import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import DataTableOrder from './data-table-order.svelte';
import DataTableCheck from '$lib/components/ui/data-table/data-table-check.svelte';
import type { Filter } from '$lib/components/ui/data-table/filter';
import type { ChapterListItemDTO } from '$lib/api_types';
import { m } from '$lib/paraglide/messages';

export const filters: Filter[] = [];

export const columns: (ColumnDef<ChapterListItemDTO> & { uniqueId?: string })[] = [
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
	{
		accessorKey: 'order',
		header: 'Order',
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
		header: m.chapter_name()
	},
	{
		accessorKey: 'visible',
		header: m.chapter_visible(),
		cell: ({ row }) => {
			return renderComponent(DataTableCheck, {
				checked: row.original.visible
			});
		}
	},
	{
		header: m.actions(),
		cell: ({ row }) => {
			return renderComponent(DataTableActions, { id: row.original.id });
		},
		uniqueId: 'actions'
	}
];
