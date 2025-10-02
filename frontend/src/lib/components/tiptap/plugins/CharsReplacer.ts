import { Extension, textInputRule, textPasteRule } from '@tiptap/core';

const rules = [
	{
		from: 'ﬂ',
		to: 'fl',
		description: 'Weird overleaf behavior'
	}
];

function escapeRegExp(str: string) {
	return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

export const CharsReplacer = Extension.create({
	name: 'charsReplacer',

	addInputRules() {
		let inputrules = [];

		for (const rule of rules) {
			inputrules.push(
				textInputRule({
					find: new RegExp(`${escapeRegExp(rule.from)}\\s?$`),
					replace: rule.to
				})
			);
		}

		return inputrules;
	},

	addPasteRules() {
		let inputrules = [];

		for (const rule of rules) {
			inputrules.push(
				textPasteRule({
					find: new RegExp(`${escapeRegExp(rule.from)}`, 'g'),
					replace: rule.to
				})
			);
		}

		return inputrules;
	}
});

export default CharsReplacer;
