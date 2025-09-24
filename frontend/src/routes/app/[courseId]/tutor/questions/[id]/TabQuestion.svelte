<script lang="ts">
	import { QuestionFormatEnum, QuestionTypeEnum } from '$lib/api_types';
	import * as Form from '$lib/components/ui/form';
	import { m } from '$lib/paraglide/messages';
	import { enumToOptions } from '$lib/utils';

	let {
		form = $bindable()
	}: {
		form: any;
	} = $props();
</script>

<div class="flex flex-col gap-4 p-2 my-4">
	<div class="grid grid-cols-12 gap-4">
		<Form.TextInput
			title={m.question_title()}
			name="title"
			id="title"
			type="text"
			class="col-span-12 sm:col-span-8"
			bind:value={form.fields.title}
			error={form.errors.title}
		></Form.TextInput>
		<Form.Checkbox
			title={m.question_active()}
			name="active"
			id="active"
			class="col-span-12 sm:col-span-4"
			wide={true}
			bind:value={form.fields.active}
			error={form.errors.active}
		></Form.Checkbox>
		<Form.SingleSelect
			title={m.question_questiontype()}
			name="questionType"
			id="questionType"
			class="col-span-12 sm:col-span-4"
			bind:value={form.fields.questionType}
			options={enumToOptions(QuestionTypeEnum, m.question_type_enum)}
			error={form.errors.questionType}
		></Form.SingleSelect>
		<Form.SingleSelect
			title={m.question_questionformat()}
			name="questionFormat"
			id="questionFormat"
			class="col-span-12 sm:col-span-4"
			bind:value={form.fields.questionFormat}
			options={enumToOptions(QuestionFormatEnum, m.question_format_enum)}
			error={form.errors.questionFormat}
		></Form.SingleSelect>
		<Form.TextInput
			title={m.question_question_timetoread()}
			name="timeToRead"
			id="timeToRead"
			type="number"
			class="col-span-12 sm:col-span-6"
			bind:value={form.fields.timeToRead}
			error={form.errors.timeToRead}
		></Form.TextInput>
		<Form.TextInput
			title={m.question_timetoprocess()}
			name="timeToProcess"
			id="timeToProcess"
			type="number"
			class="col-span-12 sm:col-span-6"
			bind:value={form.fields.timeToProcess}
			error={form.errors.timeToProcess}
		></Form.TextInput>

		<Form.Tiptap
			title={m.question_assignment()}
			name="content"
			id="content"
			class="col-span-12"
			bind:value={form.fields.content}
			error={form.errors.content}
			enableFileUpload
			enableFileLink
		></Form.Tiptap>
	</div>
</div>
