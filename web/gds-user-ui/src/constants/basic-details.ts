export const vaspCategories = [
  {
    value: 'Exchange',
    label: 'Centralized Exchange'
  },
  {
    value: 'DEX',
    label: 'Decentralized Exchange'
  },
  {
    value: 'P2P',
    label: 'Person-to-Person Exchange'
  },
  {
    value: 'Kiosk',
    label: 'Kiosk / Crypto ATM Operator'
  },
  {
    label: 'Custodian',
    value: 'Custody Provider'
  },
  {
    label: 'ODC',
    value: 'Over-The-Counter Trading Desk'
  },
  {
    value: 'Fund',
    label: 'Investment Fund - hedge funds, ETFs, and family offices'
  },
  {
    value: 'Project',
    label: 'Token Project'
  },
  {
    value: 'Gambling',
    label: 'Gambling or Gaming Site'
  },
  {
    value: 'Miner',
    label: 'Mining Pool'
  },
  {
    value: 'Mixer',
    label: 'Mixing Service'
  },
  {
    value: 'Individual',
    label: 'Legal Person'
  },
  {
    value: 'Other',
    label: 'Other'
  }
];

export const BUSINESS_CATEGORY = {
  UNKNOWN_ENTITY: 'Unknown Entity',
  PRIVATE_ORGANIZATION: 'Private Organization',
  GOVERNMENT_ENTITY: 'Government Entity',
  BUSINESS_ENTITY: 'Business Entity',
  NON_COMMERCIAL_ENTITY: 'Non Commercial Entity'
};

export const getBusinessCategoryOptions = () => {
  return Object.entries(BUSINESS_CATEGORY).map(([k, v]) => ({ value: k, label: v }));
};
