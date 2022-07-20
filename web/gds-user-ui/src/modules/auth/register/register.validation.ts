import { t } from '@lingui/macro';
import * as yup from 'yup';
const passwordRegex =
  /^(?=.*[A-Z].*[A-Z])(?=.*[!@#$&*])(?=.*[0-9].*[0-9])(?=.*[a-z].*[a-z].*[a-z]).{8}$/;
// ^                         Start anchor
// (?=.*[A-Z].*[A-Z])        Ensure string has two uppercase letters.
// (?=.*[!@#$&*])            Ensure string has one special case letter.
// (?=.*[0-9].*[0-9])        Ensure string has two digits.
// (?=.*[a-z].*[a-z].*[a-z]) Ensure string has three lowercase letters.
// .{8}                      Ensure string is of length 8.
// $                         End anchor.
export const validationSchema = yup.object().shape({
  username: yup
    .string()
    .email(t`Email is not valid`)
    .required(t`Email is required`),
  password: yup
    .string()
    .matches(
      passwordRegex,
      t`
  *At least 8 characters in length 
  * Contain at least 3 of the following 4
  types of characters: 
  * lower case letters (a-z) 
  * upper case letters (A-Z) 
  * numbers (i.e. 0-9) 
  * special characters (e.g. !@#$%^&*)
  `
    )
    .required('Password is required')
});
