<script lang="ts">
	import { API, ApiError } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import { StudyFormEnum } from '$lib/api_types';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import { m } from '$lib/paraglide/messages';
	import CourseItems from './CourseItems/CourseItems.svelte';
</script>

<div class="flex flex-col gap-8 m-8">
	<div class="flex flex-row justify-between">
		<h1 class="text-2xl">
			Course items:
		</h1>
	</div>
	<CourseItems canGroup={true} courseId={page.params.courseId} mode={StudyFormEnum.FULLTIME}
	></CourseItems>
	<CourseItems canGroup={true} courseId={page.params.courseId} mode={StudyFormEnum.COMBINED}
	></CourseItems>
</div>
