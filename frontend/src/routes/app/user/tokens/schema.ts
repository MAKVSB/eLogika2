import { renderComponent, type ColDef } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import { type Filter } from '$lib/components/ui/data-table/filter';
import type { UserAPiTokenDTO } from '$lib/api_types';
import { m } from '$lib/paraglide/messages';
import type { InitialTableState } from '@tanstack/table-core';
import { DataTableActionMode } from '$lib/components/ui/data-table/data-table-component.svelte';

export const searchParam = 'search';

export const initialState: InitialTableState = {};

export const filters: Filter[] = [];

export const columns: ColDef<UserAPiTokenDTO>[] = [
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
		accessorKey: 'name',
		columnName: m.token_name(),
		header: m.token_name()
	},
	{
		accessorKey: 'issuedAt',
		columnName: m.token_issued_at(),
		header: m.token_issued_at(),
		cell: ({ row }) => {
			return new Date(row.original.issuedAt).toLocaleString();
		}
	},
	{
		accessorKey: 'expiresAt',
		columnName: m.token_expires_at(),
		header: m.token_expires_at(),
		cell: ({ row }) => {
			return new Date(row.original.expiresAt).toLocaleString();
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
	searchParam,
	paginationMode: DataTableActionMode.DISABLED,
	sortingMode: DataTableActionMode.FRONTEND,
	filterMode: DataTableActionMode.FRONTEND
};
