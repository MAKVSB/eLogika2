<script lang="ts">
	import { page } from '$app/state';
	import { StudyFormEnum } from '$lib/api_types';
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
