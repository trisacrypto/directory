import { setupI18n } from '@lingui/core';
import { t } from '@lingui/macro';
import { format2ShortDate } from 'utils/utils';

import * as yup from 'yup';

const _i18n = setupI18n();
const minDate = new Date('01/01/1800');
const maxDate = new Date();

export const basicDetailsValidationSchema = yup.object().shape({
  website: yup
    .string()
    .url()
    .trim()
    .required(_i18n._(t`Website is a required field`)),
  established_on: yup
    .date()
    .min(
      minDate,
      t`Date of incorporation / establishment must be later than ${new Intl.DateTimeFormat([
        'ban',
        'id'
      ]).format(minDate)}`
    )
    .max(
      new Date(),
      t`Date of incorporation / establishment must be at earlier than ${new Intl.DateTimeFormat([
        'ban',
        'id'
      ]).format(maxDate)}`
    )
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
