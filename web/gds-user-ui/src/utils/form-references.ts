// default value should be initialized from localstorage and assigned to each property
export const getCertificateRegistrationDefaultValue = () => {
  return {
    entity: {
      country_of_registration: 'united states',
      name: {
        name_identifiers: [
          {
            legal_person_name: '',
            legal_person_name_identifier_type: 'LEGAL_PERSON_NAME_IDENTIFIER_LEGL'
          }
        ],
        local_name_identifiers: [],
        phonetic_name_identifiers: []
      },
      geographic_addresses: [
        { address_type: 2, address_line: ['Address 1', 'Address 2', 'Address 3'], country: 'US' },
        { address_type: 1, address_line: ['Address 1', 'Address 2', 'Address 3'], country: 'CI' }
      ],
      national_identification: {
        national_identifier: 'name identifier',
        national_identifier_type: 9,
        country_of_issue: 'country of issue',
        registration_authority: 'registration'
      }
    },
    contacts: {
      administrative: {
        name: 'John Doe',
        email: 'jdoe@example.com',
        phone: '+13565645646'
      },
      technical: {
        name: 'Jane Doe',
        email: 'jane@example.com',
        phone: '+13565645646'
      },
      billing: {},
      legal: {}
    },
    trisa_endpoint_testnet: {
      trisa_endpoint: 'TRISA Endpoint TestNet',
      common_name: 'Common name'
    },
    trisa_endpoint_mainnet: {
      trisa_endpoint: 'TRISA Endpoint MainNet',
      common_name: 'Common name'
    },
    website: '',
    business_category: 'PRIVATE_ORGANIZATION',
    vasp_categories: ['ATM'],
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
      financial_transfers_permitted: 'no',
      has_required_regulatory_program: '',
      conducts_customer_kyc: true,
      kyc_threshold: 1700,
      kyc_threshold_currency: 'USD',
      must_comply_travel_rule: true,
      applicable_regulations: [
        {
          name: 'FATF Recommendation 16'
        }
      ],
      compliance_threshold: 20000,
      compliance_threshold_currency: 'XOF',
      must_safeguard_pii: true,
      safeguards_pii: true
    }
  };
};
