
import React from 'react'
import { Col } from 'react-bootstrap'
import { NATIONAL_IDENTIFIER_TYPE } from '../../../../../constants/national-identification'
import { formatDisplayedData } from '../../../../../utils'
import countryCodeEmoji, { isoCountries } from '../../../../../utils/country'

function NationalIdentification({ data }) {
    return (
        <Col>
            {
                data ? (
                    <Col className="mt-3">
                        <p className="fw-bold mb-1">National identification</p>
                        <hr className='my-1' />
                        <p className="mb-2 fw-bold">Issued by: <span className="fw-normal">{`${formatDisplayedData(countryCodeEmoji(data?.country_of_issue))} (${formatDisplayedData(data?.country_of_issue)}) by authority ${formatDisplayedData(data?.registration_authority)}`}</span></p>
                        <p className="mb-1 fw-bold">National identification type: <span className="fw-normal badge bg-primary rounded-pill px-1">{formatDisplayedData(NATIONAL_IDENTIFIER_TYPE[data?.national_identifier_type])}</span></p>
                        <p className="mb-2 fw-bold">LEIX: <span className="fw-normal">{formatDisplayedData(data?.national_identifier)}</span></p>
                        <p className="mb-2 fw-bold">Country of registration: <span className="fw-normal">{formatDisplayedData(isoCountries[data?.country_of_issue])}</span></p>
                        <p className="mb-2 fw-bold">Customer number: <span className="fw-normal">{formatDisplayedData(data?.customer_number)}</span></p>
                    </Col>

                ) : (
                    <Col>
                        <p className="mb-1 fw-bold">National identification: <span className="fw-normal">{formatDisplayedData(data?.national_identification)}</span></p>
                    </Col>
                )
            }
        </Col>
    )
}

export default NationalIdentification
