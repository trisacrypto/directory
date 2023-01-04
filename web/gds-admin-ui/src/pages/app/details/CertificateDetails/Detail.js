import countryCodeEmoji, { getCountryName } from '@/utils/country';

function Detail({ data, title }) {
  const formatToReadableString = (data) => {
    if (Array.isArray(data)) {
      return !data.length ? 'N/A' : data.toString().split(',').join(', ');
    }
  };

  return (
      <>
      <p className="m-0">
        <span className="fw-bold">Common name:</span> {data?.common_name}
      </p>
      <p className="m-0">
        <span className="fw-bold">Country: </span>
              {formatToReadableString(data?.country.map((c) => `${countryCodeEmoji(c)} ${getCountryName(c)}`))}
        )}
      </p>
      <p className="m-0">
              <span className="fw-bold">Locality: </span>
        {formatToReadableString(data?.locality)}
      </p>
      <p className="m-0">
        <span className="fw-bold">Organization: </span>
              {formatToReadableString(data?.organization)}
      </p>
      <p className="m-0">
        <span className="fw-bold">Organizational unit: </span>
        {formatToReadableString(data?.organizational_unit)}
      </p>
      <p className="m-0">
        <span className="fw-bold">Postal code: </span>
        {formatToReadableString(data?.postal_code)}
      </p>
      <p className="m-0">
        <span className="fw-bold">Province: </span>
        {formatToReadableString(data?.province)}
      </p>
          <p className="m-0">
        <span className="fw-bold">Serial number: </span>
        {data?.serial_number ? data?.serial_number : 'N/A'}
      </p>
      <p className="m-0">
        <span className="fw-bold">Street address: </span>
        {data?.serial_number ? data?.serial_number : 'N/A'}
      </p>
    </>
  );
}

export default Detail;
