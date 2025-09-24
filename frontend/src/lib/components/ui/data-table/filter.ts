import type { SelectOptions } from '../form';

export enum FilterTypeEnum {
	STRING,
	SELECT
}

export type FilterString = {
	type: FilterTypeEnum.STRING;
	accessorKey: string;
	placeholder: string;
};

export type FilterSelect = {
	type: FilterTypeEnum.SELECT;
	accessorKey: string;
	values: SelectOptions;
	placeholder: string;
	emptyValue: string;
};

export type Filter = FilterString | FilterSelect;
