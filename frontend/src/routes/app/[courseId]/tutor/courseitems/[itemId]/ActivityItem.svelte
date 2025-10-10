<script lang="ts">
	import * as Form from '$lib/components/ui/form';
	import { m } from '$lib/paraglide/messages';

	let {
		fields = $bindable(),
		errors,
		disabled
	}: {
		fields: any;
		errors?: any;
		disabled?: boolean;
	} = $props();
</script>

<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
	{#key fields}
		<Form.Tiptap
			title={m.courseitem_group_activity_assignment()}
			name="description"
			id="description"
			bind:value={fields.description}
			error={errors.description}
			{disabled}
			enableFileUpload
			enableFileLink
		></Form.Tiptap>
		<Form.Tiptap
			title={m.courseitem_group_activity_expectedresult()}
			name="expectedResult"
			id="expectedResult"
			bind:value={fields.expectedResult}
			error={errors.expectedResult}
			{disabled}
			enableFileUpload
			enableFileLink
		></Form.Tiptap>
	{/key}
</div>
