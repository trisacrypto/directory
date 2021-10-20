
import React from 'react'
import { Card, Col, Dropdown, Row } from 'react-bootstrap';
import { Status, StatusLabel } from '../../../constants';
import { formatDisplayedData, isValidHttpUrl } from '../../../utils';
import dayjs from 'dayjs';


function BasicDetails({ data }) {
    console.log('data', data)

    const formatDate = (date) => date ? dayjs(date).format('DD-MM-YYYY') : 'N/A';
    const isNotPendingReview = () => data?.vasp?.verification_status !== Status.PENDING_REVIEW

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
                        <Dropdown.Item disabled={isNotPendingReview()}>
                            <i className="mdi mdi-card-search me-1"></i>Review
                        </Dropdown.Item>
                        <Dropdown.Item>
                            <i className="mdi mdi-square-edit-outline me-1"></i>Edit
                        </Dropdown.Item>
                        <Dropdown.Item>
                            <i className="mdi mdi-printer me-1"></i>Print
                        </Dropdown.Item>
                        <Dropdown.Item>
                            <i className="mdi mdi-email me-1"></i>Resend
                        </Dropdown.Item>
                    </Dropdown.Menu>
                </Dropdown>

                <div className='mb-3'>
                    <div className='d-flex align-items-lg-start gap-2'>
                        <h4 className="mt-0">{data?.name}</h4>
                        {data?.traveler ? <span className='badge bg-primary rounded-pill px-1'>traveler</span> : null}
                        {data?.vasp?.verification_status ? <span className='badge bg-warning rounded-pill px-1'>{StatusLabel[data?.vasp?.verification_status]}</span> : null}
                    </div>
                    {isValidHttpUrl(data?.vasp?.website) ? <a target="_blank" href={`${data?.vasp?.website}`} rel="noreferrer">{data?.vasp?.website}</a> : null}
                </div>
                <Row>
                    <Col>
                        <div>
                            <h5>Business details <button className='mdi mdi-arrow-down-bold-circle-outline border-0 bg-transparent' title="Download as JSON"></button></h5>
                            <hr className='my-1' />
                            <p className="mb-2">Business categorie(s): <span className="fw-normal">{formatDisplayedData(data?.vasp?.business_category)}</span></p>
                            <p className="mb-2">Established on: <span className="fw-normal">{formatDisplayedData(data?.vasp?.established_on)}</span></p>
                            <p className="mb-2">VASP {data?.vasp?.vasp_categories.length > 1 ? 'categories' : 'category'}: <span className="fw-normal">{formatDisplayedData(data?.vasp?.vasp_categories)}</span></p>
                        </div>
                        <div className='mt-4'>
                            <h5>TRISA Details <button className='mdi mdi-arrow-down-bold-circle-outline border-0 bg-transparent' title="Download as JSON"></button></h5>
                            <hr className='my-1' />
                            <p className="mb-2">ID: <span className="fw-normal">{formatDisplayedData(data?.vasp?.id)}</span></p>
                            <p className="mb-2">Common name: <span className="fw-normal">{formatDisplayedData(data?.vasp?.common_name)}</span></p>
                            <p className="mb-2">TRISA Endpoint: <span className="fw-normal">{formatDisplayedData(data?.vasp?.trisa_endpoint)}</span></p>
                            <p className="mb-2">Registered Directory: <span className="fw-normal">{formatDisplayedData(data?.vasp?.registered_directory)}</span></p>
                        </div>
                    </Col>
                    <Col xl={6}>

                    </Col>
                    <Col sm={12} className='d-flex justify-content-around flex-sm-wrap flex-md-nowrap text-center'>
                        <p className="fw-bold mb-2"> <span className='d-block'>First Listed</span> <span className="fw-normal">{formatDate(data?.vasp?.first_listed)}</span></p>
                        <p className="fw-bold mb-2"> <span className='d-block'>Last Updated</span> <span className="fw-normal">{formatDate(data?.vasp?.last_updated)}</span></p>
                        <p className="fw-bold mb-2"> <span className='d-block'>Verified On</span> <span className="fw-normal">{formatDate(data?.vasp?.verified_on)}</span></p>
                    </Col>
                    <Col sm={12}>
                        <hr />
                        <p className="mb-2 text-center text-muted">Version: {`${formatDisplayedData(data?.vasp?.version?.version)}.${formatDisplayedData(data?.vasp?.version?.pid)}`}</p>
                    </Col>
                </Row>

            </Card.Body>
        </Card>
    )
}

export default BasicDetails
