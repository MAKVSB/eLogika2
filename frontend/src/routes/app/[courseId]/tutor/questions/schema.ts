import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, SortButton } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import DataTableCheckedBy from './data-table-checked-by.svelte';
import DataTableCreatedBy from './data-table-created-by.svelte';
import { QuestionCheckedByFilterEnum, QuestionTypeEnum } from '$lib/api_types';
import type { QuestionListItemDTO } from '$lib/api_types';
import { FilterTypeEnum, type Filter } from '$lib/components/ui/data-table/filter';
import DataTableCheck from '$lib/components/ui/data-table/data-table-check.svelte';
import { m } from '$lib/paraglide/messages';
import { enumToOptions } from '$lib/utils';

export const filters: Filter[] = [
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
		values: enumToOptions(QuestionTypeEnum),
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
];

export const columns: (ColumnDef<QuestionListItemDTO> & { uniqueId?: string })[] = [
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
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: 'Title',
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'chapterId'
	},
	{
		accessorKey: 'chapterName',
		header: 'Chapter name'
	},
	{
		accessorKey: 'categoryId'
	},
	{
		accessorKey: 'categoryName',
		header: 'Category name'
	},
	{
		accessorKey: 'questionType',
		header: m.question_type(),
		cell: ({ row }) => {
			return m.question_type_enum({ value: row.original.questionType });
		}
	},
	{
		accessorKey: 'questionFormat',
		header: m.question_format(),
		cell: ({ row }) => {
			return m.question_format_enum({ value: row.original.questionFormat });
		}
	},
	{
		accessorKey: 'active',
		header: m.active(),
		cell: ({ row, column }) => {
			return renderComponent(DataTableCheck, {
				checked: row.original.active,
				meta: column.columnDef.meta,
				id: row.original.id
			});
		},
		uniqueId: 'active'
	},
	{
		accessorKey: 'checkedBy',
		header: m.question_checked_by(),
		cell: ({ row }) => {
			return renderComponent(DataTableCheckedBy, { users: row.original.checkedBy });
		}
	},
	{
		accessorKey: 'createdBy',
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
		cell: ({ row, column }) => {
			return renderComponent(DataTableActions, {
				id: row.original.id,
				checkedBy: row.original.checkedBy,
				meta: column.columnDef.meta
			});
		},
		uniqueId: 'actions'
	}
];
