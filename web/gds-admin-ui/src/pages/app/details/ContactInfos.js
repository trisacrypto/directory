// @flow
import React from 'react';
import { Card, Col, Row } from 'react-bootstrap';
import { formatDisplayedData } from "../../../utils"

const ContactInfos = ({ data }) => {

    return (
        <Card>
            <Card.Body>
                <h4 className="mt-0 mb-3">Contact Informations</h4>
                <Row>
                    <Col xl={6}>
                        <p className="fw-bold">Technical contact</p>
                        <hr />
                        <Row>
                            <Col xl={6}>
                                <p className="fw-bold mb-2">Email: <span className="fw-normal">{formatDisplayedData(data?.technical?.email)}</span></p>
                                <p className="fw-bold mb-2">Name: <span className="fw-normal">{formatDisplayedData(data?.technical?.name)}</span></p>
                            </Col>
                            <Col xl={6}>
                                <p className="fw-bold mb-2">Person: <span className="fw-normal">{formatDisplayedData(data?.technical?.person)}</span></p>
                                <p className="fw-bold mb-2">Phone: <span className="fw-normal">{formatDisplayedData(data?.technical?.phone)}</span></p>
                            </Col>
                        </Row>
                    </Col>
                    <Col xl={6}>
                        <p className="fw-bold">Legal contact :</p>
                        <hr />
                        <Row>
                            <Col xl={6}>
                                <p className="fw-bold mb-2">Email: <span className="fw-normal">{formatDisplayedData(data?.technical?.email)}</span></p>
                                <p className="fw-bold mb-2">Name: <span className="fw-normal">{formatDisplayedData(data?.technical?.name)}</span></p>
                            </Col>
                            <Col xl={6}>
                                <p className="fw-bold mb-2">Person: <span className="fw-normal">{formatDisplayedData(data?.technical?.person)}</span></p>
                                <p className="fw-bold mb-2">Email: <span className="fw-normal">{formatDisplayedData(data?.technical?.phone)}</span></p>
                            </Col>
                        </Row>
                    </Col>
                    <Col xl={6} className="mt-3">
                        <p className="fw-bold">Administrative contact :</p>
                        <hr />
                        <Row>
                            <Col xl={6}>
                                <p className="fw-bold mb-2">Email: <span className="fw-normal">{formatDisplayedData(data?.technical?.email)}</span></p>
                                <p className="fw-bold mb-2">Name: <span className="fw-normal">{formatDisplayedData(data?.technical?.name)}</span></p>
                            </Col>
                            <Col xl={6}>
                                <p className="fw-bold mb-2">Person: <span className="fw-normal">{formatDisplayedData(data?.technical?.person)}</span></p>
                                <p className="fw-bold mb-2">Email: <span className="fw-normal">{formatDisplayedData(data?.technical?.phone)}</span></p>
                            </Col>
                        </Row>
                    </Col>
                    <Col xl={6} className="mt-3">
                        <p className="fw-bold">Billing contact :</p>
                        <hr />
                        <Row>
                            <Col xl={6}>
                                <p className="fw-bold mb-2">Email: <span className="fw-normal">{formatDisplayedData(data?.technical?.email)}</span></p>
                                <p className="fw-bold mb-2">Name: <span className="fw-normal">{formatDisplayedData(data?.technical?.name)}</span></p>
                            </Col>
                            <Col xl={6}>
                                <p className="fw-bold mb-2">Person: <span className="fw-normal">{formatDisplayedData(data?.technical?.person)}</span></p>
                                <p className="fw-bold mb-2">Email: <span className="fw-normal">{formatDisplayedData(data?.technical?.phone)}</span></p>
                            </Col>
                        </Row>
                    </Col>
                </Row>
            </Card.Body>
        </Card>
    );
};

export default ContactInfos;
