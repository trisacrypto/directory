export const isValidIvmsAddress = (address: any) => {
  if (address) {
    return !!(address.country && address.address_type);
  }
  return false;
};

export const hasAddressLine = (address: any) => {
  if (isValidIvmsAddress(address)) {
    return Array.isArray(address.address_line) && address.address_line.length > 0;
  }
  return false;
};

export const hasAddressField = (address: any) => {
  if (isValidIvmsAddress(address) && !hasAddressLine(address)) {
    return !!(address.street_name && (address.building_number || address.building_name));
  }
  return false;
};

const hasAddressFieldAndLine = (address: any) => {
  if (hasAddressField(address) && hasAddressLine(address)) {
    console.warn('[ERROR]', 'cannot render address');
    return true;
  }
  return false;
};

export const renderLines = (address: any) => (
  <address data-testid="addressLine">
    {address.address_line.map(
      (addressLine: any, index: number) => addressLine && <div key={index}>{addressLine} </div>
    )}
    <div>{address?.country}</div>
  </address>
);

export const renderField = (address: any) => (
  <address data-testid="addressField">
    {address.sub_department ? (
      <>
        {address?.sub_department} <br />
      </>
    ) : null}
    {address.department ? (
      <>
        {address?.department} <br />
      </>
    ) : null}
    {address.building_number} {address?.street_name}
    <br />
    {address.post_box ? (
      <>
        P.O. Box: {address?.post_box} <br />
      </>
    ) : null}
    {address.floor || address.room || address.building_name ? (
      <>
        {address?.floor} {address?.room} {address?.building_name} <br />
      </>
    ) : null}
    {address.district_name ? (
      <>
        {address?.district_name} <br />
      </>
    ) : null}
    {address.town_name || address.town_location_name || address.country_sub_division ? (
      <>
        {address?.town_name} {address?.town_location_name} {address?.country_sub_division}{' '}
        {address?.post_code} <br />
      </>
    ) : null}
    {address?.country}
  </address>
);

export const renderAddress = (address: any) => {
  if (hasAddressFieldAndLine(address)) {
    console.warn('[ERROR]', 'invalid address with both fields and lines');
    return <div>Invalid Address</div>;
  }

  if (hasAddressLine(address)) {
    return renderLines(address);
  }

  if (hasAddressField(address)) {
    return renderField(address);
  }

  console.warn('[ERROR]', 'could not render address');
  return <div>Unparseable Address</div>;
};
