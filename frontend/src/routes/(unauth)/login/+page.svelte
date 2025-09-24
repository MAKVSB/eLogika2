<script lang="ts">
	import { z } from 'zod';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import Logo from '$lib/images/logo.svg';
	import * as Form from '$lib/components/ui/form';
	import { API } from '$lib/services/api.svelte';
	import { toast } from 'svelte-sonner';
	import type { ErrorObject } from '$lib/components/ui/form/types';
	import { m } from '$lib/paraglide/messages';
	import { LoginRequestSchema } from '$lib/schemas';
	import type { LoginRequest } from '$lib/api_types';


	const defaultFormData: LoginRequest = {
		email: "",
		password: "",
	}
	let form = Form.createForm(LoginRequestSchema, defaultFormData)

	async function handleSubmit(): Promise<any> {
		return API.login(form.fields)
			.then((res) => {
				if (res) {
					toast.success('Login successfull');
				} else {
					toast.error('Login failed');
				}
			})
			.catch(() => {
				toast.error('Fatal error during login');
			});
	}

	async function handleSubmitSSO(provider: string) {
		await API.request(`/api/v2/auth/login/sso`, {
			method: "POST",
			body: {
				provider
			}
		})
			.then((res) => {
				if (res.redirectUrl) {
					window.location = res.redirectUrl
				}
			})
			.catch(() => {
				toast.error('Fatal error during login');
			});
	}
</script>

{#snippet loginButtons()}
	<div class="flex flex-col w-full gap-4">
		<Form.Button text="Login" textSubmiting="Logging in" isSubmitting={form.isSubmitting}
		></Form.Button>

		<Button variant="outline" class="w-full" onclick={() => {
			handleSubmitSSO("VSB-CAS")
		}}>{m.login_button_sso()}</Button>
	</div>
{/snippet}

<div class="flex items-center justify-center w-full h-screen px-4">
	<Card.Root class="w-full max-w-sm mx-auto">
		<Card.Header>
			<Card.Title class="flex text-2xl">
				<img src={Logo} alt="eLogika logo" class="h-8 me-3" />
				<span
					class="self-center text-xl font-semibold whitespace-nowrap sm:text-2xl dark:text-white"
				>
					eLogika
				</span>
			</Card.Title>
		</Card.Header>

		<Card.Content>
			<Form.Root
				bind:form
				onsubmit={handleSubmit}
				additionalButtons={loginButtons}
				hideDefaultbutton={true}
			>
				<div class="grid gap-4">
					<Form.TextInput
						title={m.user_email()}
						name="email"
						id="email"
						type="email"
						placeholder="login@vsb.cz"
						required={!form.schema.shape.email.isOptional() &&
							!form.schema.shape.email.isNullable()}
						bind:value={form.fields.email}
						error={form.errors.email}
					></Form.TextInput>
					<Form.TextInput
						title={m.user_password()}
						name="password"
						id="password"
						type="password"
						required={!form.schema.shape.password.isOptional() &&
							!form.schema.shape.password.isNullable()}
						bind:value={form.fields.password}
						error={form.errors.password}
					></Form.TextInput>
				</div>
			</Form.Root>
		</Card.Content>
	</Card.Root>
</div>
