import { InputRule, PasteRule } from '@tiptap/core';
import { BlockMath } from '@tiptap/extension-mathematics';

export const CustomBlockMath = BlockMath.extend({
	addInputRules() {
		return [
			new InputRule({
				find: /(?<!\$)\$\$\$([^$\n]+)\$\$\$(?!\$)$/,
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
				find: /(?<!\$)\$\$\$([^$\n]+)\$\$\$(?!\$)/g,
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

export default CustomBlockMath;
