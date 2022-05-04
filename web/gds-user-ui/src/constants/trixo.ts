import currencies from 'utils/currencies.json';

const FINANCIAL_TRANSFERTS_PERMITTED_OPTIONS = {
  yes: 'Yes',
  partial: 'Partially',
  no: 'No'
};

export const getFinancialTransfertsPermittedOptions = () =>
  Object.entries(FINANCIAL_TRANSFERTS_PERMITTED_OPTIONS).map(([k, v]) => ({ value: k, label: v }));

export const getCurrenciesOptions = () => {
  return Object.entries(currencies).map(([k, v]) => ({ value: k, label: k }));
};
