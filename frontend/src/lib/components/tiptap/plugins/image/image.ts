import { mergeAttributes, Node, nodeInputRule } from '@tiptap/core';
import { SvelteNodeViewRenderer } from 'svelte-tiptap';
import ImageComponent from './image-component.svelte';

export interface ImageOptions {
	/**
	 * Controls if base64 images are allowed. Enable this if you want to allow
	 * base64 image urls in the `src` attribute.
	 * @default false
	 * @example true
	 */
	allowBase64: boolean;

	/**
	 * HTML attributes to add to the image element.
	 * @default {}
	 * @example { class: 'foo' }
	 */
	HTMLAttributes: Record<string, any>;
}

export interface SetImageOptions {
	id?: number;
	mode: string;
	type: string;
	src: string;
	originalFilename: string;
	width?: number;
}

declare module '@tiptap/core' {
	interface Commands<ReturnType> {
		image: {
			/**
			 * Add an image
			 * @param options The image attributes
			 * @example
			 * editor
			 *   .commands
			 *   .setImage({ src: 'https://tiptap.dev/logo.png', alt: 'tiptap', title: 'tiptap logo' })
			 */
			setImage: (options: SetImageOptions) => ReturnType;
		};
	}
}

/**
 * Matches an image to a ![image](src "title") on input.
 */
export const inputRegex = /(?:^|\s)(!\[(.+|:?)]\((\S+)(?:(?:\s+)["'](\S+)["'])?\))$/;

/**
 * This extension allows you to insert images.
 * @see https://www.tiptap.dev/api/nodes/image
 */
export const Image = Node.create<ImageOptions>({
	name: 'custom-image',

	addOptions() {
		return {
			allowBase64: false,
			HTMLAttributes: {}
		};
	},

	group() {
		return 'block';
	},

	draggable: true,

	selectable: true,

	addAttributes() {
		return {
			id: {
				default: null
			},
			mode: {
				default: 'url'
			},
			type: {
				default: null
			},
			src: {
				default: null
			},
			originalFilename: {
				default: null
			},
			width: {
				default: null
			}
		};
	},

	parseHTML() {
		return [
			{
				tag: this.options.allowBase64 ? 'img[src]' : 'img[src]:not([src^="data:"])'
			},
			{
				tag: this.options.allowBase64
					? 'svelte-image-component[src]'
					: 'svelte-image-component[src]:not([src^="data:"])'
			}
		];
	},

	renderHTML({ HTMLAttributes }) {
		return ['svelte-image-component', mergeAttributes(this.options.HTMLAttributes, HTMLAttributes)];
	},

	addCommands() {
		return {
			setImage:
				(options) =>
				({ commands }) => {
					return commands.insertContent({
						type: this.name,
						attrs: options
					});
				}
		};
	},

	addInputRules() {
		return [
			nodeInputRule({
				find: inputRegex,
				type: this.type,
				getAttributes: (match) => {
					const [, , alt, src, title] = match;

					return { src, alt, title };
				}
			})
		];
	},

	addNodeView() {
		return SvelteNodeViewRenderer(ImageComponent);
	}
});
