import { setupI18n } from '@lingui/core';
import { t } from '@lingui/macro';

import * as yup from 'yup';

const _i18n = setupI18n();

export const contactsValidationSchema = yup.object().shape({
  contacts: yup.object().shape({
    administrative: yup.object().shape({
      name: yup.string(),
      email: yup.string().email(_i18n._(t`Email is not valid.`)),
      phone: yup.string()
    }),
    technical: yup
      .object()
      .shape({
        name: yup.string(), // .required(),
        email: yup.string().email(_i18n._(t`Email is not valid.`)),
        // .required(_i18n._(t`Email is required.`)),
        phone: yup.string()
      })
      .required(),
    billing: yup.object().shape({
      name: yup.string(),
      email: yup.string().email(_i18n._(t`Email is not valid.`)),
      phone: yup.string()
    }),
    legal: yup.object().shape({
      name: yup.string(), // .required(),
      email: yup.string().email('Email is not valid.'), // .required('Email is required.'),
      phone: yup.string()
      // .required(
      //   'A business phone number is required to complete physical verification for MainNet registration. Please provide a phone number where the Legal/ Compliance contact can be contacted.'
      // )
    })
    // .required()
  })
});
