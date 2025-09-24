import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, SortButton } from '$lib/components/ui/data-table/index.js';
import DataTableActions from './data-table-actions.svelte';
import DataTableCreatedBy from './data-table-created-by.svelte';
import DataTableDateRange from './data-table-date-range.svelte';
import { Checkbox } from '$lib/components/ui/checkbox/index.js';
import type { TestInstanceListItemDTO, TestListItemDTO } from '$lib/api_types';
import { FilterTypeEnum, type Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';

export const filters: Filter[] = [];

export const columns: (ColumnDef<TestInstanceListItemDTO> & { uniqueId?: string })[] = [
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
		accessorKey: 'participant',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.user_username(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			}),
		cell: ({ row }) => {
			return `${row.original.participant.firstName} ${row.original.participant.familyName} (${row.original.participant.username})`;
		}
	},
	{
		accessorKey: 'state',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: 'State',
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			}),
		cell: ({ row }) => {
			return m.test_instance_state_enum({ value: row.original.state });
		}
	},
	{
		accessorKey: 'form',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: 'Form',
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'startedAt',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: 'Time',
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			}),
		cell: ({ row }) => {
			return renderComponent(DataTableDateRange, {
				start: row.original.startedAt,
				end: row.original.endedAt
			});
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
