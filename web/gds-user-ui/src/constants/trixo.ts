import { t } from '@lingui/macro';
import currencies from 'utils/currencies.json';

const FINANCIAL_TRANSFERTS_PERMITTED_OPTIONS = {
  yes: t`Yes`,
  partially: t`Partially`,
  no: t`No`
};

export const getFinancialTransfertsPermittedOptions = () =>
  Object.entries(FINANCIAL_TRANSFERTS_PERMITTED_OPTIONS).map(([k, v]) => ({ value: k, label: v }));

export const getCurrenciesOptions = () => {
  return Object.entries(currencies).map(([k]) => ({ value: k, label: k }));
};
