
import React from 'react'
import { Card, Col, Dropdown, Row } from 'react-bootstrap';
import { formatDisplayedData } from '../../../utils';


function BasicDetails({ data }) {

    return (
        <Card className="d-block">
            <Card.Body>
                <Dropdown className="float-end" align="end">
                    <Dropdown.Toggle
                        variant="link"
                        tag="a"
                        className="card-drop arrow-none cursor-pointer p-0 shadow-none">
                        <i className="dripicons-dots-3"></i>
                    </Dropdown.Toggle>
                    <Dropdown.Menu>
                        <Dropdown.Item>
                            <i className="mdi mdi-pencil me-1"></i>Edit Details
                        </Dropdown.Item>
                        <Dropdown.Item>
                            <i className="mdi mdi-pencil me-1"></i>Edit TRIXO Implementation
                        </Dropdown.Item>
                    </Dropdown.Menu>
                </Dropdown>

                <h4 className="mt-0 mb-3">Basic Details</h4>
                <Row>
                    <Col xl={6}>
                        <p className="fw-bold mb-2">ID: <span className="fw-normal">{formatDisplayedData(data?.vasp?.id)}</span></p>
                        <p className="fw-bold mb-2">Name: <span className="fw-normal">{formatDisplayedData(data?.name)}</span></p>
                        <p className="fw-bold mb-2">Common Name: <span className="fw-normal">{formatDisplayedData(data?.vasp?.common_name)}</span></p>
                        <p className="fw-bold mb-2">Verification Status: <span className="fw-normal">{formatDisplayedData(data?.vasp?.verification_status)}</span></p>
                    </Col>
                    <Col xl={6}>
                        <p className="fw-bold mb-2">TRISA Endpoint: <span className="fw-normal">{formatDisplayedData(data?.vasp?.trisa_endpoint)}</span></p>
                        <p className="fw-bold mb-2">Website: <span className="fw-normal">{formatDisplayedData(data?.vasp?.website)}</span></p>
                        <p className="fw-bold mb-2">Established On: <span className="fw-normal">{formatDisplayedData(data?.vasp?.established_on)}</span></p>
                        <p className="fw-bold mb-2">Verified On: <span className="fw-normal">{formatDisplayedData(data?.vasp?.verified_on)}</span></p>
                    </Col>
                    <Col xl={6}>
                        <p className="fw-bold mb-2">Business categories: <span className="fw-normal">{formatDisplayedData(data?.vasp?.business_category)}</span></p>
                        <p className="fw-bold mb-2">First Listed: <span className="fw-normal">{formatDisplayedData(data?.vasp?.first_listed)}</span></p>
                        <p className="fw-bold mb-2">Last Update: <span className="fw-normal">{formatDisplayedData(data?.vasp?.last_updated)}</span></p>
                        <p className="fw-bold mb-2">Registered Directory: <span className="fw-normal">{formatDisplayedData(data?.vasp?.registered_directory)}</span></p>
                        <p className="fw-bold mb-2">VASP Category: <span className="fw-normal">{formatDisplayedData(data?.vasp?.vasp_categories)}</span></p>
                    </Col>
                    <Col xl={6}>
                        <p className="fw-bold mb-2 mt-3">Version</p>
                        <hr />
                        <p className="fw-bold mb-2">PID: <span className="fw-normal">{formatDisplayedData(data?.vasp?.version?.pid)}</span></p>
                        <p className="fw-bold mb-2">Version: <span className="fw-normal">{formatDisplayedData(data?.vasp?.version?.version)}</span></p>
                    </Col>
                </Row>

            </Card.Body>
        </Card>
    )
}

export default BasicDetails
