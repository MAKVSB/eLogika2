import TextInput from './form-input-text.svelte';
import YearInput from './form-input-year.svelte';
import SingleSelect, {
	type SelectOption,
	type SelectOptions
} from './form-input-select-single.svelte';
import MultiSelect from './form-input-select-multi.svelte';
import Checkbox from './form-input-checkbox.svelte';
import Tiptap from './form-tiptap.svelte';
import Form from './form.svelte';
import Button from './form-button.svelte';
import Wrapper from './form-input-wrapper.svelte';
import TextArea from './form-input-textarea.svelte';
import type z from 'zod/v4';
import type { ErrorObject } from './types';

export {
	Form as Root,
	TextInput,
	YearInput,
	TextArea,
	SingleSelect,
	MultiSelect,
	Checkbox,
	Button,
	Tiptap,
	Wrapper,
	type SelectOptions,
	type SelectOption,
	type ErrorObject
	// Field as FormField,
	// Control as FormControl,
	// Description as FormDescription,
	// Label as FormLabel,
	// FieldErrors as FormFieldErrors,
	// Fieldset as FormFieldset,
	// Legend as FormLegend,
	// ElementField as FormElementField,
	// Button as FormButton
};

export function createForm<U, T extends z.ZodType<any, any, any>>(
	schema: T,
	defaultFormData: z.infer<T> & U
) {
	return {
		schema: schema,
		fields: defaultFormData,
		errors: {} as ErrorObject,
		isSubmitting: false
	};
}
