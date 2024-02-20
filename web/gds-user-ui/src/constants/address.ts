import { t } from '@lingui/macro';

// Use address type where users will be instructed to select a type (ie. address form component).
export const addressType = {
  ADDRESS_TYPE_CODE_MISC: t`Unspecified`,
  ADDRESS_TYPE_CODE_HOME: t`Residential`,
  ADDRESS_TYPE_CODE_BIZZ: t`Business`,
  ADDRESS_TYPE_CODE_GEOG: t`Geographic`
};

// Use addressTypeEnum when the address type code is required when sending data to the backend.
export const addressTypeEnum = {
  ADDRESS_TYPE_MISC: 'ADDRESS_TYPE_CODE_MISC',
  ADDRESS_TYPE_HOME: 'ADDRESS_TYPE_CODE_HOME',
  ADDRESS_TYPE_BIZZ: 'ADDRESS_TYPE_CODE_BIZZ',
  ADDRESS_TYPE_GEOG: 'ADDRESS_TYPE_CODE_GEOG'
};

// The addressTypeOptions is used to generate the options for the address type select form control.
export const addressTypeOptions = () => {
  return Object.entries(addressType).map(([k, v]) => ({ value: k, label: v }));
};
