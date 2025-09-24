<script lang="ts" module>
	export type ZonedDateTimeRange = {
		start?: ZonedDateTime;
		end?: ZonedDateTime;
	};
</script>

<script lang="ts">
	import { API, ApiError } from '$lib/services/api.svelte';
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
	import { TermsInsertRequestSchema } from '$lib/schemas';
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
		form.fields.activeFrom = helperFields.active.start?.toDate() ?? null!;
		form.fields.activeTo = helperFields.active.end?.toDate() ?? null!;

		form.fields.signInFrom = helperFields.signIn.start?.toDate() ?? null!;
		form.fields.signInTo = helperFields.signIn.end?.toDate() ?? null!;

		form.fields.signOutFrom = helperFields.signOut.start?.toDate() ?? null!;
		form.fields.signOutTo = helperFields.signOut.end?.toDate() ?? null!;
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

		form.fields.activeFrom = new Date(res.data.activeFrom);
		form.fields.activeTo = new Date(res.data.activeTo);
		helperFields.active = {
			start: parseAbsoluteToLocal(res.data.activeFrom),
			end: parseAbsoluteToLocal(res.data.activeTo)
		};

		helperFields.signIn = {
			start: parseAbsoluteToLocal(res.data.signInFrom),
			end: parseAbsoluteToLocal(res.data.signInTo)
		};
		form.fields.signInFrom = new Date(res.data.signInFrom);
		form.fields.signInTo = new Date(res.data.signInTo);

		helperFields.signOut = {
			start: parseAbsoluteToLocal(res.data.signOutFrom),
			end: parseAbsoluteToLocal(res.data.signOutTo)
		};
		form.fields.signOutFrom = new Date(res.data.signOutFrom);
		form.fields.signOutTo = new Date(res.data.signOutTo);

		console.log("Transfering 25")
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
		return request.then((res) => setResult(res));
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
					title="Name"
					name="name"
					id="name"
					type="text"
					bind:value={form.fields.name}
					error={form.errors.name ?? ''}
				></Form.TextInput>
				<div class="flex flex-col gap-2">
					<Label for="active">Active range</Label>
					<DateRangeField
						bind:value={helperFields.active}
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
					<Label for="signIn">Sign in range</Label>
					<DateRangeField
						bind:value={helperFields.signIn}
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
					<Label for="signOut">Sign out range</Label>
					<DateRangeField
						bind:value={helperFields.signOut}
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
					title="Classroom"
					name="classroom"
					id="classroom"
					type="text"
					bind:value={form.fields.classroom}
					error={form.errors.classroom ?? ''}
				></Form.TextInput>
				<Form.TextInput
					title="Maximum number of students"
					name="studentsMax"
					id="studentsMax"
					type="number"
					bind:value={form.fields.studentsMax}
					error={form.errors.studentsMax ?? ''}
				></Form.TextInput>
				<Form.TextInput
					title="Number of allowed tries"
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
