<script module lang="ts">
	export enum DataTableActionMode {
		BACKEND,
		FRONTEND,
		DISABLED
	}
</script>

<script lang="ts" generics="TData, TValue">
	import {
		type PaginationState,
		type RowSelectionState,
		type ColumnFiltersState,
		type SortingState,
		type TableState,
		type InitialTableState,
		type ColumnSizingState,
		getCoreRowModel,
		getFilteredRowModel,
		getPaginationRowModel,
		getSortedRowModel,
		type TableOptions,
		type VisibilityState
	} from '@tanstack/table-core';
	import {
		createSvelteTable,
		FlexRender,
		type ColDef
	} from '$lib/components/ui/data-table/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Select from '$lib/components/ui/select';
	import { type Filter, type FilterSelect, FilterTypeEnum } from './filter';
	import { m } from '$lib/paraglide/messages';
	import { page } from '$app/state';
	import { decodeBase64UrlToJson } from '$lib/services/api.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { type Table as Tabl } from '@tanstack/table-core';
	import { DataTableSearchParams } from '$lib/api_types_static';
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import Loader from '../loader/loader.svelte';

	type DataTableProps<TData, TValue> = {
		data: TData[];
		rowCount: number;
		columns: ColDef<TData, TValue>[];
		filters?: Filter[];
		refetch?: (state: TableState, queryString: string) => void;
		selectionChange?: (rowSelection: RowSelectionState, all: boolean) => void;
		initialState?: InitialTableState;
		paginationMode?: DataTableActionMode;
		sortingMode?: DataTableActionMode;
		filterMode?: DataTableActionMode;
		searchParam?: string;
		replaceState?: boolean;
	};

	let {
		data,
		columns,
		filters,
		refetch,
		selectionChange,
		initialState,
		rowCount: rc,
		paginationMode = DataTableActionMode.BACKEND,
		sortingMode = DataTableActionMode.BACKEND,
		filterMode = DataTableActionMode.BACKEND,
		searchParam,
		replaceState = false
	}: DataTableProps<TData, TValue> = $props();

	let paginationState = $state<PaginationState>({
		pageIndex: 0,
		pageSize: 25
	});
	let rowSelectionState = $state<RowSelectionState>({});
	let filtersState = $state<ColumnFiltersState>([]);
	let sortingState = $state<SortingState>([]);
	let sizingState = $state<ColumnSizingState>({});
	let visibilityState = $state<VisibilityState>({});

	let reloading = $state(true);
	let queryStringCache = $state<string | null>(null);

	$effect(() => {
		if (data) {
			reloading = false;
		}
	});

	let refetch_timer: number;
	const refetch_debounce = () => {
		clearTimeout(refetch_timer);
		refetch_timer = setTimeout(() => {
			if (!table) return;

			const state = table.getState();
			const queryString = DataTableSearchParams.fromDataTable(state).toURL();

			if (queryStringCache != queryString) {
				reloading = true;
			}
			queryStringCache = queryString;

			if (refetch) {
				refetch(state, queryString);
			}
			if (searchParam) {
				const newUrl = new URL(page.url);
				newUrl.searchParams.set(searchParam, queryString);
				goto(newUrl, { replaceState });
			}
		}, 0);
	};
	let selection_timer: number;
	const selection_debounce = () => {
		clearTimeout(selection_timer);
		selection_timer = setTimeout(() => {
			if (!table) return;
			if (selectionChange) {
				selectionChange(rowSelectionState, table.getIsAllPageRowsSelected());
			}
		}, 0);
	};

	let table: null | Tabl<TData> = $state(null);

	onMount(() => {
		if (searchParam) {
			const encodedParams = page.url.searchParams.get(searchParam);
			queryStringCache = encodedParams;
			if (encodedParams) {
				initialState = {
					...initialState,
					...decodeBase64UrlToJson(encodedParams)
				};
			}
		}

		if (initialState) {
			if (initialState.columnFilters) {
				filtersState = initialState.columnFilters;
			}
			if (initialState.sorting) {
				sortingState = initialState.sorting;
			}
			if (initialState.pagination) {
				paginationState = {
					...paginationState,
					...initialState.pagination
				};
			}
			if (initialState.columnVisibility) {
				visibilityState = initialState.columnVisibility;
			}
		}

		const tableOptions: Partial<TableOptions<TData>> = {
			get data() {
				return data;
			},
			columns,
			getCoreRowModel: getCoreRowModel(),
			getRowId(originalRow, index): string {
				return String(paginationState.pageIndex * paginationState.pageSize + index);
			},
			onColumnVisibilityChange: (updater) => {
				if (typeof updater === 'function') {
					visibilityState = updater(visibilityState);
				} else {
					visibilityState = updater;
				}
			},
			state: {
				get pagination() {
					return paginationState;
				},
				get sorting() {
					return sortingState;
				},
				get columnFilters() {
					return filtersState;
				},
				get rowSelection() {
					return rowSelectionState;
				},
				get columnSizing() {
					return sizingState;
				},
				get columnVisibility() {
					return visibilityState;
				}
			}
		};

		switch (paginationMode) {
			case DataTableActionMode.BACKEND:
				tableOptions.rowCount = rc;
				tableOptions.manualPagination = true;
				tableOptions.onPaginationChange = (updater) => {
					refetch_debounce();
					if (typeof updater === 'function') {
						paginationState = updater(paginationState);
					} else {
						paginationState = updater;
					}
				};
				break;
			case DataTableActionMode.FRONTEND:
				tableOptions.rowCount = rc;
				tableOptions.getPaginationRowModel = getPaginationRowModel();
				break;
		}

		switch (sortingMode) {
			case DataTableActionMode.BACKEND:
				tableOptions.manualSorting = true;
				tableOptions.maxMultiSortColCount = 1;
				tableOptions.onSortingChange = (updater) => {
					refetch_debounce();
					if (typeof updater === 'function') {
						sortingState = updater(sortingState);
					} else {
						sortingState = updater;
					}
				};
				break;
			case DataTableActionMode.FRONTEND:
				tableOptions.maxMultiSortColCount = 1;
				tableOptions.getSortedRowModel = getSortedRowModel();
				break;
		}

		switch (filterMode) {
			case DataTableActionMode.BACKEND:
				tableOptions.manualFiltering = true;
				tableOptions.onColumnFiltersChange = (updater) => {
					refetch_debounce();
					if (typeof updater === 'function') {
						filtersState = updater(filtersState);
					} else {
						filtersState = updater;
					}
				};
				break;
			case DataTableActionMode.FRONTEND:
				tableOptions.getFilteredRowModel = getFilteredRowModel();
				break;
		}

		if (columns.find((v) => (v.id = 'select'))) {
			tableOptions.onRowSelectionChange = (updater) => {
				selection_debounce();
				if (typeof updater === 'function') {
					rowSelectionState = updater(rowSelectionState);
				} else {
					rowSelectionState = updater;
				}
			};
		}

		table = createSvelteTable(tableOptions as TableOptions<TData>);
	});

	let filterStrings = $derived.by(() => {
		const res: { [key: string]: string } = {};

		for (const filter of filters ?? []) {
			switch (filter.type) {
				case FilterTypeEnum.SELECT:
					res[filter.accessorKey] = getFilterValue(filter);
			}
		}

		return res;
	});

	const getFilterValue = (filter: FilterSelect) => {
		if (!table) return '';
		const value = table.getColumn(filter.accessorKey)?.getFilterValue() as string;
		if (value) {
			const valueItem = filter.values.find((v) => v.value == value);
			if (valueItem) {
				return valueItem.display;
			}
			return value;
		} else {
			return filter.placeholder;
		}
	};
</script>

{#if table}
	<div class="flex flex-col gap-4">
		{#if filters && filters.length > 0}
			<div class="flex flex-col gap-2 overflow-scroll lg:flex-row lg:items-center">
				{#each filters as filter}
					{#if filter.type === FilterTypeEnum.SELECT}
						<Select.Root
							type="single"
							value={table.getColumn(filter.accessorKey)?.getFilterValue() as string}
						>
							<Select.Trigger>
								{filterStrings[filter.accessorKey]}
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
											if (table) {
												const column = table.getColumn(filter.accessorKey);
												if (column) {
													column.setFilterValue(String(value.value));
												} else {
													console.error('COLUMN NOT FOUND', filter.accessorKey);
												}
											} else {
												console.error('TABLE NOT INITIATED');
											}
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
					{#if reloading}
						<Table.Row>
							<Table.Cell colspan={200} class="w-full">
								<div class="flex justify-center mx-auto text-center">
									<Loader width={150} height={150}></Loader>
								</div>
							</Table.Cell>
						</Table.Row>
					{:else}
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
					{/if}
				</Table.Body>
			</Table.Root>
			{#key rc}
				<div class="flex w-full p-4 space-x-4 overflow-scroll justify-evenly">
					{#if paginationMode != DataTableActionMode.DISABLED}
						<p class="w-full">
							{m.page_number({
								current: table.getState().pagination.pageIndex + 1,
								total: Math.ceil(rc / paginationState.pageSize)
							})}
						</p>
					{/if}
					<DropdownMenu.Root>
						<DropdownMenu.Trigger>
							{#snippet child({ props })}
								<Button {...props} variant="outline">{m.column_visibility()}</Button>
							{/snippet}
						</DropdownMenu.Trigger>
						<DropdownMenu.Content class="w-56">
							<DropdownMenu.Group>
								{#each table.getAllColumns().filter((v) => v.getCanHide()) as column}
									<DropdownMenu.CheckboxItem
										bind:checked={
											() => column.getIsVisible(),
											(v) => {
												column.toggleVisibility(!!v);
											}
										}
									>
										{(column.columnDef as ColDef<TData>).columnName}
									</DropdownMenu.CheckboxItem>
								{/each}
							</DropdownMenu.Group>
						</DropdownMenu.Content>
					</DropdownMenu.Root>
					{#if paginationMode != DataTableActionMode.DISABLED}
						<Pagination.Root
							class="justify-end"
							count={rc}
							page={paginationState.pageIndex + 1}
							perPage={paginationState.pageSize}
							onPageChange={(e) => {
								table?.setPageIndex(e - 1);
							}}
						>
							{#snippet children({ pages, currentPage })}
								<Pagination.Content>
									<Pagination.Item>
										<Pagination.PrevButton />
									</Pagination.Item>
									{#each pages as page (page.key)}
										{#if page.type === 'ellipsis'}
											<Pagination.Item>
												<Pagination.Ellipsis />
											</Pagination.Item>
										{:else}
											<Pagination.Item>
												<Pagination.Link {page} isActive={currentPage === page.value}>
													{page.value}
												</Pagination.Link>
											</Pagination.Item>
										{/if}
									{/each}
									<Pagination.Item>
										<Pagination.NextButton />
									</Pagination.Item>
								</Pagination.Content>
							{/snippet}
						</Pagination.Root>
					{/if}
				</div>
			{/key}
		</div>
	</div>
{/if}
