import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import React from 'react';

dayjs.extend(relativeTime);

const CertificateExpirationColumn = ({ row }) => (
  <p className="m-0 d-inline-block align-middle font-16" data-testid="certificate_expiration">
    {row?.original?.certificate_expiration
      ? dayjs(row?.original?.certificate_expiration).format('MMM DD, YYYY h:mm:ss a')
      : 'N/A'}
  </p>
);

export default CertificateExpirationColumn;
