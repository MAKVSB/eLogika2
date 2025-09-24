<script lang="ts">
	import { cn, deepMerge } from '$lib/utils';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';
	import type { ErrorObject } from './types';
	import { toast } from 'svelte-sonner';
	import { Button } from '.';
	import { m } from '$lib/paraglide/messages';
	import type { Snippet } from 'svelte';

	let {
		ref = $bindable(null),
		class: className,
		children,
		onsubmit,
		form = $bindable(),
		isCreating = false,
		hideDefaultbutton = false,
		additionalButtons,
		addedValidation,
		...restProps
	}: {
		onsubmit?: (event: SubmitEvent) => Promise<any>;
		isCreating?: boolean;
		hideDefaultbutton?: boolean;
		additionalButtons?: Snippet;
		addedValidation?: (fields: any) => ErrorObject;
		form?: any; // TODO
	} & WithElementRef<HTMLAttributes<HTMLElement>> = $props();

	function validate(): boolean {
		form.errors = {};

		if (addedValidation) {
			form.errors = addedValidation(form.fields);
		}

		const result = form.schema.safeParse(form.fields);
		if (!result.success) {
			for (const issue of result.error.issues) {
				let currentLevel: ErrorObject = form.errors;

				for (let i = 0; i < issue.path.length; i++) {
					const pathPart = issue.path[i] as string;
					if (i === issue.path.length - 1) {
						currentLevel[pathPart] = issue.message;
					} else {
						if (!(pathPart in currentLevel)) {
							currentLevel[pathPart] = {};
						}
						currentLevel = currentLevel[pathPart] as ErrorObject;
					}
				}
			}
			toast.error(m.validation_error() + ' ' + JSON.stringify(form.errors));
			// TODO write whats wrong ?
			return false;
		}
		return true;
	}

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (form) {
			form.isSubmitting = true;
		}

		// Form validation
		if (form && form.schema && form.errors) {
			if (!validate()) {
				form.isSubmitting = false;
				return;
			}
		}

		if (onsubmit) {
			await onsubmit(event).catch((err) => {
				if ('formErrors' in err.data) {
					const combinedErrors = deepMerge(form.errors, err.data.formErrors);
					console.log(combinedErrors);
					form.errors = combinedErrors;
				}
			});
		}
		if (form) {
			form.isSubmitting = false;
		}
	}
</script>

<form
	onsubmit={handleSubmit}
	class={cn(
		'flex min-h-0 flex-1 flex-col gap-2 overflow-auto group-data-[collapsible=icon]:overflow-hidden',
		className
	)}
	{...restProps}
>
	{@render children?.()}

	<div class="flex justify-center gap-4">
		{#if !hideDefaultbutton && form}
			<Button
				class="w-30"
				isSubmitting={form.isSubmitting}
				text={isCreating ? m.create() : m.save()}
				textSubmiting={isCreating ? m.create_progress() : m.save_progress()}
			></Button>
		{/if}
		{@render additionalButtons?.()}
	</div>
</form>
