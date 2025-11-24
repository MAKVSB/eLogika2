import { renderComponent, SortButton, type ColDef } from '$lib/components/ui/data-table/index.js';
import type { TestInstanceEventDTO } from '$lib/api_types';
import { getLocale } from '$lib/paraglide/runtime';
import type { Filter } from '$lib/components/ui/data-table/filter';
import { m } from '$lib/paraglide/messages';
import { displayUserName } from '$lib/utils';
import type { InitialTableState } from '@tanstack/table-core';

export const searchParam = 'eventSearch';

export const initialState: InitialTableState = {
	pagination: {
		pageIndex: 0,
		pageSize: 5
	}
};

export const filters: Filter[] = [];

export const columns: ColDef<TestInstanceEventDTO>[] = [
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
		accessorKey: 'id',
		columnName: m.instanceevent_eventid(),
		header: m.instanceevent_eventid()
	},
	{
		accessorKey: 'user',
		columnName: m.instanceevent_user(),
		header: m.instanceevent_user(),
		cell: ({ row }) => {
			return `${displayUserName(row.original.user)} (${row.original.user.username})`;
		}
	},
	{
		accessorKey: 'occuredAt',
		columnName: m.instanceevent_occuredat(),
		header: m.instanceevent_occuredat(),
		cell: ({ row }) => {
			return new Date(row.original.occuredAt).toLocaleString(getLocale());
		}
	},
	{
		accessorKey: 'recievedAt',
		columnName: m.instanceevent_receivedat(),
		header: m.instanceevent_receivedat(),
		cell: ({ row }) => {
			return new Date(row.original.receivedAt).toLocaleString(getLocale());
		}
	},
	{
		accessorKey: 'eventType',
		columnName: m.instanceevent_eventtype(),
		header: m.instanceevent_eventtype(),
		cell: ({ row }) => {
			return m.test_instance_event_type_enum({ value: row.original.eventType });
		}
	},
	{
		accessorKey: 'eventData',
		columnName: m.instanceevent_data(),
		header: m.instanceevent_data(),
		cell: ({ row }) => {
			return JSON.stringify(row.original.eventData);
		}
	},
	{
		accessorKey: 'eventSource',
		columnName: m.instanceevent_source(),
		header: m.instanceevent_source()
	},
	{
		accessorKey: 'pageId',
		columnName: m.instanceevent_pageid(),
		header: m.instanceevent_pageid()
	}
];

export const tableConfig = {
	columns,
	filters,
	initialState,
	searchParam
};
