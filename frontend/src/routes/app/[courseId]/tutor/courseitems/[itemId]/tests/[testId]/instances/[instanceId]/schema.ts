import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent, SortButton } from '$lib/components/ui/data-table/index.js';
import type { TestInstanceEventDTO } from '$lib/api_types';
import { getLocale } from '$lib/paraglide/runtime';
import type { Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import { displayUserName } from '$lib/utils';

export const filters: Filter[] = [];

export const columns: (ColumnDef<TestInstanceEventDTO> & { uniqueId?: string })[] = [
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
		accessorKey: 'id',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: 'Id',
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'user',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: m.user_username(),
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			}),
		cell: ({ row }) => {
			return `${displayUserName(row.original.user)} (${row.original.user.username})`;
		}
	},
	{
		accessorKey: 'occuredAt',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: 'Occured at',
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			}),
		cell: ({ row }) => {
			return new Date(row.original.occuredAt).toLocaleString(getLocale());
		}
	},
	{
		accessorKey: 'eventType',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: 'Event type',
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			}),
		cell: ({ row }) => {
			return m.test_instance_event_type_enum({ value: row.original.eventType });
		}
	},
	{
		accessorKey: 'eventData',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: 'Event data',
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			}),
		cell: ({ row }) => {
			return JSON.stringify(row.original.eventData);
		}
	},
	{
		accessorKey: 'eventSource',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: 'Event source',
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'pageId',
		header: ({ column }) =>
			renderComponent(SortButton, {
				name: 'Page id',
				sorted: column.getIsSorted(),
				onclick: column.getToggleSortingHandler()
			})
	}
];
