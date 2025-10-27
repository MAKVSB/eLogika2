import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import { type Filter } from '$lib/components/ui/data-table/filter';
import type { UserAPiTokenDTO } from '$lib/api_types';
import { m } from '$lib/paraglide/messages';

export const filters: Filter[] = [];

export const columns: (ColumnDef<UserAPiTokenDTO> & { uniqueId?: string })[] = [
	{
		accessorKey: 'row_index',
		header: 'ID',
		cell: ({ row, table }) => {
			return (
				table.getState().pagination.pageIndex * table.getState().pagination.pageSize + row.index + 1
			);
		},
		size: 0
	},
	{
		accessorKey: 'name',
		header: 'Name'
	},
	{
		accessorKey: 'issuedAt',
		header: 'Issued at',
		cell: ({ row }) => {
			return new Date(row.original.issuedAt).toLocaleString();
		}
	},
	{
		accessorKey: 'expiresAt',
		header: 'Expires at',
		cell: ({ row }) => {
			return new Date(row.original.expiresAt).toLocaleString();
		}
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
