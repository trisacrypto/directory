import React from 'react';
import { Card, Col, Row } from 'react-bootstrap';
import Contact from './Contact';

export default function ContactList({ data }) {

    return (
        <Card>
            <Card.Body>
                <h4 className="mt-0 mb-3">Contact Information</h4>
                <Row>
                    <Col xl={6}>
                        <Contact data={data?.legal} type="Legal" />
                    </Col>
                    <Col xl={6}>
                        <Contact data={data?.administrative} type="Administrative" />
                    </Col>
                    <Col xl={6} className="mt-3">
                        <Contact data={data?.billing} type="Billing" />
                    </Col>
                    <Col xl={6} className="mt-3">
                        <Contact data={data?.technical} type="Technical" />
                    </Col>
                </Row>
            </Card.Body>
        </Card>
    );
};
