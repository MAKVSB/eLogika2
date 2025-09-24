export type ErrorObject = {
	[key: string]: string | ErrorObject;
};
export type ValidatedForm = {
	schema?: any;
	fields: any;
	errors: ErrorObject;
	isSubmitting: boolean;
};
