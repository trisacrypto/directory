export const NATIONAL_IDENTIFICATION = {
  NATIONAL_IDENTIFIER_TYPE_CODE_LEIX: 'Legal Entity Identifier (LEI)',
  NATIONAL_IDENTIFIER_TYPE_CODE_TXID: 'Tax Identification Number',
  NATIONAL_IDENTIFIER_TYPE_CODE_FIIN: 'Foreign Investment Identity Number',
  NATIONAL_IDENTIFIER_TYPE_CODE_SOCS: 'Social Security Number',
  NATIONAL_IDENTIFIER_TYPE_CODE_RAID: 'Registration Authority Identifier',
  NATIONAL_IDENTIFIER_TYPE_CODE_CCPT: 'Passport Number',
  NATIONAL_IDENTIFIER_TYPE_CODE_IDCD: 'Identity Card Number',
  NATIONAL_IDENTIFIER_TYPE_CODE_DRLC: "Driver's License Number",
  NATIONAL_IDENTIFIER_TYPE_CODE_ARNU: 'Alien Registration Number',
  NATIONAL_IDENTIFIER_TYPE_CODE_MISC: 'Unspecified'
};

export const getNationalIdentificationOptions = () =>
  Object.entries(NATIONAL_IDENTIFICATION).map(([k, v]) => ({ value: k, label: v }));

export const getNationalIdentificationLabel = (nationalIdentifierTypeCode: any) => {
  return Object.entries(NATIONAL_IDENTIFICATION).map(([k, v]) => {
    if (k === nationalIdentifierTypeCode) {
      return v;
    }
  });
};
