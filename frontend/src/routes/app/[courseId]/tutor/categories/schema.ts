import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, SortButton } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import { Checkbox } from '$lib/components/ui/checkbox/index.js';
import { type Filter } from '$lib/components/ui/data-table/filter';
import type { CategoryListItemDTO } from '$lib/api_types';
import { m } from '$lib/paraglide/messages';

export const filters: Filter[] = [];

export const columns: (ColumnDef<CategoryListItemDTO> & { uniqueId?: string })[] = [
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
				name: m.category_name(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'chapterName',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.category_chapter(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'stepsCount',
		header: m.category_steps_count()
	},
	{
		header: m.actions(),
		cell: ({ row }) => {
			return renderComponent(DataTableActions, {
				id: row.original.id
			});
		},
		uniqueId: 'actions'
	}
];
