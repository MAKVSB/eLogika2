import { InputRule, PasteRule } from '@tiptap/core';
import { InlineMath } from '@tiptap/extension-mathematics';

export const CustomInlineMath = InlineMath.extend({
	addInputRules() {
		return [
			new InputRule({
				find: /(?<!\$)\$\$([^$\n]+)\$\$(?!\$)$/,
				handler: ({ state, range, match }) => {
					const [, latex] = match;
					const { tr } = state;
					const start = range.from;
					const end = range.to;

					tr.replaceWith(start, end, this.type.create({ latex }));
				}
			}),
			new InputRule({
				find: /(?<!\$)\$([^$\n]+)\$(?!\$)$/,
				handler: ({ state, range, match }) => {
					const [, latex] = match;
					const { tr } = state;
					const start = range.from;
					const end = range.to;

					tr.replaceWith(start, end, this.type.create({ latex }));
				}
			})
		];
	},

	addPasteRules() {
		return [
			new PasteRule({
				find: /(?<!\$)\$\$([^$\n]+)\$\$(?!\$)/g,
				handler: ({ state, range, match }) => {
					const [, latex] = match;
					const { tr } = state;
					const start = range.from;
					const end = range.to;

					tr.replaceWith(start, end, this.type.create({ latex }));
				}
			}),
			new PasteRule({
				find: /(?<!\$)\$([^$\n]+)\$(?!\$)/g,
				handler: ({ state, range, match }) => {
					const [, latex] = match;
					const { tr } = state;
					const start = range.from;
					const end = range.to;

					tr.replaceWith(start, end, this.type.create({ latex }));
				}
			})
		];
	}
});

export default CustomInlineMath;
