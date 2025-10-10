<script lang="ts">
	import { page } from '$app/state';
	import { QuestionTypeEnum, type TemplateListResponse } from '$lib/api_types';
	import * as Form from '$lib/components/ui/form';
	import { m } from '$lib/paraglide/messages';
	import { API } from '$lib/services/api.svelte';
	import { enumToOptions } from '$lib/utils';
	import { onMount } from 'svelte';

	let {
		fields = $bindable(),
		errors,
		disabled
	}: {
		fields: any;
		errors?: any;
		disabled?: boolean;
	} = $props();

	let templates: Form.SelectOptions = $state([]);

	onMount(() => {
		API.request<null, TemplateListResponse>(`/api/v2/courses/${page.params.courseId}/templates`, {})
			.then((res) => {
				templates = res.items.map((t) => {
					return {
						value: t.id,
						display: t.title
					};
				});
			})
			.catch(() => {});
	});
</script>

<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
	<Form.SingleSelect
		title={m.courseitem_group_test_type()}
		name="testType"
		id="testType"
		class="sm:col-span-2 lg:col-span-3"
		options={enumToOptions(QuestionTypeEnum, m.question_type_enum)}
		bind:value={fields.testType}
		error={errors.testType}
		{disabled}
	></Form.SingleSelect>
	<Form.SingleSelect
		title={m.courseitem_group_test_template()}
		name="testTemplateId"
		id="testTemplateId"
		class="sm:col-span-2 lg:col-span-3"
		options={templates}
		bind:value={fields.testTemplateId}
		error={errors.testTemplateId}
		{disabled}
	></Form.SingleSelect>
	<Form.TextInput
		title={m.courseitem_group_test_time_limit()}
		name="timeLimit"
		id="timeLimit"
		type="number"
		class="sm:col-span-2 lg:col-span-3"
		bind:value={fields.timeLimit}
		error={errors.timeLimit ?? ''}
		{disabled}
	></Form.TextInput>
	<Form.Checkbox
		title={m.courseitem_group_test_show_results()}
		name="showResults"
		id="showResults"
		bind:value={fields.showResults}
		error={errors.showResults}
		tooltip={m.courseitem_group_test_show_results_tooltip()}
		{disabled}
	></Form.Checkbox>
	<Form.Checkbox
		title={m.courseitem_group_test_show_test()}
		name="showTest"
		id="showTest"
		bind:value={fields.showTest}
		error={errors.showTest}
		tooltip={m.courseitem_group_test_show_test_tooltip()}
		{disabled}
	></Form.Checkbox>
	<Form.Checkbox
		title={m.courseite_group_test_show_correctness()}
		name="showTest"
		id="showTest"
		bind:value={fields.showCorrectness}
		error={errors.showCorrectness}
		tooltip={m.courseitem_group_test_show_correctness_tooltip()}
		{disabled}
	></Form.Checkbox>
	<Form.Checkbox
		title={m.courseitem_group_test_offline()}
		name="allowOffline"
		id="allowOffline"
		bind:value={fields.allowOffline}
		error={errors.allowOffline}
		tooltip={m.courseitem_group_test_offline_tooltip()}
		disabled
	></Form.Checkbox>
	<Form.Checkbox
		title={m.courseitem_group_test_ispaper()}
		name="isPaper"
		id="isPaper"
		bind:value={fields.isPaper}
		error={errors.isPaper}
		tooltip={m.courseitem_group_test_ispaper_tooltip()}
		{disabled}
	></Form.Checkbox>
	<Form.TextInput
		title={m.courseitem_group_ip_range()}
		name="ipRanges"
		id="ipRanges"
		type="text"
		class="sm:col-span-2 lg:col-span-3"
		bind:value={fields.ipRanges}
		error={errors.ipRanges ?? ''}
		{disabled}
	></Form.TextInput>
</div>
