import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, SortButton } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import { Checkbox } from '$lib/components/ui/checkbox/index.js';
import { type CourseItemDTO } from '$lib/api_types';
import { type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import DataTableCheck from '$lib/components/ui/data-table/data-table-check.svelte';

export const filters: Filter[] = [];

export const columns: (ColumnDef<CourseItemDTO> & { uniqueId?: string })[] = [
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
				name: m.course_item_name(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'type',
		header: m.course_item_type(),
		cell: ({ row }) => {
			return m.course_item_type_enum({ value: row.original.type });
		}
	},
	{
		accessorKey: 'pointsMin',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: 'Points',
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			}),
		cell: ({ row }) => {
			return `${row.original.pointsMin}/${row.original.pointsMax}`;
		}
	},
	{
		accessorKey: 'mandatory',
		header: m.course_item_mandatory(),
		cell: ({ row, column }) => {
			return renderComponent(DataTableCheck, {
				checked: row.original.mandatory
			});
		}
	},
	{
		accessorKey: 'managedBy',
		header: m.course_item_managedby(),
		cell: ({ row, column }) => {
			return m.course_user_role_enum({ value: row.original.managedBy });
		}
	},
	{
		header: m.actions(),
		cell: ({ row, column }) => {
			return renderComponent(DataTableActions, {
				id: row.original.id,
				meta: column.columnDef.meta,
				type: row.original.type,
				editable: row.original.editable
			});
		},
		uniqueId: 'actions'
	}
];
