import type {
	ColumnFiltersState,
	InitialTableState,
	PaginationState,
	SortingState,
	TableState
} from '@tanstack/table-core';
import { encodeJsonToBase64Url } from './services/api.svelte';

export class DataTableSearchParams {
	pagination?: Partial<PaginationState>;
	sorting?: SortingState;
	columnFilters?: ColumnFiltersState;

	constructor(init?: Partial<DataTableSearchParams>) {
		Object.assign(this, init);
	}

	/** Convert class → encoded URL string */
	public toURL(): string {
		return encodeJsonToBase64Url({
			pagination: this.pagination,
			sorting: this.sorting,
			columnFilters: this.columnFilters
		});
	}

	/** Convert TanStack table state → class */
	public static fromDataTable(ts: TableState | InitialTableState): DataTableSearchParams {
		return new DataTableSearchParams({
			pagination: ts.pagination,
			sorting: ts.sorting,
			columnFilters: ts.columnFilters
		});
	}
}
