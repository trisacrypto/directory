import { t } from '@lingui/macro';

export const addressType = {
  ADDRESS_TYPE_CODE_MISC: t`Unspecified`,
  ADDRESS_TYPE_CODE_HOME: t`Residential`,
  ADDRESS_TYPE_CODE_BIZZ: t`Business`,
  ADDRESS_TYPE_CODE_GEOG: t`Geographic`
};

export const addressTypeOptions = () => {
  return Object.entries(addressType).map(([k, v]) => ({ value: k, label: v }));
};
