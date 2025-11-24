import { renderComponent, SortButton, type ColDef } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import { SemesterEnum, type CourseListItemDTO } from '$lib/api_types';
import { FilterTypeEnum, type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import { enumToOptions } from '$lib/utils';
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
		placeholder: m.filter_name()
	},
	{
		type: FilterTypeEnum.STRING,
		accessorKey: 'shortname',
		placeholder: m.filter_shortname()
	},
	{
		type: FilterTypeEnum.STRING,
		accessorKey: 'year',
		placeholder: m.filter_academicyearstart()
	},
	{
		type: FilterTypeEnum.SELECT,
		accessorKey: 'semester',
		values: enumToOptions(SemesterEnum),
		emptyValue: 'All',
		placeholder: m.filter_semester()
	}
];

export const columns: ColDef<CourseListItemDTO>[] = [
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
		columnName: m.course_name(),
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.course_name(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'shortname',
		columnName: m.course_shortname(),
		header: m.course_shortname()
	},
	{
		accessorKey: 'year',
		columnName: m.academic_year(),
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.academic_year(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			}),
		cell: ({ row }) => {
			return `${row.original.year}/${row.original.year + 1}`;
		}
	},
	{
		accessorKey: 'semester',
		columnName: m.semester(),
		header: m.semester(),
		cell: ({ row }) => {
			return m.semester_enum({ value: row.original.semester });
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
	searchParam
};
