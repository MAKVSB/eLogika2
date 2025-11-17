<script lang="ts">
	import DataTable from '$lib/components/ui/data-table/data-table-component.svelte';
	import { columns, filters } from './schema';
	import { type InitialTableState } from '@tanstack/table-core';
	import { buttonVariants } from '$lib/components/ui/button/button.svelte';
	import type { UserAPiTokenDTO } from '$lib/api_types';
	import { API } from '$lib/services/api.svelte';
	import { invalidateAll } from '$app/navigation';
	import CreateDialog from './CreateDialog.svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { m } from '$lib/paraglide/messages';

	let {
		tokens
	}: {
		tokens: any[]
	} = $props();
	let isLoading: boolean = $state(true);

	let dialogOpen = $state(false)

	$effect(() => {
		rowItems = tokens;
		rowCount = tokens.length;
		isLoading = false;
	});

	const actionsColumn = columns.find((c) => c.uniqueId == 'actions');
	if (actionsColumn) {
		actionsColumn.meta = {
			...(actionsColumn.meta ?? {}),
			clickEventHandler: async (event: string, id: number) => {
				switch (event) {
					case 'revoke':
						if (!confirm(m.user_token_revoke_confirm())) {
							return;
						}
						API.request<any, Blob>(
							`/api/v2/users/self/tokens/${id}`,
							{
								method: 'DELETE'
							},
							fetch
						)
							.then((res) => {
								invalidateAll();
							})
							.catch(() => {});
						break;
				}

				return true;
			}
		};
	}

	let rowItems: UserAPiTokenDTO[] = $state([]);
	let rowCount: number = $state(0);
	let initialState: InitialTableState = $state({});
</script>


<div class="flex flex-row justify-between">
	<h3 class="text-xl">{m.user_tokens()}:</h3>
	<div class="flex gap-2">
		<Dialog.Root bind:open={dialogOpen}>
			<Dialog.Trigger class={buttonVariants({ variant: 'default' })} type="button">
				{m.user_token_create_btn()}
			</Dialog.Trigger>
			{#if dialogOpen}
				<CreateDialog bind:openState={dialogOpen}></CreateDialog>
			{/if}
		</Dialog.Root>
	</div>
</div>
{#if !isLoading}
	<DataTable data={rowItems} {columns} {filters} {initialState} {rowCount} />
{/if}

