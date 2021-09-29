
import React from 'react'
import { Col, Row } from 'react-bootstrap'
import { formatDisplayedData } from '../../../../utils'

function Name({ data }) {
    return (
        <Col className="mt-3">
            <p className="fw-bold mb-2">Name</p>
            <hr />
            <Row>
                {
                    data && data.local_name_identifiers.length ? (
                        <Col xl={6}>
                            <p className="fw-bold mb-2">Local Name Identifiers</p>
                            <hr />
                            {
                                data && data.local_name_identifiers ? data.local_name_identifiers.map((identifier, index) => (
                                    <div key={index}>
                                        <p className="fw-bold mb-2">Legal Person Name: <span className="fw-normal">{formatDisplayedData(identifier?.legal_person_name)}</span></p>
                                        <p className="fw-bold mb-2">Legal Person Name Identifier Type: <span className="fw-normal">{formatDisplayedData(identifier?.legal_person_name_identifier_type)}</span></p>
                                    </div>
                                )) : (<p className="fw-bold mb-2">Legal Person Name Identifier Type: <span className="fw-normal">"N/A"</span></p>)
                            }
                        </Col>
                    ) : (
                        <Col>
                            <p className="fw-bold mb-2">Local Name Identifiers: <span className="fw-normal">{formatDisplayedData(data?.local_name_identifiers)}</span></p>
                        </Col>
                    )
                }
                {
                    data && data.phonetic_name_identifiers.length ? (
                        <Col xl={6}>
                            <p className="fw-bold mb-2">Phonetic Name Identifiers:</p>
                            <hr />
                            {
                                data && data.phonetic_name_identifiers ? data.phonetic_name_identifiers.map((identifier, index) => (
                                    <div key={index}>
                                        <p className="fw-bold mb-2">Legal Person Name: <span className="fw-normal">{formatDisplayedData(identifier?.legal_person_name)}</span></p>
                                        <p className="fw-bold mb-2">Legal Person Name Identifier Type: <span className="fw-normal">{formatDisplayedData(identifier?.legal_person_name_identifier_type)}</span></p>
                                    </div>
                                )) : (<p className="fw-bold mb-2">Legal Person Name Identifier Type: <span className="fw-normal">"N/A"</span></p>)
                            }
                        </Col>
                    ) : (
                        <Col>
                            <p className="fw-bold mb-2">Phonetic Name Identifiers: <span className="fw-normal">{formatDisplayedData(data?.phonetic_name_identifiers)}</span></p>
                        </Col>
                    )
                }


                {
                    data && data.name_identifiers.length ? (
                        <Col xl={6}>
                            <p className="fw-bold mb-2">Name Identifiers:</p>
                            <hr />
                            {
                                data && data.name_identifiers ? data.name_identifiers.map((identifier, index) => (
                                    <div key={index}>
                                        <p className="fw-bold mb-2">Legal Person Name: <span className="fw-normal">{formatDisplayedData(identifier?.legal_person_name)}</span></p>
                                        <p className="fw-bold mb-2">Legal Person Name Identifier Type: <span className="fw-normal">{formatDisplayedData(identifier?.legal_person_name_identifier_type)}</span></p>
                                    </div>
                                )) : (
                                    <Col>
                                        <p className="fw-bold mb-2">Legal Person Name Identifier Type: <span className="fw-normal">"N/A"</span></p>
                                    </Col>
                                )
                            }
                        </Col>
                    ) : (
                        <Col>
                            <p className="fw-bold mb-2">Name Identifiers: <span className="fw-normal">{formatDisplayedData(data?.name_identifiers)}</span></p>
                        </Col>
                    )
                }

            </Row>
        </Col>
    )
}

export default Name
