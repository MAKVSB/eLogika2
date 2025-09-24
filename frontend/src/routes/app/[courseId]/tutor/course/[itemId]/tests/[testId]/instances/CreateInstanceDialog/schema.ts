import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, SortButton } from '$lib/components/ui/data-table/index.js';
import { CourseUserRoleEnum, type CourseUserDTO } from '$lib/api_types';
import { FilterTypeEnum, type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import { enumToOptions } from '$lib/utils';
import { Checkbox } from '$lib/components/ui/checkbox';
import DataTableActions from './data-table-actions.svelte';

export const filters: Filter[] = [];

export const columns: (ColumnDef<CourseUserDTO> & { uniqueId?: string })[] = [
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
				id: row.original.id,
				meta: column.columnDef.meta
			});
		},
		uniqueId: 'actions'
	}
];
