import * as yup from 'yup';

const trisaEndpointPattern = /^([a-zA-Z0-9.-]+):((?!(0))[0-9]+)$/;

export const validationSchema = [
  yup.object().shape({
    website: yup.string().trim().url().required(),
    established_on: yup.date().nullable(true),
    organization_name: yup.string().trim().required(),
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
          address_line: yup.array().of(yup.string().required()),
          country: yup.string().required()
        })
      ),
      national_identification: yup.object().shape({
        national_identifier: yup.string(),
        national_identifier_type: yup.string(),
        country_of_issue: yup.string(),
        registration_authority: yup.string()
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
      common_name: yup.string()
    }),
    trisa_endpoint_mainnet: yup.object().shape({
      endpoint: yup.string().matches(trisaEndpointPattern, 'trisa endpoint is not valid'),
      common_name: yup.string()
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
      applicable_regulations: yup.array().of(
        yup.object().shape({
          name: yup.string()
        })
      ),
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
