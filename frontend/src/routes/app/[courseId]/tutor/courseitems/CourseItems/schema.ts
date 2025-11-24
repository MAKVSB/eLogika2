import { renderComponent, type ColDef } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import { type CourseItemDTO } from '$lib/api_types';
import { type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import DataTableCheck from '$lib/components/ui/data-table/data-table-check.svelte';
import type { InitialTableState } from '@tanstack/table-core';
import { DataTableActionMode } from '$lib/components/ui/data-table/data-table-component.svelte';

export const initialState: InitialTableState = {};

export const filters: Filter[] = [];

export const columns: ColDef<CourseItemDTO>[] = [
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
		columnName: m.course_item_name(),
		header: m.course_item_name()
	},
	{
		accessorKey: 'type',
		columnName: m.course_item_type(),
		header: m.course_item_type(),
		cell: ({ row }) => {
			return m.course_item_type_enum({ value: row.original.type });
		}
	},
	{
		accessorKey: 'pointsMin',
		columnName: m.course_item_points_min_max(),
		header: m.course_item_points_min_max(),
		cell: ({ row }) => {
			return `${row.original.pointsMin}/${row.original.pointsMax}`;
		}
	},
	{
		accessorKey: 'mandatory',
		columnName: m.course_item_mandatory(),
		header: m.course_item_mandatory(),
		cell: ({ row, column }) => {
			return renderComponent(DataTableCheck, {
				checked: row.original.mandatory
			});
		}
	},
	{
		accessorKey: 'managedBy',
		columnName: m.course_item_managedby(),
		header: m.course_item_managedby(),
		cell: ({ row, column }) => {
			return m.course_user_role_enum({ value: row.original.managedBy });
		}
	},
	{
		header: m.actions(),
		columnName: m.actions(),
		cell: ({ row, column }) => {
			return renderComponent(DataTableActions, {
				id: row.original.id,
				meta: column.columnDef.meta,
				type: row.original.type,
				editable: row.original.editable
			});
		},
		enableHiding: false,
		id: 'actions'
	}
];

export const tableConfig = {
	columns,
	filters,
	paginationMode: DataTableActionMode.DISABLED,
	sortingMode: DataTableActionMode.FRONTEND,
	filterMode: DataTableActionMode.FRONTEND
};
