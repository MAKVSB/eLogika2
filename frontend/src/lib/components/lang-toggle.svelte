<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { setLocale } from '$lib/paraglide/runtime';
	import { m } from '$lib/paraglide/messages';
	import * as DropdownMenu from './ui/dropdown-menu';
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<Button {...props} variant="outline">{m.language_change_button()}</Button>
		{/snippet}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content class="w-56" align="start">
		<DropdownMenu.Label>{m.language()}</DropdownMenu.Label>
		<DropdownMenu.Group>
			<DropdownMenu.Item onclick={() => setLocale('en')}>English</DropdownMenu.Item>
			<DropdownMenu.Item onclick={() => setLocale('cs')}>Česky</DropdownMenu.Item>
		</DropdownMenu.Group>
	</DropdownMenu.Content>
</DropdownMenu.Root>
