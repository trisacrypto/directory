import { setupI18n } from '@lingui/core';
import { t } from '@lingui/macro';

import * as yup from 'yup';

const _i18n = setupI18n();

export const basicDetailsValidationSchema = yup.object().shape({
  website: yup
    .string()
    .url()
    .trim()
    .required(_i18n._(t`Website is a required field`)),
  established_on: yup
    .date()
    .nullable()
    .transform((curr, orig) => (orig === '' ? null : curr))
    .required(_i18n._(t`Invalid date`))
    .test('is-invalidate-date', _i18n._(t`Invalid date / year must be 4 digit`), (value) => {
      if (value) {
        const getYear = value.getFullYear();
        if (getYear.toString().length !== 4) {
          return false;
        } else {
          return true;
        }
      }
      return false;
    })
    .required(),
  organization_name: yup
    .string()
    .trim()
    .required(_i18n._(t`Organization name is required`)),
  business_category: yup.string().nullable(true),
  vasp_categories: yup.array().of(yup.string()).nullable(true)
});
