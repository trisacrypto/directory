import { setupI18n } from '@lingui/core';
import { t } from '@lingui/macro';

import * as yup from 'yup';

const _i18n = setupI18n();

export const legalPersonValidationSchemam = yup.object().shape({
  entity: yup.object().shape({
    country_of_registration: yup.string().required(_i18n._(t`Country of registration is required`)),
    name: yup.object().shape({
      name_identifiers: yup.array(
        yup.object().shape({
          legal_person_name: yup
            .string()
            .test(
              'notEmptyIfIdentifierTypeExist',
              _i18n._(t`Legal name is required`),
              (value, ctx): any => {
                return !(ctx.parent.legal_person_name_identifier_type && !value);
              }
            ),
          legal_person_name_identifier_type: yup.string().when('legal_person_name', {
            is: (value: string) => !!value,
            then: yup.string().required(_i18n._(t`Name Identifier Type is required`))
          })
        })
      ),
      local_name_identifiers: yup.array(
        yup.object().shape({
          legal_person_name: yup
            .string()
            .test(
              'notEmptyIfIdentifierTypeExist',
              _i18n._(t`Legal name is required`),
              (value, ctx): any => {
                return !(ctx.parent.legal_person_name_identifier_type && !value);
              }
            ),
          legal_person_name_identifier_type: yup.string().when('legal_person_name', {
            is: (value: string) => !!value,
            then: yup.string().required(_i18n._(t`Name Identifier Type is required`))
          })
        })
      ),
      phonetic_name_identifiers: yup.array(
        yup.object().shape({
          legal_person_name: yup
            .string()
            .test(
              'notEmptyIfIdentifierTypeExist',
              _i18n._(t`Legal name is required`),
              (value, ctx): any => {
                return !(ctx.parent.legal_person_name_identifier_type && !value);
              }
            ),
          legal_person_name_identifier_type: yup.string().when('legal_person_name', {
            is: (value: string) => !!value,
            then: yup.string().required(_i18n._(t`Name Identifier Type is required`))
          })
        })
      )
    }),
    geographic_addresses: yup.array().of(
      yup.object().shape({
        address_line: yup.array(),
        'address_line[0]': yup
          .string()
          .test('address-line-1', 'addresse line 1', (value: any, ctx: any): any => {
            return ctx && ctx.parent && ctx.parent.address_line[0];
          }),
        country: yup.string().required(),
        town_name: yup.string().required(),
        post_code: yup.string().required(),
        country_sub_division: yup.string().required(),
        address_type: yup.string().required()
      })
    ),
    national_identification: yup.object().shape({
      national_identifier: yup.string().required(),
      national_identifier_type: yup.string(),
      country_of_issue: yup.string(),
      registration_authority: yup
        .string()
        .test(
          'registrationAuthority',
          _i18n._(t`Registration Authority cannot be left empty`),
          (value, ctx) => {
            if (
              ctx.parent.national_identifier_type !== 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX' &&
              !value
            ) {
              return false;
            }

            return true;
          }
        )
    })
  })
});
