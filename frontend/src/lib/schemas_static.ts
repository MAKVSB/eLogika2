import z from 'zod/v4';
import { en, cs } from 'zod/v4/locales';
import { getLocale } from '$lib/paraglide/runtime';
let locale = getLocale();
if (locale === 'cs') {
	z.config(cs());
} else {
	z.config(en());
}

export const UserChangePassRequestSchemaStatic = z
	.object({
		oldPassword: z.string(),
		newPassword: z.string(),
		newPasswordRep: z.string()
	})
	.refine(({ newPassword, newPasswordRep }) => {
		if (newPassword !== newPasswordRep) {
			return [
				{
					code: 'custom',
					message: 'Passwords must match',
					path: ['newPassword']
				},
				{
					code: 'custom',
					message: 'Passwords must match',
					path: ['newPasswordRep']
				}
			];
		}
		return true;
	});
export type UserChangePassRequest = z.infer<typeof UserChangePassRequestSchemaStatic>;
