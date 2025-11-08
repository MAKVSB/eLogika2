<script lang="ts">
	import type { Readable } from 'svelte/store';
	import { Editor } from 'svelte-tiptap';
	import { TiptapMenuItemType } from './menutype';

	// Shadcn
	import { cn, type WithElementRef } from '$lib/utils.js';

	// Custom elements
	import Button, { type TipTapMenuItemButton } from './button.svelte';
	import Toggle, { type TipTapMenuItemToggle } from './toggle.svelte';
	import Dropdown, { type TipTapMenuItemDropdown } from './dropdown.svelte';

	// Icons
	import UndoIcon from '@lucide/svelte/icons/undo';
	import RedoIcon from '@lucide/svelte/icons/redo';
	import BoldIcon from '@lucide/svelte/icons/bold';
	import ItalicIcon from '@lucide/svelte/icons/italic';
	import UnderlineIcon from '@lucide/svelte/icons/underline';
	import H1Icon from '@lucide/svelte/icons/heading-1';
	import H2Icon from '@lucide/svelte/icons/heading-2';
	import H3Icon from '@lucide/svelte/icons/heading-3';
	import PilcrowIcon from '@lucide/svelte/icons/pilcrow'; // Paragraph
	import StrikeIcon from '@lucide/svelte/icons/strikethrough';
	import ListIcon from '@lucide/svelte/icons/list';
	import OListIcon from '@lucide/svelte/icons/list-ordered';
	import SuperscriptIcon from '@lucide/svelte/icons/superscript';
	import SubscriptIcon from '@lucide/svelte/icons/subscript';
	import TableIcon from '@lucide/svelte/icons/table';
	import HighlighterIcon from '@lucide/svelte/icons/highlighter';
	import ImageIcon from '@lucide/svelte/icons/image';
	import BookIcon from '@lucide/svelte/icons/book';
	import PiIcon from '@lucide/svelte/icons/pi';
	import LinkIcon from '@lucide/svelte/icons/link';
	import AlignLeft from '@lucide/svelte/icons/align-left';
	import AlignRight from '@lucide/svelte/icons/align-right';
	import AlignCenter from '@lucide/svelte/icons/align-center';
	import AlignJustify from '@lucide/svelte/icons/align-justify';

	type TipTapMenuItemGroup = {
		type: TiptapMenuItemType.GROUP;
		items: (TipTapMenuItemDropdown | TipTapMenuItemButton | TipTapMenuItemToggle)[];
	};
	type TipTapMenu = (
		| TipTapMenuItemGroup
		| TipTapMenuItemDropdown
		| TipTapMenuItemButton
		| TipTapMenuItemToggle
	)[];

	let {
		ref = $bindable(null),
		editor = $bindable(),
		disabled = false,
		enableFileUpload,
		enableFileLink
	}: WithElementRef<{
		editor: Readable<Editor>;
		disabled?: boolean;
		enableFileUpload: boolean;
		enableFileLink: boolean;
	}> = $props();

	const isActive = (name: string | null, attrs = {}) => {
		if (editor) {
			if (name) {
				return $editor.isActive(name, attrs);
			} else {
				return $editor.isActive(attrs);
			}
		}
		return false;
	};

	const menuItems = $derived<TipTapMenu>([
		{
			type: TiptapMenuItemType.GROUP,
			items: [
				{
					type: TiptapMenuItemType.BUTTON,
					icon: UndoIcon,
					command: () => {
						$editor?.chain().focus().undo().run();
					},
					disabled: !$editor?.can().chain().focus().undo().run(),
					tooltip: 'Undo last change'
				},
				{
					type: TiptapMenuItemType.BUTTON,
					icon: RedoIcon,
					command: () => {
						$editor?.chain().focus().redo().run();
					},
					disabled: !$editor?.can().chain().focus().redo().run(),
					tooltip: 'Redo last change'
				}
			]
		},
		{
			type: TiptapMenuItemType.GROUP,
			items: [
				{
					type: TiptapMenuItemType.TOGGLE,
					icon: BoldIcon,
					command: () => {
						$editor.chain().focus().toggleBold().run();
					},
					active: $editor?.isActive('bold'),
					tooltip: 'Toggle bold text'
				},
				{
					type: TiptapMenuItemType.TOGGLE,
					icon: ItalicIcon,
					command: () => {
						$editor.chain().focus().toggleItalic().run();
					},
					active: $editor?.isActive('italic'),
					tooltip: 'Toggle italic text'
				},
				{
					type: TiptapMenuItemType.TOGGLE,
					icon: UnderlineIcon,
					command: () => {
						$editor.chain().focus().toggleUnderline().run();
					},
					active: $editor?.isActive('underline'),
					tooltip: 'Toggle underline text'
				},
				{
					type: TiptapMenuItemType.TOGGLE,
					icon: StrikeIcon,
					command: () => {
						$editor.chain().focus().toggleStrike().run();
					},
					active: $editor?.isActive('strike'),
					tooltip: 'Toggle strike text'
				},
				{
					type: TiptapMenuItemType.TOGGLE,
					icon: SubscriptIcon,
					command: () => {
						$editor.chain().focus().toggleSubscript().run();
					},
					active: $editor?.isActive('subscript'),
					tooltip: 'Toggle subscript text'
				},
				{
					type: TiptapMenuItemType.TOGGLE,
					icon: SuperscriptIcon,
					command: () => {
						$editor.chain().focus().toggleSuperscript().run();
					},
					active: $editor?.isActive('superscript'),
					tooltip: 'Toggle superscript text'
				},
				{
					type: TiptapMenuItemType.TOGGLE,
					icon: HighlighterIcon,
					command: () => {
						$editor.chain().focus().toggleHighlight().run();
					},
					active: $editor?.isActive('highlight'),
					tooltip: 'Highlight text'
				}
			]
		},
		{
			type: TiptapMenuItemType.DROPDOWN,
			active: isActive('paragraph') || isActive('heading'),
			tooltip: 'Change text heading/paragraph type',
			default: {
				type: TiptapMenuItemType.DROPDOWN_ITEM,
				icon: PilcrowIcon,
				disabled: !isActive('paragraph') && !isActive('heading')
			},
			items: [
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					icon: H1Icon,
					title: 'Heading level 1',
					title_class: 'text-2xl',
					command: () => {
						$editor?.chain().focus().toggleHeading({ level: 1 }).run();
					},
					active: isActive('heading', { level: 1 })
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					icon: H2Icon,
					title: 'Heading level 2',
					title_class: 'text-xl',
					command: () => {
						$editor?.chain().focus().toggleHeading({ level: 2 }).run();
					},
					active: isActive('heading', { level: 2 })
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					icon: H3Icon,
					title: 'Heading level 3',
					title_class: 'text-lg',
					command: () => {
						$editor?.chain().focus().toggleHeading({ level: 3 }).run();
					},
					active: isActive('heading', { level: 3 })
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					icon: PilcrowIcon,
					title: 'Paragraph',
					title_class: '',
					command: () => {
						$editor?.chain().focus().setParagraph().run();
					},
					active: isActive('paragraph')
				}
			]
		},
		{
			type: TiptapMenuItemType.TOGGLE,
			icon: LinkIcon,
			tooltip: 'Create or update a link',
			command: () => {
				const url = prompt('Link:', $editor.getAttributes('link').href);

				if (url === null) {
					return;
				}

				if (url === '') {
					$editor.chain().focus().extendMarkRange('link').unsetLink().run();
					return;
				}

				$editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run();
			},
			active: $editor.isActive('link')
		},
		{
			type: TiptapMenuItemType.DROPDOWN,
			active:
				isActive(null, { textAlign: 'left' }) ||
				isActive(null, { textAlign: 'right' }) ||
				isActive(null, { textAlign: 'center' }) ||
				isActive(null, { textAlign: 'justify' }),
			tooltip: 'Change text alignment',
			default: {
				type: TiptapMenuItemType.DROPDOWN_ITEM,
				icon: AlignLeft
			},
			items: [
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					icon: AlignLeft,
					command: () => {
						$editor?.chain().focus().toggleTextAlign('left').run();
					},
					active: isActive(null, { textAlign: 'left' })
				},

				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					icon: AlignRight,
					command: () => {
						$editor?.chain().focus().toggleTextAlign('right').run();
					},
					active: isActive(null, { textAlign: 'right' })
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					icon: AlignCenter,
					command: () => {
						$editor?.chain().focus().toggleTextAlign('center').run();
					},
					active: isActive(null, { textAlign: 'center' })
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					icon: AlignJustify,
					command: () => {
						$editor?.chain().focus().toggleTextAlign('justify').run();
					},
					active: isActive(null, { textAlign: 'justify' })
				}
			]
		},
		{
			type: TiptapMenuItemType.DROPDOWN,
			active: isActive('bulletList') || isActive('orderedList'),
			tooltip: 'Add/Remove list',
			default: {
				type: TiptapMenuItemType.DROPDOWN_ITEM,
				title: 'Bullet list',
				icon: ListIcon
			},
			items: [
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Bullet list',
					icon: ListIcon,
					command: () => {
						$editor?.chain().focus().toggleBulletList().run();
					},
					active: isActive('bulletList')
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Ordered List',
					icon: OListIcon,
					command: () => {
						$editor?.chain().focus().toggleOrderedList().run();
					},
					active: isActive('orderedList')
				}
			]
		},
		{
			type: TiptapMenuItemType.GROUP,
			items: [
				{
					type: TiptapMenuItemType.TOGGLE,
					icon: BookIcon,
					command: () => {
						if (isActive('detailsContent') || isActive('detailsSummary')) {
							$editor?.chain().focus().unsetDetails().run();
						} else {
							$editor?.chain().focus().setDetails().run();
						}
					},
					tooltip: 'Wrap detail',
					active: isActive('detailsContent') || isActive('detailsSummary')
				}
			]
		},
		{
			type: TiptapMenuItemType.DROPDOWN,
			active: isActive('table'),
			tooltip: 'Add or edit table',
			show_title: true,
			default: {
				type: TiptapMenuItemType.DROPDOWN_ITEM,
				icon: TableIcon,
				title: 'Table'
			},
			items: [
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Insert table',
					command: () => {
						$editor?.chain().focus().insertTable({ rows: 3, cols: 3 }).run();
					},
					disabled: isActive('table')
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Add column before',
					command: () => {
						$editor?.chain().focus().addColumnBefore().run();
					},
					disabled: !isActive('table')
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Add column after',
					command: () => {
						$editor?.chain().focus().addColumnBefore().run();
					},
					disabled: !isActive('table')
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Toggle header column',
					command: () => {
						$editor?.chain().focus().toggleHeaderColumn().run();
					},
					disabled: !isActive('table')
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Delete column',
					command: () => {
						$editor?.chain().focus().deleteColumn().run();
					},
					disabled: !isActive('table')
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Add row before',
					command: () => {
						$editor?.chain().focus().addRowBefore().run();
					},
					disabled: !isActive('table')
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Add row after',
					command: () => {
						$editor?.chain().focus().addRowAfter().run();
					},
					disabled: !isActive('table')
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Toggle header row',
					command: () => {
						$editor?.chain().focus().toggleHeaderRow().run();
					},
					disabled: !isActive('table')
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Delete row',
					command: () => {
						$editor?.chain().focus().deleteRow().run();
					},
					disabled: !isActive('table')
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Delete table',
					command: () => {
						$editor?.chain().focus().deleteTable().run();
					},
					disabled: !isActive('table')
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Merge or split',
					command: () => {
						$editor?.chain().focus().mergeOrSplit().run();
					},
					disabled: !isActive('table')
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Fix tables',
					command: () => {
						$editor?.chain().focus().fixTables().run();
					},
					disabled: !isActive('table')
				}
			]
		},
		{
			type: TiptapMenuItemType.DROPDOWN,
			active: isActive('table'),
			tooltip: 'Add math',
			show_title: true,
			default: {
				type: TiptapMenuItemType.DROPDOWN_ITEM,
				icon: PiIcon,
				title: 'Math'
			},
			items: [
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Insert inline math',
					command: () => {
						const latex = prompt('Enter block math expression:', '');

						if (latex === null) {
							return;
						}

						if (latex === '') {
							$editor.chain().deleteInlineMath().focus().run();
						}

						$editor.chain().insertInlineMath({ latex }).focus().run();
					},
					disabled: !$editor.state.selection.empty
				},
				{
					type: TiptapMenuItemType.DROPDOWN_ITEM,
					title: 'Insert block math',
					command: () => {
						const latex = prompt('Enter block math expression:', '');

						if (latex === null) {
							return;
						}

						if (latex === '') {
							$editor.chain().deleteBlockMath().focus().run();
							return;
						}

						$editor.chain().insertBlockMath({ latex }).focus().run();
					},
					disabled: !$editor.state.selection.empty
				}
			]
		},
		{
			type: TiptapMenuItemType.GROUP,
			items: [
				{
					type: TiptapMenuItemType.BUTTON,
					title: 'Link file url',
					icon: ImageIcon,
					command: () => {
						const url = prompt('Enter url:', '');

						if (url === null) {
							return;
						}

						$editor
							?.chain()
							.focus()
							.setImage({
								mode: 'url',
								type: 'unknown',
								src: url,
								originalFilename: url
							})
							.run();
					},
					tooltip: 'Add image',
					disabled: !enableFileLink
				}
			]
		},
		{
			type: TiptapMenuItemType.GROUP,
			items: [
				// {
				// 	type: TiptapMenuItemType.TOGGLE,
				// 	title: 'Code block',
				// 	icon: ImageIcon,
				// 	command: () => {
				// 		$editor.chain().focus().toggleCodeBlock().run();
				// 	},
				// 	tooltip: 'Add code',
				// 	disabled: !enableFileLink,
				// 	active: isActive('codeBlock') && !isActive('codeBlock', {language: "latex"})
				// },
				{
					type: TiptapMenuItemType.TOGGLE,
					title: 'Latex block',
					icon: ImageIcon,
					command: () => {
						$editor.chain().focus().toggleCodeBlock({language: "latex"}).run();
					},
					tooltip: 'Add latex block',
					disabled: !enableFileLink,
					active: isActive('codeBlock', {language: "latex"})
				}
			]
		}
	]);
</script>

<div class={cn('flex overflow-x-scroll rounded-md rounded-b-none border p-3 overflow-auto')}>
	<div
		class={cn('flex rounded-md rounded-b-none', disabled ? 'cursor-not-allowed opacity-50' : '')}
	>
		{#each menuItems as menuItem}
			{#if menuItem.type == TiptapMenuItemType.GROUP}
				{#each menuItem.items as menuItem2}
					{#if menuItem2.type == TiptapMenuItemType.TOGGLE}
						<Toggle {...menuItem2} disabled={disabled || menuItem2.disabled}></Toggle>
					{:else if menuItem2.type == TiptapMenuItemType.BUTTON}
						<Button {...menuItem2} disabled={disabled || menuItem2.disabled}></Button>
					{:else if menuItem2.type == TiptapMenuItemType.DROPDOWN}
						<Dropdown {...menuItem2} disabled={disabled || menuItem2.disabled}></Dropdown>
					{/if}
				{/each}
			{:else if menuItem.type == TiptapMenuItemType.DROPDOWN}
				<Dropdown {...menuItem} disabled={disabled || menuItem.disabled}></Dropdown>
			{:else if menuItem.type == TiptapMenuItemType.TOGGLE}
				<Toggle {...menuItem} disabled={disabled || menuItem.disabled}></Toggle>
			{:else if menuItem.type == TiptapMenuItemType.BUTTON}
				<Button {...menuItem} disabled={disabled || menuItem.disabled}></Button>
			{/if}
		{/each}
	</div>
</div>
