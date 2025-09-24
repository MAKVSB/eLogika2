<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button/index.js';
	import { m } from '$lib/paraglide/messages';
	import { API, ApiError } from '$lib/services/api.svelte';
	import type { TermsJoinResponse } from '$lib/api_types';
	import { invalidateAll } from '$app/navigation';

	let {
		userId
	}: {
		userId: number | string;
	} = $props();

	async function termSignOut() {
		await API.request<any, TermsJoinResponse>(
			`/api/v2/courses/${page.params.courseId}/items/${page.params.itemId}/terms/${page.params.termId}/students`,
			{
				method: 'DELETE',
				body: {
					userId: userId
				}
			}
		)
			.then((res) => {
				invalidateAll();
			})
			.catch(() => {});
	}
</script>

<div class="flex justify-between">
	<Button variant="ghost" class="relative" onclick={() => termSignOut()}>
		<span>{m.term_signout_student()}</span>
	</Button>
</div>
