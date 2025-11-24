import type {
	CategoryListResponse,
	ChapterListResponse,
	QuestionListResponse
} from '$lib/api_types.js';
import { DataTableSearchParams } from '$lib/api_types_static';
import { API } from '$lib/services/api.svelte.js';
import { tableConfig } from './schema.svelte';

export async function load({ fetch, params, url }) {
	const search =
		url.searchParams.get(tableConfig.searchParam) ??
		DataTableSearchParams.fromDataTable(tableConfig.initialState).toURL();

	const data = API.request<null, QuestionListResponse>(
		`/api/v2/courses/${params.courseId}/questions`,
		{
			searchParams: {
				search
			}
		},
		fetch
	);

	const chapterData = API.request<null, ChapterListResponse>(
		`/api/v2/courses/${params.courseId}/chapters`,
		{},
		fetch
	);

	const categoryData = API.request<null, CategoryListResponse>(
		`/api/v2/courses/${params.courseId}/categories`,
		{},
		fetch
	);

	return { data, chapterData, categoryData };
}
