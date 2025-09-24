import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, SortButton } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import { SemesterEnum, type CourseListItemDTO } from '$lib/api_types';
import { FilterTypeEnum, type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import { enumToOptions } from '$lib/utils';

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

export const columns: ColumnDef<CourseListItemDTO>[] = [
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
				name: m.course_name(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'shortname',
		header: m.course_shortname()
	},
	{
		accessorKey: 'year',
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
		header: m.semester(),
		cell: ({ row }) => {
			return m.semester_enum({ value: row.original.semester });
		}
	},
	{
		header: m.actions(),
		cell: ({ row }) => {
			return renderComponent(DataTableActions, { id: row.original.id });
		}
	}
];
