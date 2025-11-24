import SortButton from './data-table-sort-button.svelte';
import DataTable from './data-table-component.svelte';
import type { ColumnDef, RowData } from '@tanstack/table-core';

export { default as FlexRender } from './flex-render.svelte';
export { renderComponent, renderSnippet } from './render-helpers.js';
export { createSvelteTable } from './data-table.svelte';

export { SortButton, DataTable };

export type ColDef<T extends RowData, V = unknown> = ColumnDef<T, V> & {
	accessorKey?: string;
	columnName: string;
};
