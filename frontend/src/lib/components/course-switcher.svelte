<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { useSidebar } from '$lib/components/ui/sidebar/index.js';
	import type { CourseUserRoleEnum, LoggedUserCourseDTO2 } from '$lib/api_types';
	import ChevronsUpDownIcon from '@lucide/svelte/icons/chevrons-up-down';
	import AudioWaveformIcon from '@lucide/svelte/icons/audio-waveform';
	import CalendarDaysIcon from '@lucide/svelte/icons/calendar-days';
	import { m } from '$lib/paraglide/messages';
	import { page } from '$app/state';
	import {} from '$app/stores';
	import { goto } from '$app/navigation';
	import GlobalState from '$lib/shared.svelte';
	import { base } from '$app/paths';

	let {
		courses
	}: {
		courses: LoggedUserCourseDTO2[];
	} = $props();

	function onlyUnique(value: any, index: number, array: number[]) {
		return array.indexOf(value) === index;
	}

	const defaultYear =
		new Date().getMonth() >= 8 ? new Date().getFullYear() : new Date().getFullYear() - 1;

	let selectedYear = $state(GlobalState.activeCourse?.year ?? defaultYear);
	let years = $derived(courses.map((c) => c.year).filter(onlyUnique));

	let coursesInYear = $derived(courses.filter((c) => c.year == selectedYear));

	function getRedirectPathOnCourseChange(
		currentPath: string,
		newCourse: LoggedUserCourseDTO2
	): string | false {
		let match;

		// TODO add every possible route to re-route on course change

		// admin course
		match = currentPath.match(/^\/app\/admin(.*)$/);
		if (match) {
			// Add logic to verify if questionId is valid for newCourseId
			return false;
		}

		// Question management
		match = currentPath.match(/^\/app\/(\d+)\/tutor\/questions(\/(\d+))?$/);
		if (match) {
			// Add logic to verify if questionId is valid for newCourseId
			return base+`/app/${newCourse.id}/tutor/questions/`;
		}
		// Chapter management
		match = currentPath.match(/^\/app\/(\d+)\/tutor\/chapters\/(\d+)$/);
		if (match) {
			// TODO Add logic to verify if questionId is valid for newCourseId
			return base+`/app/${newCourse.id}/tutor/chapters/${newCourse.chapterId}`;
		}

		// TODO add whole menu here (or figure out a way to append it to menu item structure)

		// Add more route matching logic here
		return base+`/app/${newCourse.id}`;
	}

	function yearChange(year: number) {
		selectedYear = year;

		const sameCourse = coursesInYear.find((c) => c.name === GlobalState.activeCourse?.name);

		GlobalState.activeCourse = sameCourse ?? coursesInYear[0];
		let redirectTo = getRedirectPathOnCourseChange(page.url.pathname, coursesInYear[0]);
		if (redirectTo) {
			console.log('Transfering 8');
			goto(base + redirectTo);
		}
	}

	function courseChange(course: LoggedUserCourseDTO2) {
		GlobalState.activeCourse = course;
		let redirectTo = getRedirectPathOnCourseChange(page.url.pathname, course);
		if (redirectTo) {
			console.log('Transfering 9');
			goto(redirectTo);
		}
	}

	const sidebar = useSidebar();
</script>

<Sidebar.Menu>
	<Sidebar.MenuItem>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger class="w-full">
				<Sidebar.MenuButton
					size="lg"
					class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
				>
					<div
						class="flex items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground aspect-square size-8"
					>
						<CalendarDaysIcon class="size-4" />
					</div>
					{#if selectedYear}
						<div class="grid flex-1 text-sm leading-tight text-left">
							<span class="font-medium truncate">
								{m.academic_year()}
							</span>
							<span class="text-xs truncate">
								{selectedYear}/{selectedYear + 1}
								{selectedYear == defaultYear ? '[Current]' : ''}
							</span>
						</div>
					{:else}
						<div class="grid flex-1 text-sm leading-tight text-left">
							<span class="font-medium truncate">{m.no_course_available()}</span>
						</div>
					{/if}
					<ChevronsUpDownIcon class="ml-auto" />
				</Sidebar.MenuButton>
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
				align="start"
				side={sidebar.isMobile ? 'bottom' : 'right'}
				sideOffset={4}
			>
				<DropdownMenu.Label class="text-xs text-muted-foreground">Academic years</DropdownMenu.Label
				>
				{#each years as year}
					<DropdownMenu.Item
						onSelect={() => {
							yearChange(year);
						}}
						class="gap-2 p-2 {selectedYear === year ? 'font-bold' : ''}"
					>
						{year}/{year + 1}
						{year == defaultYear ? '[Current]' : ''}
					</DropdownMenu.Item>
				{/each}
			</DropdownMenu.Content>
		</DropdownMenu.Root>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger class="w-full">
				<Sidebar.MenuButton
					size="lg"
					class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
				>
					<div
						class="flex items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground aspect-square size-8"
					>
						<AudioWaveformIcon class="size-4" />
					</div>
					{#if GlobalState.activeCourse}
						<div class="grid flex-1 text-sm leading-tight text-left">
							<span class="font-medium truncate"> Course </span>
							<span class="text-xs truncate">{GlobalState.activeCourse.name}</span>
						</div>
					{:else}
						<div class="grid flex-1 text-sm leading-tight text-left">
							<span class="font-medium truncate">{m.no_course_available()}</span>
						</div>
					{/if}
					<ChevronsUpDownIcon class="ml-auto" />
				</Sidebar.MenuButton>
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
				align="start"
				side={sidebar.isMobile ? 'bottom' : 'right'}
				sideOffset={4}
			>
				<DropdownMenu.Label class="text-xs text-muted-foreground">{m.courses()}</DropdownMenu.Label>
				{#each coursesInYear as course}
					<DropdownMenu.Item
						onSelect={() => {
							courseChange(course);
						}}
						class="gap-2 p-2 {GlobalState.activeCourse?.name === course.name ? 'font-bold' : ''}"
					>
						{course.name}
					</DropdownMenu.Item>
				{/each}
			</DropdownMenu.Content>
		</DropdownMenu.Root>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger class="w-full">
				<Sidebar.MenuButton
					size="lg"
					class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
				>
					<div
						class="flex items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground aspect-square size-8"
					>
						<AudioWaveformIcon class="size-4" />
					</div>
					{#if GlobalState.activeRole}
						<div class="grid flex-1 text-sm leading-tight text-left">
							<span class="font-medium truncate"> Role </span>
							<span class="text-xs truncate">{GlobalState.activeRole}</span>
						</div>
					{:else}
						<div class="grid flex-1 text-sm leading-tight text-left">
							<span class="font-medium truncate">No role available</span>
						</div>
					{/if}
					<ChevronsUpDownIcon class="ml-auto" />
				</Sidebar.MenuButton>
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
				align="start"
				side={sidebar.isMobile ? 'bottom' : 'right'}
				sideOffset={4}
			>
				<DropdownMenu.Label class="text-xs text-muted-foreground">Roles</DropdownMenu.Label>
				{#each GlobalState.activeCourse?.roles ?? [] as role}
					<DropdownMenu.Item
						onSelect={() => {
							GlobalState.activeRole = role as CourseUserRoleEnum;
						}}
						class="gap-2 p-2"
					>
						{role}
					</DropdownMenu.Item>
				{/each}
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>
