<script lang="ts">
	import { API } from '$lib/services/api.svelte';
	import { onMount } from 'svelte';
	import GlobalState from '$lib/shared.svelte';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';

	onMount(() => {
		const params = new URLSearchParams(window.location.search);
		const provider = params.get('provider');

		if (provider == 'VSBCAS') {
			const ticket = params.get('ticket');

			if (ticket != '') {
				API.request(`/api/v2/auth/login/sso/callback`, {
					method: 'POST',
					body: {
						provider,
						ticket
					}
				})
					.then((res) => {
						GlobalState.accessToken = res.accessToken;
					})
					.catch(() => {
						const errMsg = 'Failed to validate token, please try again';
						toast.error(errMsg);
						console.error(errMsg);
						setTimeout(() => {
							goto(base + '/login');
						}, 2000);
					});
			} else {
				const errMsg = 'Auth token not found, redirecting to login page';
				toast.error(errMsg);
				console.error(errMsg);
				setTimeout(() => {
					goto(base + '/login');
				}, 2000);
				
			}
		}
	});
</script>
