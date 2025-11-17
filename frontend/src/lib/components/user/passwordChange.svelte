<script lang="ts">
	import { API } from '$lib/services/api.svelte';
	import * as Form from '$lib/components/ui/form';
	import { m } from '$lib/paraglide/messages';
	import { IdentityProviderEnum, UserTypeEnum, type UserChangePassRequest } from '$lib/api_types';
	import GlobalState from '$lib/shared.svelte';
	import { UserChangePassRequestSchemaStatic } from '$lib/schemas_static.js';
	import { Button } from '../ui/button';
	import { toast } from 'svelte-sonner';

	let {
		id,
		identityProvider
	}: {
		id: number;
		identityProvider: IdentityProviderEnum;
	} = $props();

	let emptyPassForm = {
		oldPassword: '',
		newPassword: '',
		newPasswordRep: ''
	};
	let passwordForm = $state(Form.createForm(UserChangePassRequestSchemaStatic, emptyPassForm));

	async function handleSubmitPassword(): Promise<any> {
		let request = API.request<UserChangePassRequest, null>(`/api/v2/users/self/password`, {
			method: 'PUT',
			body: {
				...passwordForm.fields,
				generate: false
			}
		});

		passwordForm.fields = emptyPassForm;

		return request.then(() => {});
	}

	async function handleGeneratePassword(): Promise<any> {
		let request = API.request<UserChangePassRequest, null>(`/api/v2/users/self/password`, {
			method: 'PUT',
			body: {
				generate: true
			}
		});

		return request.then(() => {
			toast.info('New password has been sent to email of the user');
		});
	}
</script>

{#if identityProvider == IdentityProviderEnum.INT}
	<div>
		<div class="flex flex-row justify-between">
			<h1 class="text-2xl">Change password:</h1>
			<div>
				{#if id != GlobalState.loggedUser?.id}
					<Button onclick={() => handleGeneratePassword()}>Generovat nové heslo</Button>
				{/if}
			</div>
		</div>
		{#if id == GlobalState.loggedUser?.id}
			<Form.Root bind:form={passwordForm} onsubmit={handleSubmitPassword} isCreating={false}>
				<div class="flex flex-col gap-4 p-2">
					<div
						class="grid {GlobalState.loggedUser?.type == UserTypeEnum.ADMIN
							? 'grid-cols-3'
							: 'grid-cols-2'} gap-4"
					>
						<Form.TextInput
							title={m.user_password_reset_current()}
							name="oldPassword"
							id="oldPassword"
							type="password"
							bind:value={passwordForm.fields.oldPassword}
							error={passwordForm.errors.oldPassword}
						></Form.TextInput>
						<Form.TextInput
							title={m.user_password_reset_new()}
							name="newPassword"
							id="newPassword"
							type="password"
							bind:value={passwordForm.fields.newPassword}
							error={passwordForm.errors.newPassword}
						></Form.TextInput>
						<Form.TextInput
							title={m.user_password_reset_newrepeat()}
							name="newPasswordRep"
							id="newPasswordRep"
							type="password"
							bind:value={passwordForm.fields.newPasswordRep}
							error={passwordForm.errors.newPasswordRep}
						></Form.TextInput>
					</div>
				</div>
			</Form.Root>
		{/if}
	</div>
{/if}
