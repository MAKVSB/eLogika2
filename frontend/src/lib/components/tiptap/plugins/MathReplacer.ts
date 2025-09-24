import { Extension, textInputRule } from '@tiptap/core';

const rules = [
	{
		from: '@forall',
		to: '∀',
		description: 'Universal quantifier'
	},
	{
		from: '@exists',
		to: '∃',
		description: 'Existential quantifier'
	},
	{
		from: '@notexists',
		to: '∄',
		description: 'There does not exist'
	},
	{
		from: '@not',
		to: '¬',
		description: 'Logical NOT'
	},
	{
		from: '@and',
		to: '∧',
		description: 'Logical AND'
	},
	{
		from: '@or',
		to: '∨',
		description: 'Logical OR'
	},
	{
		from: '@xor',
		to: '⊻',
		description: 'Exclusive OR'
	},
	{
		from: '@implies',
		to: '→',
		description: 'Implies'
	},
	{
		from: '@iff',
		to: '↔',
		description: 'If and only if'
	},
	{
		from: '@true',
		to: '⊤',
		description: 'Logical truth'
	},
	{
		from: '@false',
		to: '⊥',
		description: 'Logical falsehood'
	},
	{
		from: '@entails',
		to: '⊢',
		description: 'Syntactic entailment'
	},
	{
		from: '@models',
		to: '⊨',
		description: 'Semantic entailment'
	},
	{
		from: '@forces',
		to: '⊩',
		description: 'Forcing (set theory)'
	},
	{
		from: '@in',
		to: '∈',
		description: 'Element of'
	},
	{
		from: '@notin',
		to: '∉',
		description: 'Not an element of'
	},
	{
		from: '@ni',
		to: '∋',
		description: 'Contains as member'
	},
	{
		from: '@notni',
		to: '∌',
		description: 'Does not contain as member'
	},
	{
		from: '@empty',
		to: '∅',
		description: 'Empty set'
	},
	{
		from: '@intersect',
		to: '∩',
		description: 'Set intersection'
	},
	{
		from: '@union',
		to: '∪',
		description: 'Set union'
	},
	{
		from: '@subset',
		to: '⊂',
		description: 'Proper subset'
	},
	{
		from: '@subseteq',
		to: '⊆',
		description: 'Subset or equal'
	},
	{
		from: '@notsubset',
		to: '⊄',
		description: 'Not a subset'
	},
	{
		from: '@supset',
		to: '⊃',
		description: 'Superset'
	},
	{
		from: '@supseteq',
		to: '⊇',
		description: 'Superset or equal'
	},
	{
		from: '@plus',
		to: '+',
		description: 'Addition'
	},
	{
		from: '@minus',
		to: '−',
		description: 'Subtraction'
	},
	{
		from: '@times',
		to: '×',
		description: 'Multiplication'
	},
	{
		from: '@divide',
		to: '÷',
		description: 'Division'
	},
	{
		from: '@asterisk',
		to: '∗',
		description: 'Asterisk operator'
	},
	{
		from: '@dot',
		to: '⋅',
		description: 'Dot operator'
	},
	{
		from: '@slash',
		to: '∕',
		description: 'Division slash'
	},
	{
		from: '@sum',
		to: '∑',
		description: 'Summation'
	},
	{
		from: '@product',
		to: '∏',
		description: 'Product'
	},
	{
		from: '@sqrt',
		to: '√',
		description: 'Square root'
	},
	{
		from: '@cuberoot',
		to: '∛',
		description: 'Cube root'
	},
	{
		from: '@fourthroot',
		to: '∜',
		description: 'Fourth root'
	},
	{
		from: '@infinity',
		to: '∞',
		description: 'Infinity'
	},
	{
		from: '@proportional',
		to: '∝',
		description: 'Proportional to'
	},
	{
		from: '@setminus',
		to: '∖',
		description: 'Set difference'
	},
	{
		from: '@neq',
		to: '≠',
		description: 'Not equal'
	},
	{
		from: '@eq',
		to: '=',
		description: 'Equal'
	},
	{
		from: '@sim',
		to: '∼',
		description: 'Similar'
	},
	{
		from: '@equiv',
		to: '≡',
		description: 'Equivalent'
	},
	{
		from: '@congruent',
		to: '≅',
		description: 'Congruent'
	},
	{
		from: '@xor',
		to: '⊕',
		description: 'Exclusive OR'
	},
	{
		from: '@oplux',
		to: '⊕',
		description: 'Direct sum'
	},
	{
		from: '@partial',
		to: '∂',
		description: 'Partial derivative'
	},
	{
		from: '@nabla',
		to: '∇',
		description: 'Nabla operator (gradient)'
	},
	{
		from: '@integral',
		to: '∫',
		description: 'Integral'
	},
	{
		from: '@doubleintegral',
		to: '∬',
		description: 'Double integral'
	},
	{
		from: '@tripleintegral',
		to: '∭',
		description: 'Triple integral'
	},
	{
		from: '@contourintegral',
		to: '∮',
		description: 'Contour integral'
	},
	{
		from: '@surfaceintegral',
		to: '∯',
		description: 'Surface integral'
	},
	{
		from: '@volumeintegral',
		to: '∰',
		description: 'Volume integral'
	},
	{
		from: '@equal',
		to: '=',
		description: 'Equal to'
	},
	{
		from: '@approx',
		to: '≈',
		description: 'Approximately equal to'
	},
	{
		from: '@leq',
		to: '≤',
		description: 'Less than or equal to'
	},
	{
		from: '@geq',
		to: '≥',
		description: 'Greater than or equal to'
	},
	{
		from: '@lt',
		to: '<',
		description: 'Less than'
	},
	{
		from: '@gt',
		to: '>',
		description: 'Greater than'
	},
	{
		from: '@angle',
		to: '∠',
		description: 'Angle'
	},
	{
		from: '@therefore',
		to: '∴',
		description: 'Therefore'
	},
	{
		from: '@because',
		to: '∵',
		description: 'Because'
	},
	{
		from: '@parallel',
		to: '∥',
		description: 'Parallel to'
	},
	{
		from: '@otimes',
		to: '⊗',
		description: 'Tensor product'
	},
	{
		from: '@boxplus',
		to: '⊞',
		description: 'Box plus'
	},
	{
		from: '@boxminus',
		to: '⊟',
		description: 'Box minus'
	},
	{
		from: '@boxtimes',
		to: '⊠',
		description: 'Box times'
	},
	{
		from: '@boxdot',
		to: '⊡',
		description: 'Box dot'
	},
	{
		from: '@circledcirc',
		to: '⊚',
		description: 'Circled circle'
	},
	{
		from: '@circledast',
		to: '⊛',
		description: 'Circled asterisk'
	},
	{
		from: '@circleeq',
		to: '⊜',
		description: 'Circle equals'
	},
	{
		from: '@circleddash',
		to: '⊝',
		description: 'Circled dash'
	}
];

function escapeRegExp(str: string) {
	return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&') + ' ';
}

export const MathReplacer = Extension.create({
	name: 'mathReplacer',

	addInputRules() {
		let inputrules = [];

		for (const rule of rules) {
			new RegExp(`${rule.from}\\s?$`);
			inputrules.push(
				textInputRule({
					find: new RegExp(`${escapeRegExp(rule.from)}\\s?$`),
					replace: rule.to
				})
			);
		}

		return inputrules;
	}
});

export default MathReplacer;
