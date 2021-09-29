
import React from 'react'
import { Col, Row } from 'react-bootstrap';
import { formatDisplayedData } from '../../../../utils';

function Contact({ data, type }) {
    return (
        <>
            <p className="fw-bold">{type} contact:</p>
            <hr />
            <Row>
                <Col xl={6}>
                    <p className="fw-bold mb-2">Email: <span className="fw-normal">{formatDisplayedData(data?.email)}</span></p>
                    <p className="fw-bold mb-2">Name: <span className="fw-normal">{formatDisplayedData(data?.name)}</span></p>
                </Col>
                <Col xl={6}>
                    <p className="fw-bold mb-2">Person: <span className="fw-normal">{data?.person ? 'Has IVMS101 Record' : 'No IVMS101 Data'}</span></p>
                    <p className="fw-bold mb-2">Phone: <span className="fw-normal">{formatDisplayedData(data?.phone)}</span></p>
                </Col>
            </Row>
        </>
    )
}

export default Contact
