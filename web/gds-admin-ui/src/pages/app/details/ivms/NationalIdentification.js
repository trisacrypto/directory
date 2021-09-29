
import React from 'react'
import { Col } from 'react-bootstrap'
import { formatDisplayedData } from '../../../../utils'

function NationalIdentification({ data }) {
    return (
        <>
            {
                data ? (
                    <Col lg={6} className="mt-3">
                        <p className="fw-bold mb-2">National Identification</p>
                        <hr />
                        <p className="fw-bold mb-2">Country of Issue: <span className="fw-normal">{formatDisplayedData(data?.country_of_issue)}</span></p>
                        <p className="fw-bold mb-2">National Identifier: <span className="fw-normal">{formatDisplayedData(data?.national_identifier)}</span></p>
                        <p className="fw-bold mb-2">National Identification Type: <span className="fw-normal">{formatDisplayedData(data?.national_identifier_type)}</span></p>
                        <p className="fw-bold mb-2">Registration Authority: <span className="fw-normal">{formatDisplayedData(data?.registration_authority)}</span></p>
                    </Col>

                ) : (
                    <Col>
                        <p className="fw-bold mb-2">National Identification: <span className="fw-normal">{formatDisplayedData(data?.national_identification)}</span></p>
                    </Col>
                )
            }
        </>
    )
}

export default NationalIdentification
