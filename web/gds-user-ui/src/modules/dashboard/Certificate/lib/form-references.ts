// default value should be initialized from localstorage and assigned to each property
export const getCertificateRegistrationDefaultValue = () => {
  return {
    entity: {
      country_of_registration: '',
      name: {
        name_identifiers: [
          {
            legal_person_name: '',
            legal_person_name_identifier_type: ''
          }
        ],
        local_name_identifiers: [],
        phonetic_name_identifiers: []
      },
      geographic_addresses: [],
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
      conducts_customer_kyc: false,
      kyc_threshold: 0,
      kyc_threshold_currency: 'USD',
      must_comply_travel_rule: false,
      applicable_regulations: [],
      compliance_threshold: 0,
      compliance_threshold_currency: '',
      must_safeguard_pii: false,
      safeguards_pii: false
    }
  };
};
