import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, SortButton } from '$lib/components/ui/data-table/index.js';
import { QuestionCheckedByFilterEnum, QuestionTypeEnum } from '$lib/api_types';
import type { JoinedStudentDTO, QuestionListItemDTO } from '$lib/api_types';
import { FilterTypeEnum, type Filter } from '$lib/components/ui/data-table/filter';
import DataTableCheck from '$lib/components/ui/data-table/data-table-check.svelte';
import { m } from '$lib/paraglide/messages';
import { enumToOptions } from '$lib/utils';
import DataTableActions from './data-table-actions.svelte';
import DataTableDate from '$lib/components/ui/data-table/data-table-date.svelte';

export const filters: Filter[] = [];

export const columns: (ColumnDef<JoinedStudentDTO> & { uniqueId?: string })[] = [
	{
		accessorKey: 'row_index',
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
		accessorKey: 'username',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.user_username(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'familyName',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.user_family_name(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'firstName',
		header: m.user_first_name()
	},
	{
		accessorKey: 'degreeBefore',
		header: m.user_degree_before()
	},
	{
		accessorKey: 'degreeAfter',
		header: m.user_degree_after()
	},
	{
		accessorKey: 'createdAt',
		header: m.course_item_term_user_signedinat(),
		cell: ({ row }) => {
			return renderComponent(DataTableDate, {
				dateTime: row.original.createdAt
			});
		}
	},
	{
		accessorKey: 'deletedAt',
		header: m.course_item_term_user_signedoutat(),
		cell: ({ row }) => {
			return renderComponent(DataTableDate, {
				dateTime: row.original.deletedAt
			});
		}
	},
	{
		accessorKey: 'email',
		header: m.user_email()
	},
	{
		header: m.actions(),
		cell: ({ row, column }) => {
			return renderComponent(DataTableActions, {
				userId: row.original.userId,
				isJoined: row.original.deletedAt == null
			});
		},
		uniqueId: 'actions'
	}
];
