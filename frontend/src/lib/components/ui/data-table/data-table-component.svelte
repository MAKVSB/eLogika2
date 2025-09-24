<script lang="ts" generics="TData, TValue">
	import {
		getCoreRowModel,
		type ColumnDef,
		type PaginationState,
		type RowSelectionState,
		type ColumnFiltersState,
		type SortingState,
		type TableState,
		type InitialTableState,
		type ColumnSizingState,
		getFilteredRowModel,
		getPaginationRowModel,
		getSortedRowModel
	} from '@tanstack/table-core';
	import { createSvelteTable, FlexRender } from '$lib/components/ui/data-table/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Select from '$lib/components/ui/select';
	import { type Filter, type FilterSelect, FilterTypeEnum } from './filter';
	import { m } from '$lib/paraglide/messages';
	import { page } from '$app/state';
	import { decodeBase64UrlToJson, encodeJsonToBase64Url } from '$lib/services/api.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { type Table as Tabl } from '@tanstack/table-core';

	type DataTableProps<TData, TValue> = {
		columns: ColumnDef<TData, TValue>[];
		filters?: Filter[];
		data: TData[];
		refetch?: (state: TableState, queryString: string) => void;
		selection?: (rowSelection: RowSelectionState, all: boolean) => void;
		initialState?: InitialTableState;
		rowCount?: number;
		paginationEnabled?: boolean;
		sortingEnabled?: boolean;
		filterEnabled?: boolean;
		selectionEnabled?: boolean;
		queryParam?: string;
		frontEndMode?: boolean;
	};

	let {
		data,
		columns,
		filters,
		refetch,
		queryParam,
		selection,
		initialState,
		rowCount: rc,
		paginationEnabled = true,
		sortingEnabled = true,
		filterEnabled = true,
		selectionEnabled = true,
		frontEndMode = false
	}: DataTableProps<TData, TValue> = $props();

	let pagination = $state<PaginationState>({ pageIndex: 0, pageSize: 25 });
	let rowSelection = $state<RowSelectionState>({});
	let columnFilters = $state<ColumnFiltersState>([]);
	let sorting = $state<SortingState>([]);
	let columnSizing = $state<ColumnSizingState>({});

	let refetch_timer: number;
	const refetch_debounce = () => {
		clearTimeout(refetch_timer);
		refetch_timer = setTimeout(() => {
			if (!table) return;

			const newUrl = new URL(page.url);
			const state = table.getState();
			const queryParams = {
				...(state.pagination ? { pagination: state.pagination } : {}),
				...(state.sorting ? { sorting: state.sorting } : {}),
				...(state.columnFilters ? { columnFilters: state.columnFilters } : {})
			};
			const queryString = encodeJsonToBase64Url(queryParams);

			if (refetch) {
				refetch(table.getState(), queryString);
			}
			if (queryParam) {
				newUrl.searchParams.set(queryParam, queryString);
				console.log("Transfering 11")
				goto(newUrl);
			}
		}, 0);
	};
	let selection_timer: number;
	const selection_debounce = () => {
		clearTimeout(selection_timer);
		selection_timer = setTimeout(() => {
			if (!table) return;
			if (selection) {
				selection(rowSelection, table.getIsAllPageRowsSelected());
			}
		}, 0);
	};

	$effect(() => {
		if (!table) return;
		table.setOptions((prev: any) => ({
			...prev,
			get data() {
				return data;
			},
			...(paginationEnabled
				? {
						get rowCount() {
							return rc;
						},
						get pageCount() {
							console.log(rc, pagination.pageSize);
							return rc ? Math.ceil(rc / pagination.pageSize) : -1;
						},
						onPaginationChange: (updater: any) => {
							refetch_debounce();
							if (typeof updater === 'function') {
								pagination = updater(pagination);
							} else {
								pagination = updater;
							}
						},
						manualPagination: true
					}
				: {})
		}));
		console.log(table.getCanNextPage());
		// table.setPageIndex(table.getState().pagination.pageIndex)
	});

	let table: null | Tabl<TData> = $state(null);

	onMount(() => {
		if (queryParam) {
			const encodedParams = page.url.searchParams.get(queryParam);
			if (encodedParams) {
				initialState = decodeBase64UrlToJson(encodedParams);
			}
		}

		table = createSvelteTable({
			get data() {
				return data;
			},
			columns,
			...(paginationEnabled
				? {
						get rowCount() {
							return rc;
						},
						get pageCount() {
							return rc ? Math.ceil(rc / pagination.pageSize) : 1;
						},
						onPaginationChange: (updater) => {
							refetch_debounce();
							if (typeof updater === 'function') {
								pagination = updater(pagination);
							} else {
								pagination = updater;
							}
						},
						manualPagination: true
					}
				: {}),
			...(sortingEnabled
				? {
						onSortingChange: (updater) => {
							refetch_debounce();
							if (typeof updater === 'function') {
								sorting = updater(sorting);
							} else {
								sorting = updater;
							}
						},
						manualSorting: true,
						maxMultiSortColCount: 1
					}
				: {}),
			...(filterEnabled
				? {
						onColumnFiltersChange: (updater) => {
							refetch_debounce();
							if (typeof updater === 'function') {
								columnFilters = updater(columnFilters);
							} else {
								columnFilters = updater;
							}
						},
						manualFiltering: true
					}
				: {}),
			...(selectionEnabled
				? {
						onRowSelectionChange: (updater) => {
							selection_debounce();
							if (typeof updater === 'function') {
								rowSelection = updater(rowSelection);
							} else {
								rowSelection = updater;
							}
						}
					}
				: {}),
			getCoreRowModel: getCoreRowModel(),
			...(frontEndMode
				? {
						...(paginationEnabled
							? {
									getPaginationRowModel: getPaginationRowModel()
								}
							: {}),
						...(sortingEnabled
							? {
									getSortedRowModel: getSortedRowModel()
								}
							: {}),
						...(filterEnabled
							? {
									getFilteredRowModel: getFilteredRowModel()
								}
							: {})
					}
				: {}),
			initialState,
			getRowId(originalRow, index, parent) {
				if ('id' in (originalRow as any)) {
					return (originalRow as any).id;
				}
				return index;
			},
			state: {
				get pagination() {
					return pagination;
				},
				get sorting() {
					return sorting;
				},
				get columnFilters() {
					return columnFilters;
				},
				get rowSelection() {
					return rowSelection;
				},
				get columnSizing() {
					return columnSizing;
				}
			}
		});

		if (initialState) {
			if (initialState.columnFilters) {
				table.setColumnFilters(initialState.columnFilters);
			}
			if (initialState.sorting) {
				table.setSorting(initialState.sorting);
			}
			if (initialState.pagination) {
				if (initialState.pagination.pageIndex) {
					table.setPageIndex(initialState.pagination.pageIndex);
				}
				if (initialState.pagination.pageSize) {
					table.setPageSize(initialState.pagination.pageSize);
				}
			}
		}
	});

	const getFilterValue = (filter: FilterSelect, accessorKey: string) => {
		if (!table) return;
		const value = table.getColumn(accessorKey)?.getFilterValue() as string;
		if (value) {
			const valueItem = filter.values.find((v) => v.value == value);
			if (valueItem) {
				return valueItem.display;
			} else {
				return value;
			}
		} else {
			return filter.placeholder;
		}
	};
</script>

{#if table}
	<div class="flex flex-col gap-4">
		{#if filters && filters.length > 0}
			<div class="flex flex-col gap-2 lg:flex-row lg:items-center">
				{#each filters as filter}
					{#if filter.type === FilterTypeEnum.SELECT}
						<Select.Root
							type="single"
							value={table.getColumn(filter.accessorKey)?.getFilterValue() as string}
						>
							<Select.Trigger>
								{getFilterValue(filter, filter.accessorKey)}
							</Select.Trigger>
							<Select.Content>
								<Select.Item
									value=""
									onclick={() => {
										table && table.getColumn(filter.accessorKey)?.setFilterValue(null);
									}}>{filter.emptyValue}</Select.Item
								>
								{#each filter.values as value}
									<Select.Item
										value={value.value as string}
										onclick={() => {
											table && table.getColumn(filter.accessorKey)?.setFilterValue(value.value);
										}}>{value.display}</Select.Item
									>
								{/each}
							</Select.Content>
						</Select.Root>
					{:else if filter.type === FilterTypeEnum.STRING}
						<Input
							placeholder={filter.placeholder}
							value={(table.getColumn(filter.accessorKey)?.getFilterValue() as string) ?? ''}
							onchange={(e) => {
								table && table.getColumn(filter.accessorKey)?.setFilterValue(e.currentTarget.value);
							}}
							class="max-w-sm"
						/>
					{/if}
				{/each}
			</div>
		{/if}

		<div class="border rounded-md">
			<Table.Root>
				<Table.Header>
					{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
						<Table.Row>
							{#each headerGroup.headers as header (header.id)}
								<Table.Head style={`width: ${header.getSize()}px`} class="relative">
									{#if !header.isPlaceholder}
										<FlexRender
											content={header.column.columnDef.header}
											context={header.getContext()}
										/>
									{/if}
									<!-- {#if header.column.getCanResize()}
								<div
									aria-label="resize"
									onmousedown={header.getResizeHandler()}
									ontouchstart={header.getResizeHandler()}
									ondblclick={() => header.column.resetSize()}
									class="absolute top-0 right-0 h-full"
								>
									&nbsp;
								</div>
							{/if} -->
								</Table.Head>
							{/each}
						</Table.Row>
					{/each}
				</Table.Header>
				<Table.Body>
					{#each table.getRowModel().rows as row (row.id)}
						<Table.Row data-state={row.getIsSelected() && 'selected'}>
							{#each row.getVisibleCells() as cell (cell.id)}
								<Table.Cell>
									<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
								</Table.Cell>
							{/each}
						</Table.Row>
					{:else}
						<Table.Row>
							<Table.Cell colspan={columns.length} class="h-24 text-center"
								>{m.datatable_noresults()}</Table.Cell
							>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
			{#key rc}
				<!-- {rc} {table.getPageCount()} -->
				{#if paginationEnabled}
					<div class="flex items-center justify-end p-4 space-x-4">
						<p>
							{m.page_number({
								current: table.getState().pagination.pageIndex + 1,
								total: table.getPageCount()
							})}
						</p>
						<Button
							variant="outline"
							size="sm"
							onclick={() => table && table.previousPage()}
							disabled={!table.getCanPreviousPage()}
						>
							{m.pagination_previous()}
						</Button>
						<Button
							variant="outline"
							size="sm"
							onclick={() => table && table.nextPage()}
							disabled={!table.getCanNextPage()}
						>
							{m.pagination_next()}
						</Button>
					</div>
				{/if}
			{/key}
		</div>
	</div>
{/if}
