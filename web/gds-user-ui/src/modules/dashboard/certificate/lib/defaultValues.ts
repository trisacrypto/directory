// default value should be initialized from localstorage and assigned to each property
export const getRegistrationDefaultValue = () => {
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
          address_type: '',
          address_line: ['', '', ''],
          country: ''
        }
      ],
      national_identification: {
        national_identifier: '',
        national_identifier_type: null,
        country_of_issue: '',
        registration_authority: ''
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
    trisa_endpoint_testnet: {
      trisa_endpoint: '',
      common_name: ''
    },
    trisa_endpoint_mainnet: {
      trisa_endpoint: '',
      common_name: ''
    },
    website: '',
    business_category: '',
    vasp_categories: [],
    established_on: '',
    organization_name: '',
    trixo: {
      primary_national_jurisdiction: '',
      primary_regulator: '',
      other_jurisdictions: [
        {
          country: '',
          regulator_name: ''
        }
      ],
      financial_transfers_permitted: '',
      has_required_regulatory_program: '',
      conducts_customer_kyc: null,
      kyc_threshold: 0,
      kyc_threshold_currency: 'USD',
      must_comply_travel_rule: null,
      applicable_regulations: [],
      compliance_threshold: 0,
      compliance_threshold_currency: '',
      must_safeguard_pii: null,
      safeguards_pii: null
    }
  };
};
