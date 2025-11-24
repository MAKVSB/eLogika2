<script lang="ts">
	import { type TokenCreateRequest, type TokenCreateResponse } from '$lib/api_types';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { API } from '$lib/services/api.svelte';
	import { Label } from '$lib/components/ui/label/index.js';
	import { toast } from 'svelte-sonner';
	import { invalidateAll } from '$app/navigation';
	import * as Form from '$lib/components/ui/form';
	import { m } from '$lib/paraglide/messages';
	import DateField from '$lib/components/ui/date-field/date-field.svelte';
	import { getLocalTimeZone, now, type ZonedDateTime } from '@internationalized/date';
	import { TokenCreateRequestSchema } from '$lib/schemas';
	import { getLocale } from '$lib/paraglide/runtime';
	import Button from '$lib/components/ui/button/button.svelte';

	let locale = getLocale();

	let {
		openState = $bindable()
	}: {
		openState: boolean;
	} = $props();

	let response: TokenCreateResponse | null = $state(null);

	const defaultFormData: TokenCreateRequest = {
		name: '',
		expiresAt: new Date()
	};
	let form = $state(Form.createForm(TokenCreateRequestSchema, defaultFormData));

	let helperFieldExpiresAt: ZonedDateTime = $state(now(getLocalTimeZone()).add({ days: 7 }));

	function helperFieldChanged() {
		if (helperFieldExpiresAt) {
			form.fields.expiresAt = helperFieldExpiresAt?.toAbsoluteString();
		} else {
			form.fields.expiresAt = '';
		}
	}

	async function handleSubmit(): Promise<any> {
		let request = API.request<TokenCreateRequest, TokenCreateResponse>(
			`/api/v2/users/self/tokens`,
			{
				method: 'POST',
				body: form.fields
			},
			fetch
		);

		return request
			.then((res) => {
				response = res;
			})
			.catch(() => {});
	}
</script>

<Dialog.Content class="max-h-full w-300 overflow-scroll sm:max-h-[90%] sm:max-w-[90%]">
	{#if response}
		<Dialog.Header>
			<Dialog.Title>New api token</Dialog.Title>
		</Dialog.Header>

		<Form.Root>
			<Form.TextInput
				title={m.token_name()}
				name="name"
				id="name"
				type="text"
				bind:value={(response as TokenCreateResponse).name}
				error=""
				disabled
			></Form.TextInput>

			<Form.TextInput
				title={m.token_value()}
				name="value"
				id="value"
				type="value"
				bind:value={(response as TokenCreateResponse).value}
				error=""
			></Form.TextInput>
		</Form.Root>

		<Button
			onclick={() => {
				invalidateAll();
				openState = false;
			}}
		>
			Close
		</Button>
	{:else}
		<Dialog.Header>
			<Dialog.Title>Create new API token</Dialog.Title>
		</Dialog.Header>

		<Form.Root bind:form onsubmit={handleSubmit} isCreating={true}>
			<Form.TextInput
				title={m.token_name()}
				name="name"
				id="name"
				type="text"
				bind:value={form.fields.name}
				error={form.errors.name}
			></Form.TextInput>

			<div class="flex flex-col gap-2">
				<Label for="active">{m.token_expires_at()}</Label>
				<DateField
					bind:value={helperFieldExpiresAt}
					{locale}
					granularity="minute"
					onValueChange={helperFieldChanged}
				></DateField>
				{#if form.errors.activeFrom}
					<p class="text-sm text-red-500">{form.errors.activeFrom}</p>
				{/if}
				{#if form.errors.activeTo}
					<p class="text-sm text-red-500">{form.errors.activeTo}</p>
				{/if}
			</div>
		</Form.Root>
	{/if}
</Dialog.Content>
