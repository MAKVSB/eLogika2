<script lang="ts" module>
	export type ZonedDateTimeRange = {
		start?: ZonedDateTime;
		end?: ZonedDateTime;
	};
</script>

<script lang="ts">
	import { API } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import { getLocale } from '$lib/paraglide/runtime';
	import * as Form from '$lib/components/ui/form';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import {
		getLocalTimeZone,
		ZonedDateTime,
		now,
		parseAbsoluteToLocal
	} from '@internationalized/date';
	import type {
		TermDTO,
		TermsGetByIdResponse,
		TermsInsertRequest,
		TermsInsertResponse,
		TermsUpdateRequest,
		TermsUpdateResponse
	} from '$lib/api_types';
	import DateRangeField from '$lib/components/ui/date-range-field/date-range-field.svelte';
	import { Label } from '$lib/components/ui/label';
	import { TermsInsertRequestSchema } from '$lib/schemas_static.js';
	import { m } from '$lib/paraglide/messages.js';
	import type { DateRange } from 'bits-ui';

	let locale = getLocale();

	let courseId = $derived<string>(page.params.courseId);
	let itemId = $derived<string>(page.params.itemId);
	let termId = $derived<string>(page.params.termId);

	let { data } = $props();
	$effect(() => {
		if (data.term) {
			data.term.then((data) => setResult(data)).catch(() => {});
		} else {
			form.fields = defaultFormData();
		}
	});

	let helperFields: {
		active: ZonedDateTimeRange;
		signIn: ZonedDateTimeRange;
		signOut: ZonedDateTimeRange;
	} = $state({
		active: {
			start: now(getLocalTimeZone()),
			end: now(getLocalTimeZone()).add({ days: 7 })
		},
		signIn: {
			start: now(getLocalTimeZone()),
			end: now(getLocalTimeZone()).add({ days: 7 })
		},
		signOut: {
			start: now(getLocalTimeZone()),
			end: now(getLocalTimeZone()).add({ days: 7 })
		}
	});

	function helperFieldChanged() {
		if (helperFields.active.start) {
			form.fields.activeFrom = helperFields.active.start?.toAbsoluteString();
		} else {
			form.fields.activeFrom = '';
		}
		if (helperFields.active.end) {
			form.fields.activeTo = helperFields.active.end?.toAbsoluteString();
		} else {
			form.fields.activeTo = '';
		}

		if (helperFields.signIn.start) {
			form.fields.signInFrom = helperFields.signIn.start?.toAbsoluteString();
		} else {
			form.fields.signInFrom = '';
		}
		if (helperFields.signIn.end) {
			form.fields.signInTo = helperFields.signIn.end?.toAbsoluteString();
		} else {
			form.fields.signInTo = '';
		}

		if (helperFields.signOut.start) {
			form.fields.signOutFrom = helperFields.signOut.start?.toAbsoluteString();
		} else {
			form.fields.signOutFrom = '';
		}
		if (helperFields.signOut.end) {
			form.fields.signOutTo = helperFields.signOut.end?.toAbsoluteString();
		} else {
			form.fields.signOutTo = '';
		}
	}

	function defaultFormData(): TermDTO {
		return {
			id: 0,
			version: 0,
			courseId: Number(page.params.courseId),
			courseItemId: Number(page.params.itemId),
			name: '',

			activeFrom: new Date().toISOString(),
			activeTo: new Date().toISOString(),

			requiresSign: true,
			signInFrom: new Date().toISOString(),
			signInTo: new Date().toISOString(),

			signOutFrom: new Date().toISOString(),
			signOutTo: new Date().toISOString(),

			// OfflineTo:    new Date(), // TODO
			classroom: '',
			studentsMax: 0,
			studentsJoined: 0,
			tries: 1
		};
	}
	let form = $state(Form.createForm(TermsInsertRequestSchema, defaultFormData()));
	function addedValidation(formFields: TermDTO): Form.ErrorObject {
		const errors: Form.ErrorObject = {};

		if (new Date(formFields.activeFrom).getTime() > new Date(formFields.activeTo).getTime()) {
			errors.activeFrom = 'Time "from" must be smaller than time "to"';
		}

		if (new Date(formFields.signInFrom).getTime() > new Date(formFields.signInTo).getTime()) {
			errors.signInFrom = 'Time "from" must be smaller than time "to"';
		}

		if (new Date(formFields.signOutFrom).getTime() > new Date(formFields.signOutTo).getTime()) {
			errors.signOutFrom = 'Time "from" must be smaller than time "to"';
		}

		return errors;
	}

	function setResult(res: TermsGetByIdResponse | TermsInsertResponse | TermsUpdateResponse) {
		form.fields = res.data;

		helperFields.active = {
			start: parseAbsoluteToLocal(res.data.activeFrom),
			end: parseAbsoluteToLocal(res.data.activeTo)
		};

		helperFields.signIn = {
			start: parseAbsoluteToLocal(res.data.signInFrom),
			end: parseAbsoluteToLocal(res.data.signInTo)
		};

		helperFields.signOut = {
			start: parseAbsoluteToLocal(res.data.signOutFrom),
			end: parseAbsoluteToLocal(res.data.signOutTo)
		};

		console.log('Transfering 25');
		goto(String(res.data.id), {
			replaceState: true
		});
	}

	async function handleSubmit(): Promise<any> {
		let request;
		if (data.creating) {
			request = API.request<TermsInsertRequest, TermsInsertResponse>(
				`/api/v2/courses/${courseId}/items/${itemId}/terms`,
				{
					method: 'POST',
					body: form.fields
				}
			);
		} else {
			request = API.request<TermsUpdateRequest, TermsUpdateResponse>(
				`/api/v2/courses/${courseId}/items/${itemId}/terms/${termId}`,
				{
					method: 'PUT',
					body: form.fields
				}
			);
		}
		return request.then((res) => {
			setResult(res);
			if (data.creating) {
				toast.success('Created succesfully');
			} else {
				toast.success('Saved succesfully');
			}
		});
	}
</script>

<div class="flex flex-col gap-8 m-8">
	{#await data.term}
		<Pageloader></Pageloader>
	{:then staticResourceData}
		<div class="flex flex-row justify-between">
			<h1 class="mb-8 text-2xl">
				Term management:
				<b>
					{staticResourceData?.data?.name ?? 'New term'}
				</b>
			</h1>
		</div>
		<Form.Root bind:form onsubmit={handleSubmit} isCreating={data.creating} {addedValidation}>
			<div class="flex flex-col gap-4 p-2">
				<Form.TextInput
					title={m.term_name()}
					name="name"
					id="name"
					type="text"
					bind:value={form.fields.name}
					error={form.errors.name ?? ''}
				></Form.TextInput>
				<div class="flex flex-col gap-2">
					<Label for="active">{m.term_active_range()}</Label>
					<DateRangeField
						bind:value={helperFields.active as DateRange}
						{locale}
						granularity="minute"
						onValueChange={helperFieldChanged}
					></DateRangeField>
					{#if form.errors.activeFrom}
						<p class="text-sm text-red-500">{form.errors.activeFrom}</p>
					{/if}
					{#if form.errors.activeTo}
						<p class="text-sm text-red-500">{form.errors.activeTo}</p>
					{/if}
				</div>
				<div class="flex flex-col gap-2">
					<Label for="signIn">{m.term_signin_range()}</Label>
					<DateRangeField
						bind:value={helperFields.signIn as DateRange}
						{locale}
						granularity="minute"
						onValueChange={helperFieldChanged}
					></DateRangeField>
					{#if form.errors.signInFrom}
						<p class="text-sm text-red-500">{form.errors.signInFrom}</p>
					{/if}
					{#if form.errors.signInTo}
						<p class="text-sm text-red-500">{form.errors.signInTo}</p>
					{/if}
				</div>
				<Form.Checkbox
					title="Requires sign in"
					name="requiresSign"
					id="requiresSign"
					bind:value={form.fields.requiresSign}
					error={form.errors.requiresSign}
				></Form.Checkbox>
				<div class="flex flex-col gap-2">
					<Label for="signOut">{m.term_signout_range()}</Label>
					<DateRangeField
						bind:value={helperFields.signOut as DateRange}
						{locale}
						granularity="minute"
						onValueChange={helperFieldChanged}
					></DateRangeField>
					{#if form.errors.signOutFrom}
						<p class="text-sm text-red-500">{form.errors.signOutFrom}</p>
					{/if}
					{#if form.errors.signOutTo}
						<p class="text-sm text-red-500">{form.errors.signOutTo}</p>
					{/if}
				</div>
				<Form.TextInput
					title={m.term_classroom()}
					name="classroom"
					id="classroom"
					type="text"
					bind:value={form.fields.classroom}
					error={form.errors.classroom ?? ''}
				></Form.TextInput>
				<Form.TextInput
					title={m.term_maximumstudents()}
					name="studentsMax"
					id="studentsMax"
					type="number"
					bind:value={form.fields.studentsMax}
					error={form.errors.studentsMax ?? ''}
				></Form.TextInput>
				<Form.TextInput
					title={m.term_allowed_tries()}
					name="tries"
					id="tries"
					type="number"
					bind:value={form.fields.tries}
					error={form.errors.tries ?? ''}
				></Form.TextInput>
			</div>
		</Form.Root>
	{/await}
</div>
