import { BsCartXFill } from 'react-icons/bs';
import * as yup from 'yup';
import _ from 'lodash';

const trisaEndpointPattern = /^([a-zA-Z0-9.-]+):((?!(0))[0-9]+)$/;
const commonNameRegex = /^[A-Za-z0-9\s]+\.[A-Za-z0-9\s]+$/;

export const validationSchema = [
  yup.object().shape({
    website: yup.string().url().trim().required(),
    established_on: yup
      .date()
      .nullable()
      .transform((curr, orig) => (orig === '' ? null : curr))
      .required('Invalid date')
      .test('is-invalidate-date', 'Invalid date / year must be 4 digit ', (value) => {
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
      .notRequired(),
    organization_name: yup.string().trim().required('Organization name is required'),
    business_category: yup.string().nullable(true),
    vasp_categories: yup.array().of(yup.string()).nullable(true)
  }),
  yup.object().shape({
    entity: yup.object().shape({
      name: yup.object().shape({
        name_identifiers: yup.array().of(
          yup.object().shape({
            legal_person_name: yup.string(),
            legal_person_name_identifier_type: yup.string()
          })
        ),
        local_name_identifiers: yup.array().of(
          yup.object().shape({
            legal_person_name: yup.string(),
            legal_person_name_identifier_type: yup.string()
          })
        ),
        phonetic_name_identifiers: yup.array().of(
          yup.object().shape({
            legal_person_name: yup.string(),
            legal_person_name_identifier_type: yup.string()
          })
        )
      }),
      geographic_addresses: yup.array().of(
        yup.object().shape({
          address_type: yup.string().required(),
          address_line: yup.array(),
          'address_line[0]': yup
            .string()
            .test('test-0', 'addresse line 0', (value: any, ctx: any): any => {
              return ctx && ctx.parent && ctx.parent.address_line[0];
            }),
          'address_line[2]': yup
            .string()
            .test('test-0', 'addresse line 0', (value: any, ctx: any): any => {
              return ctx && ctx.parent && ctx.parent.address_line[2];
            }),
          country: yup.string().required()
        })
      ),
      national_identification: yup.object().shape({
        national_identifier: yup.string(),
        national_identifier_type: yup.string().required('National identification type is required'),
        country_of_issue: yup.string(),
        registration_authority: yup
          .string()
          .test(
            'registrationAuthority',
            'Registration Authority cannot be left empty',
            (value, ctx) => {
              console.log('ctex', ctx.parent.national_identifier_type);
              console.log('ctex value', typeof value);
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
  }),
  yup.object().shape({
    contacts: yup.object().shape({
      administrative: yup.object().shape({
        name: yup.string(),
        email: yup.string().email(),
        phone: yup.string()
      }),
      technical: yup
        .object()
        .shape({
          name: yup.string().required(),
          email: yup.string().email().required(),
          phone: yup.string()
        })
        .required(),
      billing: yup.object().shape({
        name: yup.string(),
        email: yup.string().email(),
        phone: yup.string()
      }),
      legal: yup
        .object()
        .shape({
          name: yup.string().required(),
          email: yup.string().email().required(),
          phone: yup.string()
        })
        .required()
    })
  }),
  yup.object().shape({
    trisa_endpoint: yup.string().trim(),
    trisa_endpoint_testnet: yup.object().shape({
      endpoint: yup.string().matches(trisaEndpointPattern, 'trisa endpoint is not valid'),
      common_name: yup
        .string()
        .matches(
          commonNameRegex,
          'Common name should not contain special characters, no spaces and must have a dot(.) in it'
        )
    }),
    trisa_endpoint_mainnet: yup.object().shape({
      endpoint: yup
        .string()
        .test(
          'uniqueMainetEndpoint',
          'TestNet and MainNet endpoints should not be the same',
          (value, ctx: any): any => {
            return ctx.from[1].value.trisa_endpoint_testnet.endpoint !== value;
          }
        )
        .matches(trisaEndpointPattern, 'trisa endpoint is not valid'),
      common_name: yup
        .string()
        .matches(
          commonNameRegex,
          'Common name should not contain special characters, no spaces and must have a dot(.) in it'
        )
    })
  }),
  yup.object().shape({
    trixo: yup.object().shape({
      primary_national_jurisdiction: yup.string(),
      primary_regulator: yup.string(),
      other_jurisdictions: yup.array().of(
        yup.object().shape({
          country: yup.string(),
          regulator_name: yup.string()
        })
      ),
      financial_transfers_permitted: yup.string(),
      has_required_regulatory_program: yup.string(),
      conducts_customer_kyc: yup.boolean(),
      kyc_threshold: yup.number(),
      kyc_threshold_currency: yup.string(),
      must_comply_travel_rule: yup.boolean(),
      applicable_regulations: yup
        .array()
        .of(
          yup.object().shape({
            name: yup.string()
          })
        )
        .transform((value, originalValue) => {
          if (originalValue) {
            return originalValue.filter((item: any) => item.name.length > 0);
          }
          return value;

          // remove empty items
        }),
      compliance_threshold: yup.number(),
      compliance_threshold_currency: yup.string(),
      must_safeguard_pii: yup.boolean(),
      safeguards_pii: yup.boolean()
    })
  })
];

// export const certificateRegistrationValidationSchema = yup.object().shape({
//   entity: yup.object().shape({
//     country_of_registration: yup.string(),
//     name: yup.object().shape({
//       name_identifiers: yup.array().of(
//         yup.object().shape({
//           legal_person_name: yup.string(),
//           legal_person_name_identifier_type: yup.string()
//         })
//       ),
//       local_name_identifiers: yup.array().of(
//         yup.object().shape({
//           legal_person_name: yup.string(),
//           legal_person_name_identifier_type: yup.string()
//         })
//       ),
//       phonetic_name_identifiers: yup.array().of(
//         yup.object().shape({
//           legal_person_name: yup.string(),
//           legal_person_name_identifier_type: yup.string()
//         })
//       )
//     }),
//     geographic_addresses: yup.array().of(
//       yup.object().shape({
//         address_type: yup.number(),
//         address_line: yup.array().of(yup.string())
//       })
//     ),
//     national_identification: yup.object().shape({
//       national_identifier: yup.string(),
//       national_identifier_type: yup.number(),
//       country_of_issue: yup.string(),
//       registration_authority: yup.string()
//     })
//   }),
//   contacts: yup.object().shape({
//     administrative: yup.object().shape({
//       name: yup.string(),
//       email: yup.string().email(),
//       phone: yup.string()
//     }),
//     technical: yup
//       .object()
//       .shape({
//         name: yup.string().required(),
//         email: yup.string().email().required(),
//         phone: yup.string()
//       })
//       .required(),
//     billing: yup.object().shape({
//       name: yup.string(),
//       email: yup.string().email(),
//       phone: yup.string()
//     }),
//     legal: yup
//       .object()
//       .shape({
//         name: yup.string().required(),
//         email: yup.string().email().required(),
//         phone: yup.string()
//       })
//       .required()
//   }),
//   trisa_endpoint: yup.string().trim().matches(trisaEndpointPattern, 'trisa endpoint is not valid'),
//   common_name: yup.string(),
//   website: yup
//     .string()
//     .trim()
//     .url()
//     .test('empty-check', 'Website is required', (value) => value !== '')
//     .required(),
//   business_category: yup.string().nullable(true),
//   vasp_categories: yup.array().of(yup.string()).nullable(true),
//   established_on: yup.date().nullable(true),
//   trixo: yup.object().shape({
//     primary_national_jurisdiction: yup.string(),
//     primary_regulator: yup.string(),
//     other_jurisdictions: yup.array().of(
//       yup.object().shape({
//         country: yup.string(),
//         regulator_name: yup.string()
//       })
//     ),
//     financial_transfers_permitted: yup.string(),
//     has_required_regulatory_program: yup.string(),
//     conducts_customer_kyc: yup.boolean(),
//     kyc_threshold: yup.number(),
//     kyc_threshold_currency: yup.string(),
//     must_comply_travel_rule: yup.boolean(),
//     applicable_regulations: yup.array().of(
//       yup.object().shape({
//         name: yup.string()
//       })
//     ),
//     compliance_threshold: yup.number(),
//     compliance_threshold_currency: yup.string(),
//     must_safeguard_pii: yup.boolean(),
//     safeguards_pii: yup.boolean()
//   })
// });
