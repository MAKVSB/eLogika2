<script lang="ts" module>
	import UsersIcon from '@lucide/svelte/icons/users';
	import QuestionIcon from '@lucide/svelte/icons/file-question';
	import PenLineIcon from '@lucide/svelte/icons/pen-line';
	import NotebookPenIcon from '@lucide/svelte/icons/notebook-pen';
	import CalendarIcon from '@lucide/svelte/icons/calendar-check';
	import LifeBuoyIcon from '@lucide/svelte/icons/life-buoy';
	import GalleryVerticalEndIcon from '@lucide/svelte/icons/gallery-vertical-end';
	import AudioWaveformIcon from '@lucide/svelte/icons/audio-waveform';
	import GlobalState from '$lib/shared.svelte';
	import { base } from '$app/paths';

	let menuTree: (SidebarCategory | SidebarFlex)[] = $derived.by(() => {
		return [
			{
				name: m.menu_administrator(),
				type: MenuTreeType.CATEGORY,
				items: [
					{
						title: m.menu_administrator_courses(),
						url: base+`/app/admin/courses`,
						icon: GalleryVerticalEndIcon
					},
					{
						title: m.menu_administrator_users(),
						url: base+`/app/admin/users`,
						icon: UsersIcon
					}
				],
				requiredRoles: [CourseUserRoleEnum.ADMIN],
				noCourse: true
			},
			{
				name: m.menu_tutor(),
				type: MenuTreeType.CATEGORY,
				items: [
					{
						title: m.menu_tutor_course(),
						url: base+`/app/${GlobalState.activeCourse?.id}/tutor/course`,
						icon: GalleryVerticalEndIcon
					},
					{
						title: m.menu_tutor_courseitems(),
						url: base+`/app/${GlobalState.activeCourse?.id}/tutor/courseitems`,
						icon: GalleryVerticalEndIcon
					},
					{
						title: m.menu_tutor_classes(),
						url: base+`/app/${GlobalState.activeCourse?.id}/tutor/classes`,
						icon: UsersIcon
					},
					{
						title: m.menu_tutor_users(),
						url: base+`/app/${GlobalState.activeCourse?.id}/tutor/users`,
						icon: UsersIcon
					},
					{
						title: m.menu_tutor_questions(),
						url: base+`/app/${GlobalState.activeCourse?.id}/tutor/questions/`,
						icon: QuestionIcon
					},
					{
						title: m.menu_tutor_categories(),
						url: base+`/app/${GlobalState.activeCourse?.id}/tutor/categories/`,
						icon: QuestionIcon
					},
					{
						title: m.menu_tutor_chapters(),
						url: base+`/app/${GlobalState.activeCourse?.id}/tutor/chapters/${GlobalState.activeCourse?.chapterId}`,
						icon: AudioWaveformIcon //TODO
					},
					{
						title: m.menu_tutor_templates(),
						url: base+`/app/${GlobalState.activeCourse?.id}/tutor/templates`,
						icon: AudioWaveformIcon //TODO
					},
					// {
					// 	title: m.menu_tutor_consultation(),
					// 	url: base+`/app/${GlobalState.activeCourse?.id}/tutor/consultation`,
					// 	icon: CalendarIcon
					// }
				],
				requiredRoles: [
					CourseUserRoleEnum.ADMIN,
					CourseUserRoleEnum.GARANT,
					CourseUserRoleEnum.TUTOR
				]
			},
			{
				type: MenuTreeType.CATEGORY,
				name: m.menu_student(),
				items: [
					{
						title: m.menu_student_results(),
						url: base+`/app/${GlobalState.activeCourse?.id}/student/results`,
						icon: NotebookPenIcon //TODO
					},
					{
						title: m.menu_student_materials(),
						url: base+`/app/${GlobalState.activeCourse?.id}/student/materials/${GlobalState.activeCourse?.chapterId}`,
						icon: AudioWaveformIcon //TODO
					},
					{
						title: m.menu_student_terms(),
						url: base+`/app/${GlobalState.activeCourse?.id}/student/terms`,
						icon: CalendarIcon //TODO
					},
					{
						title: m.menu_student_tests(),
						url: base+`/app/${GlobalState.activeCourse?.id}/student/tests`,
						icon: PenLineIcon //TODO
					},
					{
						title: m.menu_student_homeworks(),
						url: base+`/app/${GlobalState.activeCourse?.id}/student/activities`,
						icon: PenLineIcon //TODO
					},
					// {
					// 	title: m.menu_student_consultation(),
					// 	url: base+`/app/${GlobalState.activeCourse?.id}/student/consultations`,
					// 	icon: CalendarIcon
					// }
				],
				requiredRoles: [CourseUserRoleEnum.STUDENT]
			},
			{
				type: MenuTreeType.FLEX_GROW
			},
			{
				type: MenuTreeType.CATEGORY,
				name: ``,
				items: [
					{
						title: m.menu_support(),
						url: base+`/app/support`,
						icon: LifeBuoyIcon
					}
				],
				requiredRoles: []
			}
		];
	});
</script>

<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import NavUser from './nav-user.svelte';
	import CourseSwitcher from './course-switcher.svelte';
	import NavCategory, {
		MenuTreeType,
		type SidebarCategory,
		type SidebarFlex
	} from './nav-category.svelte';
	import { m } from '$lib/paraglide/messages';
	import { CourseUserRoleEnum } from '$lib/api_types';

	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();
</script>

<Sidebar.Root class="top-(--header-height) h-[calc(100svh-var(--header-height))]!" {...restProps}>
	<Sidebar.Header>
		{#if GlobalState.loggedUser != null}
			<CourseSwitcher
				courses={GlobalState.availableCourses}
			/>
		{/if}
	</Sidebar.Header>
	<Sidebar.Content>
		{#each menuTree as menuCategory}
			{#if menuCategory.type == MenuTreeType.CATEGORY}
				<NavCategory {...menuCategory}></NavCategory>
			{:else if menuCategory.type == MenuTreeType.FLEX_GROW}
				<div class="flex-grow"></div>
			{/if}
		{/each}
	</Sidebar.Content>
	<Sidebar.Footer>
		{#if GlobalState.loggedUser != null}
			<NavUser bind:user={GlobalState.loggedUser} />
		{/if}
	</Sidebar.Footer>
</Sidebar.Root>
