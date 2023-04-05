const requiredFieldsByKey = {
  1: ['Organization Name', 'Website', 'Date of Incorporation / Establishment'],
  2: [
    'Name Identifiers',
    'Legal Name Identifiers',
    'Georaphic Addresses ( Country , City , Address, postcode )',
    'National Identifier',
    'Country of Registration'
  ],
  3: [
    'Technical Contact Information',
    'Administrative Contact Information',
    'Legal Contact Information'
  ],
  4: ['Mainnet Common Name', 'Mainnet Endpoint', 'Testnet Common Name', 'Testnet Endpoint']
} as any;

export const getRequiredFieldsByStepKey = (step: number) => {
  return requiredFieldsByKey[step] || [];
};
