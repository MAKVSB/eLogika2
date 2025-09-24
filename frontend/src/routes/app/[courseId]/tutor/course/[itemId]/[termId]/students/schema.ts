import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, SortButton } from '$lib/components/ui/data-table/index.js';
import { QuestionCheckedByFilterEnum, QuestionTypeEnum } from '$lib/api_types';
import type { JoinedStudentDTO, QuestionListItemDTO } from '$lib/api_types';
import { FilterTypeEnum, type Filter } from '$lib/components/ui/data-table/filter';
import DataTableCheck from '$lib/components/ui/data-table/data-table-check.svelte';
import { m } from '$lib/paraglide/messages';
import { enumToOptions } from '$lib/utils';
import DataTableActions from './data-table-actions.svelte';

export const filters: Filter[] = [];

export const columns: (ColumnDef<JoinedStudentDTO> & { uniqueId?: string })[] = [
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
		header: m.user_username()
	},
	{
		accessorKey: 'firstName',
		header: m.user_first_name()
	},
	{
		accessorKey: 'familyName',
		header: m.user_family_name()
	},
	{
		accessorKey: 'email',
		header: m.user_email()
	},
	{
		header: m.actions(),
		cell: ({ row, column }) => {
			return renderComponent(DataTableActions, {
				userId: row.original.id
			});
		},
		uniqueId: 'actions'
	}
];
