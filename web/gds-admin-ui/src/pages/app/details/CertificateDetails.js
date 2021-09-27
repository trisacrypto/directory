
import React from 'react'
import { Card, Col, Row } from 'react-bootstrap';
import { formatDisplayedData } from '../../../utils';

function CertificateDetails({ data }) {
    return (
        <Card>
            <Card.Body>
                <h4 className="mt-0 mb-3">Certificate details</h4>
                <p className="fw-bold mb-2">Signature Algorithm: <span className="fw-normal">{formatDisplayedData(data?.signature_algorithm)}</span></p>
                <p className="fw-bold mb-2">Version: <span className="fw-normal">{formatDisplayedData(data?.version)}</span></p>
                <p className="fw-bold mb-2">Not Before: <span className="fw-normal">{formatDisplayedData(data?.not_before)}</span></p>
                <p className="fw-bold mb-2">Not After: <span className="fw-normal">{formatDisplayedData(data?.not_after)}</span></p>
                <p className="fw-bold mb-2">Public Key Algorithm: <span className="fw-normal">{formatDisplayedData(data?.public_key_algorithm)}</span></p>
                <p className="fw-bold mb-2">Revoked: <span className="fw-normal">{formatDisplayedData(data?.revoked)}</span></p>
                <p className="fw-bold mb-2">Serial Number: <span className="fw-normal">{formatDisplayedData(data?.serial_number)}</span></p>

                <Row>
                    <Col>
                        <p className="fw-bold mb-2 mt-3">Issuer</p>
                        <hr />
                        <p className="fw-bold mb-2">Common Name: <span className="fw-normal">{formatDisplayedData(data?.issuer?.common_name)}</span></p>
                        <p className="fw-bold mb-2">Country: <span className="fw-normal">{formatDisplayedData(data?.issuer?.country)}</span></p>
                        <p className="fw-bold mb-2">Locality: <span className="fw-normal">{formatDisplayedData(data?.issuer?.locality
                        )}</span></p>
                        <p className="fw-bold mb-2">Organisation: <span className="fw-normal">{formatDisplayedData(data?.issuer?.organization)}</span></p>
                        <p className="fw-bold mb-2">Organisation Unit: <span className="fw-normal">{formatDisplayedData(data?.issuer?.organizational_unit)}</span></p>
                        <p className="fw-bold mb-2">Postal Code: <span className="fw-normal">{formatDisplayedData(data?.issuer?.postal_code)}</span></p>
                        <p className="fw-bold mb-2">Province: <span className="fw-normal">{formatDisplayedData(data?.issuer?.province)}</span></p>
                        <p className="fw-bold mb-2">Serial Number: <span className="fw-normal">{formatDisplayedData(data?.issuer?.serial_number)}</span></p>
                        <p className="fw-bold mb-2">Street Adress: <span className="fw-normal">{formatDisplayedData(data?.issuer?.street_address)}</span></p>
                    </Col>
                    <Col>
                        <p className="fw-bold mb-2 mt-3">Subject</p>
                        <hr />
                        <p className="fw-bold mb-2">Common Name: <span className="fw-normal">{formatDisplayedData(data?.subject?.common_name)}</span></p>
                        <p className="fw-bold mb-2">Country: <span className="fw-normal">{formatDisplayedData(data?.subject?.country)}</span></p>
                        <p className="fw-bold mb-2">Locality: <span className="fw-normal">{formatDisplayedData(data?.subject?.locality)}</span></p>
                        <p className="fw-bold mb-2">Organisation: <span className="fw-normal">{formatDisplayedData(data?.subject?.organization)}</span></p>
                        <p className="fw-bold mb-2">Organisation Unit: <span className="fw-normal">{formatDisplayedData(data?.subject?.organizational_unit)}</span></p>
                        <p className="fw-bold mb-2">Postal Code: <span className="fw-normal">{formatDisplayedData(data?.subject?.postal_code)}</span></p>
                        <p className="fw-bold mb-2">Province: <span className="fw-normal">{formatDisplayedData(data?.subject?.province)}</span></p>
                        <p className="fw-bold mb-2">Serial Number: <span className="fw-normal">{formatDisplayedData(data?.serial_number)}</span></p>
                        <p className="fw-bold mb-2">Street Adress: <span className="fw-normal">{formatDisplayedData(data?.subject?.street_address)}</span></p>
                    </Col>

                </Row>
                {/* <p className="fw-bold mb-2">Subject: <span className="fw-normal">{formatDisplayedData(data?.subject)}</span></p> */}
                <p className="fw-bold mb-2">Signature: <span className="fw-normal">{formatDisplayedData(data?.signature)}</span></p>
                <p className="fw-bold mb-2">Chain: <span className="fw-normal">{formatDisplayedData(data?.chain)}</span></p>
                <p className="fw-bold mb-2">Data: <span className="fw-normal">{formatDisplayedData(data?.data)}</span></p>
            </Card.Body>
        </Card>
    )
}

export default CertificateDetails
