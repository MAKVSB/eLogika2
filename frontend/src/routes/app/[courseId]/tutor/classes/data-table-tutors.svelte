<script lang="ts">
	import type { ClassUserDTO } from '$lib/api_types';
	import { displayUserName } from '$lib/utils';

	let { tutors }: { tutors: ClassUserDTO[] } = $props();
</script>

{#each tutors as user}
	<p>
		{displayUserName(user)}
	</p>
{/each}
