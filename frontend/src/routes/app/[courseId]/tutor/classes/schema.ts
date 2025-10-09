import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, SortButton } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import DataTableTutors from './data-table-tutors.svelte';
import type { ClassListItemDTO } from '$lib/api_types';
import { type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';

export const filters: Filter[] = [];

export const columns: (ColumnDef<ClassListItemDTO> & { uniqueId?: string })[] = [
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
				name: m.classes_name(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'tutors',
		header: m.classes_tutors(),
		cell: ({ row }) => {
			return renderComponent(DataTableTutors, { tutors: row.original.tutors });
		}
	},
	{
		accessorKey: 'room',
		header: m.classes_room()
	},
	{
		accessorKey: 'type',
		header: m.classes_type(),
		cell: ({ row }) => {
			return m.class_type_enum({ value: row.original.type });
		}
	},
	{
		accessorKey: 'studyForm',
		header: m.classes_studyform(),
		cell: ({ row }) => {
			return m.study_form_enum({ value: row.original.studyForm });
		}
	},
	{
		accessorKey: 'timeFrom',
		header: m.classes_timefromto(),
		cell: ({ row }) => {
			return `${row.original.timeFrom} - ${row.original.timeTo}`;
		}
	},
	{
		accessorKey: 'day',
		header: m.classes_day(),
		cell: ({ row }) => {
			return m.week_day_enum({ value: row.original.day });
		}
	},
	{
		accessorKey: 'weekParity',
		header: m.classes_weekparity(),
		cell: ({ row }) => {
			return m.week_parity_enum({ value: row.original.weekParity });
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
