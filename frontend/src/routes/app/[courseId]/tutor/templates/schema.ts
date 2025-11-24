import { renderComponent, SortButton, type ColDef } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import DataTableCreatedBy from '$lib/components/ui/data-table/data-table-created-by.svelte';
import { type Filter } from '$lib/components/ui/data-table/filter';
import type { TemplateListItemDTO } from '$lib/api_types';
import { m } from '$lib/paraglide/messages';

export const searchParam = 'search';

export const initialState = {};

export const filters: Filter[] = [];

export const columns: ColDef<TemplateListItemDTO>[] = [
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
		accessorKey: 'title',
		columnName: m.template_title(),
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.template_title(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'description',
		columnName: m.template_description(),
		header: m.template_description()
	},
	{
		accessorKey: 'blocks',
		columnName: m.template_block_count(),
		header: m.template_block_count()
	},
	{
		accessorKey: 'createdBy',
		columnName: m.template_created_by(),
		header: m.template_created_by(),
		cell: ({ row }) => {
			return renderComponent(DataTableCreatedBy, {
				createdBy: row.original.createdBy,
				createdAt: row.original.createdAt
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
