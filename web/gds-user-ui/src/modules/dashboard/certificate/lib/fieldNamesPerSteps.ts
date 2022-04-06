export const fieldNamesPerSteps = {
  basicDetails: ['website', 'established_on', 'business_category', 'vasp_categories'],
  legalPerson: [
    'entity.name.name_identifiers',
    'entity.name.local_name_identifiers',
    'entity.name.phonetic_name_identifiers',
    'entity.geographic_addresses',
    'entity.country_of_registration',
    'entity.national_identification.national_identifier',
    'entity.national_identification.national_identifier_type',
    'entity.national_identification.country_of_issue',
    'entity.national_identification.registration_authority'
  ],
  contacts: [
    ...['administrative', 'technical', 'billing', 'legal'].flatMap((value) => [
      `contacts.${value}.name`,
      `contacts.${value}.email`,
      `contacts.${value}.phone`
    ])
  ],
  trisaImplementation: [
    ...['trisa_endpoint_testnet', 'trisa_endpoint_mainnet'].flatMap((value) => [
      `${value}.common_name`,
      `${value}.endpoint`
    ])
  ],
  trixoImplementation: [
    'trixo.primary_national_jurisdiction',
    'trixo.primary_regulator',
    'trixo.financial_transfers_permitted',
    'trixo.has_required_regulatory_program',
    'trixo.conducts_customer_kyc',
    'trixo.kyc_threshold',
    'trixo.kyc_threshold_currency',
    'trixo.must_comply_travel_rule',
    'trixo.compliance_threshold',
    'trixo.compliance_threshold_currency',
    'trixo.must_safeguard_pii',
    'trixo.safeguards_pii'
  ]
};
