import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, SortButton } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import { Checkbox } from '$lib/components/ui/checkbox/index.js';
import { FilterTypeEnum, type Filter } from '$lib/components/ui/data-table/filter';
import type { UserListItemDTO } from '$lib/api_types';
import { m } from '$lib/paraglide/messages';

export const filters: Filter[] = [
	{
		type: FilterTypeEnum.STRING,
		accessorKey: 'username',
		placeholder: m.filter_username()
	},
	{
		type: FilterTypeEnum.STRING,
		accessorKey: 'firstName',
		placeholder: m.filter_firstname()
	},
	{
		type: FilterTypeEnum.STRING,
		accessorKey: 'familyName',
		placeholder: m.filter_lastname()
	}
];

export const columns: ColumnDef<UserListItemDTO>[] = [
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
		accessorKey: 'firstName',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.user_first_name(),
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
		accessorKey: 'degreeBefore',
		header: m.user_degree_before()
	},
	{
		accessorKey: 'degreeAfter',
		header: m.user_degree_after()
	},
	{
		accessorKey: 'email',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.user_email(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		header: m.actions(),
		cell: ({ row }) => {
			return renderComponent(DataTableActions, { id: row.original.id });
		}
	}
];
