<script lang="ts">
	import * as Item from '$lib/components/ui/item/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { base } from '$app/paths';
	import { m } from '$lib/paraglide/messages';
</script>

<div class="flex flex-col gap-4 m-8">
	<h1 class="text-2xl">eLogika 2.0 {m.recognizer()}:</h1>

	<div class="flex flex-col w-full gap-6 md:flex-row">
		<Item.Root variant="outline">
			<Item.Content>
				<Item.Title>Windows</Item.Title>
                <Item.Description>{m.recognizer_download_description({os: "windows"})}</Item.Description>
			</Item.Content>
			<Item.Actions>
				<Button variant="outline" size="sm" target="_blank" href="{base}/eLogikaScanner.exe">{m.download()}</Button>
			</Item.Actions>
		</Item.Root>
		<Item.Root variant="outline">
			<Item.Content>
				<Item.Title>Linux</Item.Title>
                <Item.Description>{m.recognizer_download_description({os: "linux"})}</Item.Description>
			</Item.Content>
			<Item.Actions>
				<Button variant="outline" size="sm" target="_blank" href="{base}/eLogikaScanner">{m.download()}</Button>
			</Item.Actions>
		</Item.Root>
		<Item.Root variant="outline">
			<Item.Content>
				<Item.Title>{m.recognizer_source()}</Item.Title>
                <Item.Description>{m.recognizer_source_description()}</Item.Description>
			</Item.Content>
			<Item.Actions>
				<Button variant="outline" size="sm" target="_blank" href="https://bitbucket.org/makvsb/sp/src/main/">{m.view()}</Button>
			</Item.Actions>
		</Item.Root>
	</div>
</div>
