export const NAME_IDENTIFIER_TYPE = {
  LEGAL_PERSON_NAME_TYPE_CODE_LEGL: 'Legal',
  LEGAL_PERSON_NAME_TYPE_CODE_SHRT: 'Short',
  LEGAL_PERSON_NAME_TYPE_CODE_TRAD: 'Trading',
  LEGAL_PERSON_NAME_TYPE_CODE_MISC: 'Misc'
};

export const getNameIdentiferTypeOptions = () => {
  return Object.entries(NAME_IDENTIFIER_TYPE).map(([k, v], index) => ({
    label: v,
    value: k
  }));
};

export const getNameIdentiferTypeLabel = (value: any) => {
  return Object.entries(NAME_IDENTIFIER_TYPE).map(([k, v]) => {
    if (k === value) {
      return v;
    }
  });
};
