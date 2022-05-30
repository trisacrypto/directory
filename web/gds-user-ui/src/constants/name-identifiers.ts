import { t } from '@lingui/macro';

export const NAME_IDENTIFIER_TYPE = {
  LEGAL_PERSON_NAME_TYPE_CODE_LEGL: t`Legal`,
  LEGAL_PERSON_NAME_TYPE_CODE_SHRT: t`Short`,
  LEGAL_PERSON_NAME_TYPE_CODE_TRAD: t`Trading`,
  LEGAL_PERSON_NAME_TYPE_CODE_MISC: t`Misc`
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
