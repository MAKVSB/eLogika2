<script lang="ts">
	import { API } from '$lib/services/api.svelte';
	import { onMount } from 'svelte';
	import GlobalState from '$lib/shared.svelte';

	onMount(() => {
		const params = new URLSearchParams(window.location.search);
		const provider = params.get('provider');

		if (provider == 'VSBCAS') {
			const ticket = params.get('ticket');

			if (ticket) {
				API.request(`/api/v2/auth/login/sso/callback`, {
					method: 'POST',
					body: {
						provider: 'VSB-CAS',
						ticket
					}
				})
					.then((res) => {
						GlobalState.accessToken = res.accessToken;
					})
					.catch(() => {});
			}
		}
	});
</script>
