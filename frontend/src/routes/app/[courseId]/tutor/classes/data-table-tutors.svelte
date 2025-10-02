<script lang="ts">
	import type { ClassUserDTO } from '$lib/api_types';

	let { tutors }: { tutors: ClassUserDTO[] } = $props();
</script>

{#each tutors as user}
	<p>
		{user.degreeBefore}
		{user.firstName}
		{user.familyName}
		{user.degreeAfter}
	</p>
{/each}
