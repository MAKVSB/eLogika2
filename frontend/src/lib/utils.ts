import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, 'child'> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any } ? Omit<T, 'children'> : T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & { ref?: U | null };

export function enumToOptions<T extends Record<string, string>>(
	en: T,
	langFunc?: (inputs: any, options?: any) => string,
	langParams?: any
) {
	return Object.values(en).map((value) => {
		return {
			value: value,
			display: langFunc ? langFunc({ value: value, ...langParams }) : value
		};
	});
}

export function deepMerge(target: any, source: any) {
	for (const key in source) {
		if (source[key] && typeof source[key] === 'object' && !Array.isArray(source[key])) {
			if (!target[key]) target[key] = {};
			deepMerge(target[key], source[key]);
		} else if (Array.isArray(source[key]) && Array.isArray(target[key])) {
			target[key] = [...target[key], ...source[key]];
		} else {
			target[key] = source[key];
		}
	}
	return target;
}

export function displayUserName(pud: {
	degreeBefore: string;
	firstName: string;
	familyName: string;
	degreeAfter: string;
}) {
	return `${pud.degreeBefore} ${pud.firstName} ${pud.familyName}${pud.degreeAfter && pud.degreeAfter != '' ? ', ' + pud.degreeAfter : ''}`;
}
