import { t } from '@lingui/macro';

export const NATIONAL_IDENTIFICATION = {
  NATIONAL_IDENTIFIER_TYPE_CODE_LEIX: t`Legal Entity Identifier (LEI)`,
  NATIONAL_IDENTIFIER_TYPE_CODE_TXID: t`Tax Identification Number`,
  NATIONAL_IDENTIFIER_TYPE_CODE_RAID: t`Registration Authority Identifier`,
  NATIONAL_IDENTIFIER_TYPE_CODE_MISC: t`Unspecified`,
  NATIONAL_IDENTIFIER_TYPE_CODE_FIIN: t`Foreign Investment Identity Number`,
  NATIONAL_IDENTIFIER_TYPE_CODE_SOCS: t`Social Security Number`,
  NATIONAL_IDENTIFIER_TYPE_CODE_CCPT: t`Passport Number`,
  NATIONAL_IDENTIFIER_TYPE_CODE_IDCD: t`Identity Card Number`,
  NATIONAL_IDENTIFIER_TYPE_CODE_DRLC: t`Driver's License Number`,
  NATIONAL_IDENTIFIER_TYPE_CODE_ARNU: t`Alien Registration Number`
};

export const disabledIdentifiers = [
  'NATIONAL_IDENTIFIER_TYPE_CODE_SOCS',
  'NATIONAL_IDENTIFIER_TYPE_CODE_FIIN',
  'NATIONAL_IDENTIFIER_TYPE_CODE_CCPT',
  'NATIONAL_IDENTIFIER_TYPE_CODE_IDCD',
  'NATIONAL_IDENTIFIER_TYPE_CODE_DRLC',
  'NATIONAL_IDENTIFIER_TYPE_CODE_ARNU'
];

export const getNationalIdentificationOptions = () =>
  Object.entries(NATIONAL_IDENTIFICATION).map(([k, v]) => ({
    value: k,
    label: v,
    isDisabled: disabledIdentifiers.includes(k)
  }));

export const getNationalIdentificationLabel = (nationalIdentifierTypeCode: any) => {
  return Object.entries(NATIONAL_IDENTIFICATION).map(([k, v]) => {
    if (k === nationalIdentifierTypeCode) {
      return v;
    }
  });
};
