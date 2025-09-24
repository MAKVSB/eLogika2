<script lang="ts">
	import type { ClassUserDTO } from '$lib/api_types';

	let { tutors }: { tutors: ClassUserDTO[] } = $props();
</script>

{#each tutors as user}
	<p>
		{user.firstName}
		{user.familyName}
	</p>
{/each}
