import type { InitialTableState } from '@tanstack/table-core';
import { renderComponent, SortButton, type ColDef } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import DataTableCheckedBy from './data-table-checked-by.svelte';
import DataTableCreatedBy from '$lib/components/ui/data-table/data-table-created-by.svelte';
import { QuestionCheckedByFilterEnum, QuestionTypeEnum } from '$lib/api_types';
import type { QuestionListItemDTO } from '$lib/api_types';
import { FilterTypeEnum, type Filter } from '$lib/components/ui/data-table/filter';
import DataTableCheck from '$lib/components/ui/data-table/data-table-check.svelte';
import { m } from '$lib/paraglide/messages';
import { enumToOptions } from '$lib/utils';

export const searchParam = 'search';

export const initialState: InitialTableState = {
	pagination: {
		pageIndex: 0,
		pageSize: 25
	}
};

export const filters: Filter[] = $state([
	{
		type: FilterTypeEnum.STRING,
		accessorKey: 'title',
		placeholder: m.filter_title()
	},
	{
		type: FilterTypeEnum.SELECT,
		accessorKey: 'chapterId',
		values: [],
		emptyValue: 'No filter',
		placeholder: m.filter_chapter()
	},
	{
		type: FilterTypeEnum.SELECT,
		accessorKey: 'categoryId',
		values: [],
		emptyValue: 'No filter',
		placeholder: m.filter_category()
	},
	{
		type: FilterTypeEnum.SELECT,
		accessorKey: 'questionType',
		values: enumToOptions(QuestionTypeEnum, m.question_type_enum),
		emptyValue: 'No filter',
		placeholder: m.filter_questiontype()
	},
	{
		type: FilterTypeEnum.SELECT,
		accessorKey: 'checkedBy',
		values: enumToOptions(QuestionCheckedByFilterEnum),
		emptyValue: 'No filter',
		placeholder: m.filter_checkedstate()
	}
]);

export const columns: ColDef<QuestionListItemDTO>[] = [
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
		columnName: m.question_title(),
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.question_title(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'chapterId',
		columnName: m.question_chapter(),
		header: m.question_chapter(),
		cell: ({ row }) => {
			return row.original.chapterName;
		}
	},
	{
		accessorKey: 'categoryId',
		columnName: m.question_category(),
		header: m.question_category(),
		cell: ({ row }) => {
			return row.original.categoryName;
		}
	},
	{
		accessorKey: 'questionType',
		columnName: m.question_type(),
		header: m.question_type(),
		cell: ({ row }) => {
			return m.question_type_enum({ value: row.original.questionType });
		}
	},
	{
		accessorKey: 'questionFormat',
		columnName: m.question_format(),
		header: m.question_format(),
		cell: ({ row }) => {
			return m.question_format_enum({ value: row.original.questionFormat });
		}
	},
	{
		accessorKey: 'active',
		columnName: m.active(),
		header: m.active(),
		cell: ({ row, column }) => {
			return renderComponent(DataTableCheck, {
				checked: row.original.active,
				meta: column.columnDef.meta,
				id: row.original.id
			});
		},
		id: 'active'
	},
	{
		accessorKey: 'checkedBy',
		columnName: m.question_checked_by(),
		header: m.question_checked_by(),
		cell: ({ row }) => {
			return renderComponent(DataTableCheckedBy, { users: row.original.checkedBy });
		}
	},
	{
		accessorKey: 'createdBy',
		columnName: m.question_created_by(),
		header: m.question_created_by(),
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
				checkedBy: row.original.checkedBy,
				meta: column.columnDef.meta
			});
		},
		enableHiding: false,
		id: 'actions'
	}
];

export const tableConfig = $state({
	columns,
	filters,
	initialState,
	searchParam
});
