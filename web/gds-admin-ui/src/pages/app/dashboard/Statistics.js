import React from 'react';
import { Card, Row, Col } from 'react-bootstrap';


const Statistics = ({ data }) => {


    return (
        <>
            <Row>
                <Col>
                    <Card className="widget-inline">
                        <Card.Body className="p-0">
                            <Row className="g-0">
                                <Col sm={6} xl={3}>
                                    <Card className="shadow-none m-0">
                                        <Card.Body className="text-center">
                                            <i className="dripicons-briefcase text-muted font-24"></i>
                                            <h3>
                                                {data?.vasps_count}
                                            </h3>
                                            <p className="text-muted font-15 mb-0">All VASPs</p>
                                        </Card.Body>
                                    </Card>
                                </Col>

                                <Col sm={6} xl={3}>
                                    <Card className="card shadow-none m-0 border-start">
                                        <Card.Body className="text-center">
                                            <i className="dripicons-checklist text-muted font-24"></i>
                                            <h3>
                                                <span>
                                                    {data?.pending_registrations}
                                                </span>
                                            </h3>
                                            <p className="text-muted font-15 mb-0">Pending Registrations</p>
                                        </Card.Body>
                                    </Card>
                                </Col>

                                <Col sm={6} xl={3}>
                                    <Card className="card shadow-none m-0 border-start">
                                        <Card.Body className="text-center">
                                            <i className="dripicons-user-group text-muted font-24"></i>
                                            <h3>
                                                <span>
                                                    {data?.verified_contacts}
                                                </span>
                                            </h3>
                                            <p className="text-muted font-15 mb-0">Verified Contacts</p>
                                        </Card.Body>
                                    </Card>
                                </Col>

                                <Col sm={6} xl={3}>
                                    <Card className="card shadow-none m-0 border-start">
                                        <Card.Body className="text-center">
                                            <i className="dripicons-copy text-muted font-24"></i>
                                            <h3>
                                                {data?.certificates_issued}
                                            </h3>
                                            <p className="text-muted font-15 mb-0">Certificates Issued</p>
                                        </Card.Body>
                                    </Card>
                                </Col>
                            </Row>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </>
    );
};

export default Statistics;
