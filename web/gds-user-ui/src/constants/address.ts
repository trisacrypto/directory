export const addressType = {
  ADDRESS_TYPE_CODE_MISC: 'Unspecified',
  ADDRESS_TYPE_CODE_HOME: 'Residential',
  ADDRESS_TYPE_CODE_BIZZ: 'Business',
  ADDRESS_TYPE_CODE_GEOG: 'Geographic'
};

export const addressTypeOptions = () => {
  return Object.entries(addressType).map(([k, v]) => ({ value: k, label: v }));
};
