import { setupI18n } from '@lingui/core';
import { t } from '@lingui/macro';
import dayjs from 'dayjs';

import * as yup from 'yup';

const DATE_FORMAT = 'DD/MM/YYYY';

const _i18n = setupI18n();
const minDate = '1970-01-01'; // it fix this issue https://github.com/jquense/yup/issues/325
const fromMinDate = '1800-01-01';

export const basicDetailsValidationSchema = yup.object().shape({
  website: yup.string().url().trim(),

  established_on: yup
    .date()
    .transform((value, originalValue, schema) => {
      if (schema.isType(value)) {
        return value;
      }

      return dayjs(originalValue).format(DATE_FORMAT);
    })
    .min(
      fromMinDate,
      t`Date of incorporation / establishment must be later than` +
        ` ${dayjs(minDate).format(DATE_FORMAT)}.` +
        `Please select a date no earlier than 1800-01-01.`
    )
    .min(
      minDate,
      t`Date of incorporation / establishment must be later than` +
        ` ${dayjs(minDate).format(DATE_FORMAT)}.`
    )
    .max(new Date(), t`Date of incorporation / establishment must be earlier than current date.`)
    .nullable()
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
    }),

  organization_name: yup.string().trim(),
  // .required(_i18n._(t`Organization name is required.`)),
  business_category: yup.string().nullable(true),
  vasp_categories: yup.array().of(yup.string()).nullable(true)
});
