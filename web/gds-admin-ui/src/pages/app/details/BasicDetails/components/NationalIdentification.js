
import React from 'react'
import { Col } from 'react-bootstrap'
import { NATIONAL_IDENTIFIER_TYPE } from 'constants/national-identification'
import { formatDisplayedData } from 'utils'
import countryCodeEmoji, { isoCountries } from 'utils/country'
import PropTypes from 'prop-types';


  function getDataFormatted(data){
      const issued_country_emoji = data?.national_identification?.country_of_issue ? formatDisplayedData(countryCodeEmoji(data.national_identification.country_of_issue)) : '';
      const issued_country_code = formatDisplayedData(data?.national_identification?.country_of_issue);
      const issued_authority =  formatDisplayedData(data?.national_identification?.registration_authority);
      const nat_ident_type = data?.national_identification?.national_identifier_type ? formatDisplayedData(NATIONAL_IDENTIFIER_TYPE[data.national_identification.national_identifier_type]) : 'N/A';
      const leix = formatDisplayedData(data?.national_identification?.national_identifier)
      const country_of_issue = data?.national_identification?.country_of_issue ? formatDisplayedData(isoCountries[data.national_identification.country_of_issue]) : 'N/A';
      const customer_number = formatDisplayedData(data?.customer_number);
        return {
            issued_country_emoji,
            issued_country_code,
            issued_authority,
            nat_ident_type,
            country_of_issue,
            leix,
            customer_number
        };
    };
function NationalIdentification({ data }) {
 
    return (
        <Col>
            {data ? (
                <Col className="mt-3">
                    <p className="fw-bold mb-1">National Identification</p>
                    <hr className="my-1" />
                    <p className="mb-2 fw-bold">
                        Issued by :{' '}
                        <span className="fw-normal">
                            {`${getDataFormatted(data).issued_country_emoji} 
                            (${getDataFormatted(data).issued_country_code})
                             by authority ${getDataFormatted(data).issued_authority}`
                            } `
                        </span>
                    </p>
                    <p className="mb-1 fw-bold">
                        National identification type:{' '}
                        <span className="fw-normal badge bg-primary rounded-pill px-1">
                            {getDataFormatted(data).nat_ident_type}
                        </span>
                    </p>
                    <p className="mb-2 fw-bold">
                        LEIX:{' '}
                        <span className="fw-normal">
                            {getDataFormatted(data).leix}
                        </span>
                    </p>
                    <p className="mb-2 fw-bold">
                        Country of registration:{' '}
                        <span className="fw-normal">
                            {getDataFormatted(data).country_of_issue}
                        </span>
                    </p>
                    <p className="mb-2 fw-bold">
                        Customer number: <span className="fw-normal">{getDataFormatted(data).customer_number}</span>
                    </p>
                </Col>
            ) : (
                <Col>
                    <p className="mb-1 fw-bold">
                        National Identification:{' '}
                        <span className="fw-normal">{formatDisplayedData(data?.national_identification)}</span>
                    </p>
                </Col>
            )}
        </Col>
    );
}

NationalIdentification.propTypes = {
    data: PropTypes.object
}

export default NationalIdentification
