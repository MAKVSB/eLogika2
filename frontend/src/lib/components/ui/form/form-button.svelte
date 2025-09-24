<script lang="ts">
	import type { ButtonProps } from '../button';
	import { Button } from '../button';
	import Loader from '../loader/loader.svelte';
	import type { WithoutChildren } from '$lib/utils';

	let {
		disabled,
		isSubmitting = $bindable(false),
		text = '',
		textSubmiting = '',
		type = 'submit',
		...restProps
	}: WithoutChildren<ButtonProps> & {
		isSubmitting?: boolean;
		text: string;
		textSubmiting: string;
	} = $props();
</script>

<Button {type} disabled={disabled || isSubmitting} {...restProps}>
	{#if isSubmitting}
		<Loader></Loader>{textSubmiting}...
	{:else}
		{text}
	{/if}
</Button>
