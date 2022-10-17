import { setupI18n } from '@lingui/core';
import { t } from '@lingui/macro';
import dayjs from 'dayjs';

import * as yup from 'yup';

const DATE_FORMAT = 'DD/MM/YYYY';

const _i18n = setupI18n();
const minDate = '1970-01-01'; // it fix this issue https://github.com/jquense/yup/issues/325

export const basicDetailsValidationSchema = yup.object().shape({
  website: yup
    .string()
    .url()
    .trim()
    .required(_i18n._(t`Website is a required field`)),
  established_on: yup
    .date()
    .transform((value, originalValue, schema) => {
      if (schema.isType(value)) {
        return value;
      }

      return dayjs(originalValue).format(DATE_FORMAT);
    })
    .min(
      minDate,
      t`Date of incorporation / establishment must be later than` +
        ` ${dayjs(minDate).format(DATE_FORMAT)}` +
        '.'
    )
    .max(new Date(), t`Date of incorporation / establishment must be earlier than current date.`)
    .nullable()
    .typeError(t`Invalid date.`)
    .required(),
  organization_name: yup
    .string()
    .trim()
    .required(_i18n._(t`Organization name is required.`)),
  business_category: yup.string().nullable(true),
  vasp_categories: yup.array().of(yup.string()).nullable(true)
});
