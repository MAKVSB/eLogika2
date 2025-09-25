<script lang="ts">
	import { API, ApiError } from '$lib/services/api.svelte';
	import { page } from '$app/state';
	import * as Form from '$lib/components/ui/form';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';
	import Pageloader from '$lib/components/ui/loader/pageloader.svelte';
	import { m } from '$lib/paraglide/messages';
	import {
		IdentityProviderEnum,
		UserTypeEnum,
		type UserDTO,
		type UserGetByIdResponse,
		type UserInsertRequest,
		type UserInsertResponse,
		type UserUpdateRequest,
		type UserUpdateResponse
	} from '$lib/api_types';
	import { enumToOptions } from '$lib/utils';
	import GlobalState from '$lib/shared.svelte';
	import type { ErrorObject } from '$lib/components/ui/form/types';
	import { UserInsertRequestSchema } from '$lib/schemas.js';

	let { data } = $props();

	$effect(() => {
		if (data.question) {
			data.question
				.then((data) => setResult(data))
				.catch(() => {});
		} else {
			form.fields = defaultFormData;
		}
	});

	const defaultFormData: UserDTO = {
		id: 0,
		version: 0,
		firstName: '',
		familyName: '',
		username: '',
		email: '',
		notification: {
			discord: {
				level: {
					results: false,
					messages: false,
					terms: false
				},
				userId: ''
			},
			email: {
				level: {
					results: true,
					messages: true,
					terms: true
				}
			},
			push: {
				level: {
					results: true,
					messages: true,
					terms: true
				}
			}
		},
		type: UserTypeEnum.NORMAL,
		identityProvider: IdentityProviderEnum.INT,
		courses: []
	};
	let form = $state(Form.createForm(UserInsertRequestSchema, defaultFormData));

	function setResult(res: UserGetByIdResponse | UserInsertResponse | UserUpdateResponse) {
		form.fields = res.data;
		console.log("Transfering 34")
		goto(String(res.data.id), {
			replaceState: true
		});
	}

	async function handleSubmit(): Promise<any> {
		let request;
		if (data.creating) {
			request = API.request<UserInsertRequest, UserInsertResponse>(`/api/v2/users`, {
				method: 'POST',
				body: form.fields
			});
		} else {
			request = API.request<UserUpdateRequest, UserUpdateResponse>(
				`/api/v2/users/${page.params.id}`,
				{
					method: 'PUT',
					body: form.fields
				}
			);
		}

		return request.then((res) => setResult(res));
	}
</script>

<div class="m-8">
	{#await data.question}
		<Pageloader></Pageloader>
	{:then staticResourceData}
		<div class="flex flex-row justify-between">
			<h1 class="mb-8 text-2xl">
				Admin user management:
				<b>
					{staticResourceData?.data?.username ?? 'New user'}
				</b>
			</h1>
		</div>
		<Form.Root bind:form onsubmit={handleSubmit} isCreating={data.creating}>
			<div class="flex flex-col gap-4 p-2">
				<div
					class="grid {GlobalState.loggedUser?.type == UserTypeEnum.ADMIN
						? 'grid-cols-3'
						: 'grid-cols-2'} gap-4"
				>
					<Form.TextInput
						title={m.user_first_name()}
						name="firstName"
						id="firstName"
						type="text"
						bind:value={form.fields.firstName}
						error={form.errors.firstName}
					></Form.TextInput>
					<Form.TextInput
						title={m.user_family_name()}
						name="familyName"
						id="familyName"
						type="text"
						bind:value={form.fields.familyName}
						error={form.errors.firstName}
					></Form.TextInput>
					{#if GlobalState.loggedUser?.type == UserTypeEnum.ADMIN}
						<Form.SingleSelect
							title="User type"
							name="type"
							id="type"
							options={enumToOptions(UserTypeEnum, m.user_type_enum)}
							error={form.errors.type}
							bind:value={form.fields.type}
						></Form.SingleSelect>
					{/if}
				</div>
				<div class="grid grid-cols-3 gap-4">
					<Form.TextInput
						title={m.user_username()}
						name="username"
						id="username"
						type="text"
						bind:value={form.fields.username}
						error={form.errors.username}
					></Form.TextInput>
					<Form.TextInput
						title={m.user_email()}
						name="email"
						id="email"
						type="email"
						class="col-span-2"
						bind:value={form.fields.email}
						error={form.errors.email}
					></Form.TextInput>
				</div>
				<div class="grid grid-cols-1 gap-4">
					<h3 class="text-xl">Notification settings:</h3>

					<div class="px-4">
						<h4 class="text-lg">Email:</h4>
						<div class="grid grid-cols-3 gap-4 sm:grid-cols-6">
							<Form.Checkbox
								title="Results"
								name="email-results"
								id="email-results"
								value={form.fields.notification.email.level.results}
								error={(
									((form.errors.notification as ErrorObject)?.Email as ErrorObject)
										?.LevTermsel as ErrorObject
								)?.Results}
							></Form.Checkbox>
							<Form.Checkbox
								title="Messages"
								name="email-messages"
								id="email-messages"
								value={form.fields.notification.email.level.messages}
								error={(
									((form.errors.notification as ErrorObject)?.Email as ErrorObject)
										?.LevTermsel as ErrorObject
								)?.Messages}
							></Form.Checkbox>
							<Form.Checkbox
								title="Terms"
								name="email-terms"
								id="email-terms"
								value={form.fields.notification.email.level.terms}
								error={(
									((form.errors.notification as ErrorObject)?.Email as ErrorObject)
										?.LevTermsel as ErrorObject
								)?.Terms}
							></Form.Checkbox>
						</div>
					</div>
					<div class="px-4">
						<h4 class="text-lg">Push:</h4>
						<div class="grid grid-cols-3 gap-4 sm:grid-cols-6">
							<Form.Checkbox
								title="Results"
								name="push-results"
								id="push-results"
								value={form.fields.notification.push.level.results}
								error={(
									((form.errors.notification as ErrorObject)?.Push as ErrorObject)
										?.LevTermsel as ErrorObject
								)?.Results}
							></Form.Checkbox>
							<Form.Checkbox
								title="Messages"
								name="push-messages"
								id="push-messages"
								value={form.fields.notification.push.level.messages}
								error={(
									((form.errors.notification as ErrorObject)?.Push as ErrorObject)
										?.LevTermsel as ErrorObject
								)?.Messages}
							></Form.Checkbox>
							<Form.Checkbox
								title="Terms"
								name="push-terms"
								id="push-terms"
								value={form.fields.notification.push.level.terms}
								error={(
									((form.errors.notification as ErrorObject)?.Push as ErrorObject)
										?.LevTermsel as ErrorObject
								)?.Terms}
							></Form.Checkbox>
						</div>
					</div>
					<div class="px-4">
						<h4 class="text-lg">Discord:</h4>
						<div class="grid grid-cols-3 gap-4 sm:grid-cols-6">
							<Form.Checkbox
								title="Results"
								name="discord-results"
								id="discord-results"
								value={form.fields.notification.discord.level.results}
								error={(
									((form.errors.notification as ErrorObject)?.Discord as ErrorObject)
										?.Level as ErrorObject
								)?.Results}
							></Form.Checkbox>
							<Form.Checkbox
								title="Messages"
								name="discord-messages"
								id="discord-messages"
								value={form.fields.notification.discord.level.messages}
								error={(
									((form.errors.notification as ErrorObject)?.Discord as ErrorObject)
										?.Level as ErrorObject
								)?.Messages}
							></Form.Checkbox>
							<Form.Checkbox
								title="Terms"
								name="discord-terms"
								id="discord-terms"
								value={form.fields.notification.discord.level.terms}
								error={(
									((form.errors.notification as ErrorObject)?.Discord as ErrorObject)
										?.LevTermsel as ErrorObject
								)?.Messages}
							></Form.Checkbox>
							<Form.TextInput
								title="Discord user id"
								name="discord-userid"
								id="discord-userid"
								type="number"
								class="col-span-3"
								bind:value={form.fields.notification.discord.userId}
								error={((form.errors.notification as ErrorObject)?.Discord as ErrorObject)
									?.UserID as ErrorObject}
							></Form.TextInput>
						</div>
					</div>
				</div>
			</div>
		</Form.Root>
	{/await}
</div>
