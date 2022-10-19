
import React from 'react'
import { Card, Col, Row } from 'react-bootstrap';
import { Status, StatusLabel } from 'constants/index';
import { formatDisplayedData, getStatusClassName, isValidHttpUrl } from 'utils';
import dayjs from 'dayjs';
import Name from './components/Name';
import NationalIdentification from './components/NationalIdentification';
import { BUSINESS_CATEGORY } from 'constants/basic-details';
import Geographic from './components/Geographic';
import countryCodeEmoji from 'utils/country';
import { downloadFile } from 'helpers/api/utils';
import classNames from 'classnames';
import BasicDetailsDropDown from './components/BasicDetailsDropdown';
import PropTypes from 'prop-types'
import TrisaDetails from './components/TrisaDetails';


function BasicDetails({ data }) {
    const userLang = navigator.language || navigator.userLanguage;
    const formatDate = (date) => {
        if (date) {
            const dateObj = new Date(date);
            return dateObj.toLocaleDateString(userLang, {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit'
            });
        }
        return 'N/A';
    };
    const isNotPendingReview = () => data?.vasp?.verification_status !== Status.PENDING_REVIEW

    const handleIvmsJsonExportClick = () => {
        if (data && data.vasp && data.vasp.entity) {
            const filename = `${dayjs().format("YYYY-MM-DD")}-ivms.json`
            const mime = `data:text/json;charset=utf-8`
            const file = JSON.stringify(data.vasp?.entity)

            downloadFile(file, filename, mime)
        }
    }

    const handleTrisaJsonExportClick = () => {
        if (data && data.vasp && data.vasp.entity) {
            const filename = `${dayjs().format("YYYY-MM-DD")}-trisa.json`
            const mime = `data:text/json;charset=utf-8`

            const trisaData = {
                id: data?.vasp?.id,
                common_name: data?.vasp?.common_name,
                trisa_endpoint: data?.vasp?.trisa_endpoint,
                verified_on: data?.vasp?.verified_on,
                registered_directory: data?.vasp?.registered_directory
            }
            const file = JSON.stringify(trisaData)

            downloadFile(file, filename, mime)
        }
    }

    return (
        <>
            <Card className="d-block">
                <Card.Body>
                    <BasicDetailsDropDown isNotPendingReview={isNotPendingReview} vasp={data} />
                    <div>
                        <div>
                            <h3 className="m-0 d-inline-block text-dark">{data?.name}</h3>
                            {data?.traveler ? <span className='badge bg-primary rounded-pill px-1 ms-1 align-text-bottom'>Traveler</span> : null}
                            {data?.vasp?.verification_status ? <span className={classNames('badge rounded-pill px-1 ms-1 align-text-bottom', getStatusClassName(data?.vasp?.verification_status))}>{StatusLabel[data?.vasp?.verification_status]}</span> : null}
                        </div>
                        <div className='d-flex align-items-center'>
                            <span className="fw-normal d-block me-1" style={{ fontSize: '2rem' }} data-testid="country-flag">{countryCodeEmoji(data?.vasp?.entity?.country_of_registration || data?.vasp?.entity?.geographic_addresses[0]?.country)}</span>
                            {isValidHttpUrl(data?.vasp?.website) ? <a target="_blank" href={`${data?.vasp?.website}`} rel="noreferrer">{data?.vasp?.website}</a> : null}
                        </div>
                    </div>
                    <Row>
                        <Col>
                            <h4 className='text-dark mb-0'>Business details <button onClick={handleIvmsJsonExportClick} className='mdi mdi-arrow-down-bold-circle-outline border-0 bg-transparent' title="Download as JSON"></button></h4>
                            <p className="mb-2">
                                {
                                    Array.isArray(data?.vasp?.vasp_categories) && data?.vasp?.vasp_categories.map((category, index) => <span key={index} className='badge bg-success rounded-pill px-1 me-1 fw-normal'>{category}</span>)
                                }
                            </p>
                            <hr className='m-0' />
                            <Row>
                                <Col>
                                    <Name data={data?.vasp?.entity?.name} />
                                    <NationalIdentification data={data?.vasp?.entity} />
                                </Col>
                                <Col>
                                    <p className="mb-2 mt-md-3 mt-lg-3 fw-bold">Established on: <span className="fw-normal">{formatDisplayedData(data?.vasp?.established_on)}</span></p>
                                    <h5 className='mt-3'>Address(es):</h5>
                                    <hr className='m-0 mb-1' />
                                    <Geographic data={data?.vasp?.entity?.geographic_addresses || []} />
                                </Col>
                            </Row>
                            <Col>
                                <p className="mb-2 fw-bold">Business Category/Categories: <span className="badge bg-primary rounded-pill px-1">{BUSINESS_CATEGORY[data?.vasp?.business_category]}</span></p>
                            </Col>
                            <Col className='mt-4'>
                                <TrisaDetails data={data} handleTrisaJsonExportClick={handleTrisaJsonExportClick} />
                            </Col>
                        </Col>
                        <Col sm={12} className='d-flex justify-content-around flex-sm-wrap flex-md-nowrap text-center'>
                            <p className="fw-bold mb-2 text-dark"> <span className='d-block'>First listed</span> <span className="fw-normal">{formatDate(data?.vasp?.first_listed)}</span></p>
                            <p className="fw-bold mb-2 text-dark"> <span className='d-block'>Last updated</span> <span className="fw-normal">{formatDate(data?.vasp?.last_updated)}</span></p>
                            <p className="fw-bold mb-2 text-dark"> <span className='d-block'>Verified on</span> <span className="fw-normal">{formatDate(data?.vasp?.verified_on)}</span></p>
                        </Col>
                    </Row>
                    <Col sm={12}>
                        <hr />
                        <p className="mb-2 text-center text-muted ">Version: {`${formatDisplayedData(data?.vasp?.version?.version)}.${formatDisplayedData(data?.vasp?.version?.pid)}`}</p>
                    </Col>
                </Card.Body>
            </Card>
        </>

    )
}

BasicDetails.propTypes = {
    data: PropTypes.object.isRequired
}

export default BasicDetails
