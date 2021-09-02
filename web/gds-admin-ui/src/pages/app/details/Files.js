// @flow
import React from 'react';
import { Card, Row } from 'react-bootstrap';


const Files = (): React$Element<any> => {
    return (
        <Card>
            <Card.Body>
                <h5 className="card-title mb-3">Files</h5>

                <Card className="mb-1 shadow-none border">
                    <div className="p-2">
                        <Row className="align-items-center">
                            <div className="col-auto">
                                <div className="avatar-sm">
                                    <span className="avatar-title rounded">.ZIP</span>
                                </div>
                            </div>
                            <div className="col ps-0">
                                <a href="/" className="text-muted fw-bold">
                                    Hyper-admin-design.zip
                                </a>
                                <p className="mb-0">2.3 MB</p>
                            </div>
                            <div className="col-auto">
                                <a href="/" className="btn btn-link btn-lg text-muted">
                                    <i className="dripicons-download"></i>
                                </a>
                            </div>
                        </Row>
                    </div>
                </Card>

                <Card className="mb-1 shadow-none border">
                    <div className="p-2">
                        <Row className="align-items-center">
                            <div className="col ps-0">
                                <a href="/" className="text-muted fw-bold">
                                    Dashboard-design.jpg
                                </a>
                                <p className="mb-0">3.5 MB</p>
                            </div>
                            <div className="col-auto">
                                <a href="/" className="btn btn-link btn-lg text-muted">
                                    <i className="dripicons-download"></i>
                                </a>
                            </div>
                        </Row>
                    </div>
                </Card>

                <Card className="mb-1 shadow-none border">
                    <div className="p-2">
                        <Row className="align-items-center">
                            <div className="col-auto">
                                <div className="avatar-sm">
                                    <span className="avatar-title bg-secondary rounded">.MP4</span>
                                </div>
                            </div>
                            <div className="col ps-0">
                                <a href="/" className="text-muted fw-bold">
                                    Admin-bug-report.mp4
                                </a>
                                <p className="mb-0">7.05 MB</p>
                            </div>
                            <div className="col-auto">
                                <a href="/" className="btn btn-link btn-lg text-muted">
                                    <i className="dripicons-download"></i>
                                </a>
                            </div>
                        </Row>
                    </div>
                </Card>
            </Card.Body>
        </Card>
    );
};

export default Files;
