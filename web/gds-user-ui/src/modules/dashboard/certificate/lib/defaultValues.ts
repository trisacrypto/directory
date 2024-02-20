import { addressTypeEnum } from "constants/address";

export const getRegistrationDefaultValues = () => {
  return {
    entity: {
      country_of_registration: '',
      name: {
        name_identifiers: [
          {
            legal_person_name: '',
            legal_person_name_identifier_type: 'LEGAL_PERSON_NAME_TYPE_CODE_LEGL'
          }
        ],
        local_name_identifiers: [],
        phonetic_name_identifiers: []
      },
      geographic_addresses: [
        {
          address_type: addressTypeEnum.ADDRESS_TYPE_BIZZ,
          address_line: [],
          country: ''
        }
      ],
      national_identification: {
        national_identifier: undefined,
        national_identifier_type: 'NATIONAL_IDENTIFIER_TYPE_CODE_LEIX',
        country_of_issue: '',
        registration_authority: 'RA777777'
      }
    },
    contacts: {
      administrative: {
        name: '',
        email: '',
        phone: ''
      },
      technical: {
        name: '',
        email: '',
        phone: ''
      },
      billing: {},
      legal: {}
    },
    testnet: {
      endpoint: '',
      common_name: ''
    },
    mainnet: {
      endpoint: '',
      common_name: ''
    },
    website: '',
    business_category: 'UNKNOWN_ENTITY',
    vasp_categories: [],
    established_on: '',
    organization_name: '',
    trixo: {
      primary_national_jurisdiction: '',
      primary_regulator: '',
      other_jurisdictions: [],
      financial_transfers_permitted: 'no',
      has_required_regulatory_program: 'no',
      conducts_customer_kyc: false,
      kyc_threshold: undefined,
      kyc_threshold_currency: 'USD',
      must_comply_travel_rule: false,
      applicable_regulations: ['FATF Recommendation 16'],
      compliance_threshold: undefined,
      compliance_threshold_currency: 'USD',
      must_safeguard_pii: false,
      safeguards_pii: false
    }
  };
};

export type RegistrationFormValues = ReturnType<typeof getRegistrationDefaultValues>;
