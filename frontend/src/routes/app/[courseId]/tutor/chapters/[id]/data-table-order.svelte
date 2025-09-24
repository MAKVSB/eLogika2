<script lang="ts">
	import { MoveDirectionEnum } from '$lib/api_types';
	import { Button } from '$lib/components/ui/button/index.js';

	import ArrDown from '@lucide/svelte/icons/arrow-down-from-line';
	import ArrUp from '@lucide/svelte/icons/arrow-up-from-line';

	let {
		id,
		order = $bindable(),
		meta,
		isFirst,
		isLast
	}: {
		id: number | string;
		order: number;
		meta: any;
		isFirst: boolean;
		isLast: boolean;
	} = $props();

	function move(dir: MoveDirectionEnum) {
		if ('changeEventHandler' in meta) {
			meta.changeEventHandler(id, dir);
		}
	}
</script>

<div class="flex items-center">
	<div class="flex">
		<Button
			variant="ghost"
			class="relative w-min"
			onclick={() => move(MoveDirectionEnum.UP)}
			disabled={isFirst}
		>
			<ArrUp />
		</Button>
		<Button
			variant="ghost"
			class="relative w-min"
			onclick={() => move(MoveDirectionEnum.DOWN)}
			disabled={isLast}
		>
			<ArrDown />
		</Button>
	</div>
</div>
