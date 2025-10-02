<script lang="ts">
	import 'katex/dist/katex.min.css';

	import { onDestroy, onMount } from 'svelte';
	import { createLowlight, common as LLCommon } from 'lowlight';
	import { createEditor, Editor, EditorContent } from 'svelte-tiptap';

	import StarterKit from '@tiptap/starter-kit';
	import Subscript from '@tiptap/extension-subscript';
	import Superscript from '@tiptap/extension-superscript';
	import { TableKit } from '@tiptap/extension-table';
	import Hightlight from '@tiptap/extension-highlight';
	import CustomInlineMath from './plugins/InlineMath';
	import CustomBlockMath from './plugins/BlockMath';
	import FileHandler from '@tiptap/extension-file-handler';
	import TipTapImage from './plugins/image';
	import { Details, DetailsContent, DetailsSummary } from '@tiptap/extension-details';
	import Typography from '@tiptap/extension-typography';
	import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight';
	import TextAlign from '@tiptap/extension-text-align';

	import MathReplacer from './plugins/MathReplacer';
	import CharsReplacer from './plugins/CharsReplacer';

	interface TipTapAttributes extends HTMLAttributes<HTMLTextAreaElement> {
		autocomplete?: FullAutoFill | undefined | null;
		cols?: number | undefined | null;
		dirname?: string | undefined | null;
		disabled?: boolean | undefined | null;
		form?: string | undefined | null;
		maxlength?: number | undefined | null;
		minlength?: number | undefined | null;
		name?: string | undefined | null;
		placeholder?: string | undefined | null;
		readonly?: boolean | undefined | null;
		required?: boolean | undefined | null;
		rows?: number | undefined | null;
		value?: any; // TODO
		wrap?: 'hard' | 'soft' | undefined | null;

		'on:change'?: ChangeEventHandler<HTMLTextAreaElement> | undefined | null;
		onchange?: ChangeEventHandler<HTMLTextAreaElement> | undefined | null;

		'bind:value'?: any;

		enableFileUpload?: boolean;
		enableFileLink?: boolean;
		enabledExtensions?: string[];
	}

	let {
		ref = $bindable(null),
		value = $bindable(TipTapDefaultContent),
		class: className,
		disabled = false,
		enableFileUpload = true,
		enabledExtensions = ['all'],
		enableFileLink = false,
		...restProps
	}: WithoutChildren<WithElementRef<TipTapAttributes>> = $props();

	let editor = $state() as Readable<Editor>;
	let isUpdatingFromOutside = false;

	async function uploadFile(file: File): Promise<FileUploadResponse> {
		const formData = new FormData();
		formData.append('file', file);

		try {
			const res = await API.request<FormData, FileUploadResponse>('/api/v2/files', {
				method: 'POST',
				body: formData
			});
			return res;
		} catch (err) {
			console.error(err);
			throw new Error(`Upload failed: ${err}`);
		}
	}

	onMount(() => {
		editor = createEditor({
			extensions: [
				StarterKit.configure({
					heading: {
						levels: [1, 2, 3]
					},
					bulletList: {
						HTMLAttributes: {
							class: 'list-disc'
						}
					},
					orderedList: {
						HTMLAttributes: {
							class: 'list-decimal'
						}
					},
					listItem: {
						HTMLAttributes: {
							class: 'ml-4'
						}
					},
					link: {
						openOnClick: false,
						enableClickSelection: true,
						HTMLAttributes: {
							class: 'underline'
						}
					},

					blockquote: false, // TODO
					codeBlock: false, // TODO
					code: false // TODO
				}),
				Subscript,
				Superscript,
				Hightlight,

				TableKit.configure({
					table: {
						resizable: true
					}
				}),

				CustomInlineMath.configure({
					onClick: (node, pos) => {
						const newCalculation = prompt('Enter new calculation:', node.attrs.latex);
						if (newCalculation) {
							editor;
							$editor
								.chain()
								.setNodeSelection(pos)
								.updateInlineMath({ latex: newCalculation })
								.focus()
								.run();
						} else {
							$editor.chain().deleteInlineMath().focus().run();
						}
					}
				}),
				CustomBlockMath.configure({
					onClick: (node, pos) => {
						const newCalculation = prompt('Enter new calculation:', node.attrs.latex);
						if (newCalculation) {
							$editor
								.chain()
								.setNodeSelection(pos)
								.updateBlockMath({ latex: newCalculation })
								.focus()
								.run();
						} else {
							$editor.chain().deleteBlockMath().focus().run();
						}
					}
				}),

				FileHandler.configure({
					onDrop: (currentEditor: any, files: File[], pos: number) => {
						files.forEach(async (file: File) => {
							if (enableFileUpload) {
								if (enabledExtensions.includes('all') || enabledExtensions.includes(file.type)) {
									let res = await uploadFile(file);
									currentEditor
										.chain()
										.insertContentAt(pos, {
											type: 'custom-image',
											attrs: {
												id: res.id,
												mode: 'storage',
												src: res.filename,
												type: file.type,
												originalFilename: res.originalName,
												width: 1000
											}
										})
										.focus()
										.run();
								}
							}
						});
					},
					onPaste: (currentEditor: any, files: File[], htmlContent: string | undefined) => {
						files.forEach(async (file) => {
							if (enableFileUpload) {
								if (enabledExtensions.includes('all') || enabledExtensions.includes(file.type)) {
									let res = await uploadFile(file);
									currentEditor
										.chain()
										.insertContentAt(currentEditor.state.selection.anchor, {
											type: 'custom-image',
											attrs: {
												id: res.id,
												mode: 'storage',
												src: res.filename,
												type: file.type,
												originalFilename: res.originalName,
												width: 1000
											}
										})
										.focus()
										.run();
								}
							}
						});
					}
				}),

				Details.configure({
					openClassName: 'is-open',
					HTMLAttributes: {
						class: 'details'
					}
				}),
				DetailsContent,
				DetailsSummary,

				Typography,
				MathReplacer,
				CharsReplacer,
				TextAlign.configure({
					types: ['heading', 'paragraph']
				}),
				// CodeBlockLowlight.configure({
				// 	lowlight: createLowlight(LLCommon)
				// }),
				TipTapImage
			],
			content: value,
			onTransaction({ editor }) {
				if (!isUpdatingFromOutside) {
					value = editor.getJSON();
				}
			},
			editable: !disabled
		});
	});

	onDestroy(() => {
		if (editor) {
			$editor.destroy();
		}
	});

	// Shadcn
	import { cn, type WithElementRef, type WithoutChildren } from '$lib/utils.js';
	import type { ChangeEventHandler, FullAutoFill, HTMLAttributes } from 'svelte/elements';
	import type { Readable } from 'svelte/store';
	import Menu from './menu/menu.svelte';
	import { TipTapDefaultContent } from '$lib/constants';
	import { API } from '$lib/services/api.svelte';
	import type { FileUploadResponse } from '$lib/api_types';
</script>

<div>
	{#if editor}
		<Menu bind:editor disabled={disabled ?? false} {enableFileUpload} {enableFileLink}></Menu>
	{/if}
	<div
		class={cn(
			'border-input placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:bg-input/30 field-sizing-content w-full rounded-md rounded-t-none border bg-transparent text-base shadow-xs transition-[color,box-shadow] outline-none focus-visible:ring-[3px] disabled:cursor-not-allowed disabled:opacity-50 md:text-sm',
			className
		)}
	>
		<EditorContent editor={$editor} class="overflow-auto" />
	</div>
</div>
