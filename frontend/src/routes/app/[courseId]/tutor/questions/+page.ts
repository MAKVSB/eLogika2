import type {
	CategoryListResponse,
	ChapterListResponse,
	QuestionListResponse
} from '$lib/api_types.js';
import { API, encodeJsonToBase64Url } from '$lib/services/api.svelte.js';

export async function load({ fetch, params, url }) {
	const search = url.searchParams.get('search');
	const disablePaginationSearchParams = {
		search: encodeJsonToBase64Url({
			pagination: {
				pageIndex: 0,
				pageSize: 100000
			}
		})
	};

	const data = API.request<null, QuestionListResponse>(
		`/api/v2/courses/${params.courseId}/questions`,
		{
			searchParams: {
				...(search ? { search: search } : {})
			}
		},
		fetch
	);

	const chapterData = API.request<null, ChapterListResponse>(
		`/api/v2/courses/${params.courseId}/chapters`,
		{
			searchParams: disablePaginationSearchParams
		},
		fetch
	);

	const categoryData = API.request<null, CategoryListResponse>(
		`/api/v2/courses/${params.courseId}/categories`,
		{
			searchParams: disablePaginationSearchParams
		},
		fetch
	);

	return { data, chapterData, categoryData };
}
