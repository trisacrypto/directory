import React from 'react'
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime'
dayjs.extend(relativeTime)



const CertificateExpirationColumn = ({ row }) => {

    return (
        <React.Fragment>
            <p className="m-0 d-inline-block align-middle font-16" data-testid="certificate_expiration">
                {row?.original?.certificate_expiration ? dayjs(row?.original?.certificate_expiration).format("MMM DD, YYYY h:mm:ss a") : 'N/A'}
            </p>
        </React.Fragment>
    );
};

export default CertificateExpirationColumn