<script lang="ts">
	import type { ErrorResponse } from '$lib/api_types';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages';

	let errorData: ErrorResponse = $props();
</script>

<div>
	{errorData.message}
</div>
<div>
	{errorData.details}
</div>
<div>
	{#if errorData.fileData}
		<Button
			variant="destructive"
			size="sm"
			onclick={() => {
				if (errorData.fileData) {
					const binaryString = atob(errorData.fileData.content);
					const len = binaryString.length;

					const bytes = new Uint8Array(len);
					for (let i = 0; i < len; i++) {
						bytes[i] = binaryString.charCodeAt(i);
					}

					const blob = new Blob([bytes], { type: errorData.fileData.mimeType });

					window.open(URL.createObjectURL(blob));
				}
			}}
		>
			({m.error_log_download()})
		</Button>
	{/if}
</div>
