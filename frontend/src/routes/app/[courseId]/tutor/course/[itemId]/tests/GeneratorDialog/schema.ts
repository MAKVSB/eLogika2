import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, SortButton } from '$lib/components/ui/data-table/index.js';
import { CourseUserRoleEnum, type CourseUserDTO } from '$lib/api_types';
import { FilterTypeEnum, type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import { enumToOptions } from '$lib/utils';
import { Checkbox } from '$lib/components/ui/checkbox';

export const filters: Filter[] = [];

export const columns: (ColumnDef<CourseUserDTO> & { uniqueId?: string })[] = [
	{
		id: 'select',
		header: ({ table }) =>
			renderComponent(Checkbox, {
				checked: table.getIsAllPageRowsSelected(),
				indeterminate: table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected(),
				onCheckedChange: (value: boolean) => table.toggleAllPageRowsSelected(!!value),
				'aria-label': m.select_all()
			}),
		cell: ({ row }) =>
			renderComponent(Checkbox, {
				checked: row.getIsSelected(),
				onCheckedChange: (value: boolean) => row.toggleSelected(!!value),
				'aria-label': m.select_row()
			}),
		enableSorting: false,
		enableHiding: false
	},
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
		accessorKey: 'degreeBefore',
		header: m.user_degree_before()
	},
	{
		accessorKey: 'degreeAfter',
		header: m.user_degree_after()
	},
	{
		accessorKey: 'email',
		header: m.user_email()
	}
];
