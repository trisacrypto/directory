
import React from 'react'
import FileInformationCard from 'components/FileInformationCard';
import PropTypes from 'prop-types'
import { Card, Col, Dropdown, Row } from 'react-bootstrap';
import { copyToClipboard, formatDisplayedData } from 'utils';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime'
dayjs.extend(relativeTime)


export const IdentityCertificateDropDown = ({ handleCopySignatureClick, handleCopySerialNumberClick, handleViewDetailsClick }) => {

    return (
        <Dropdown className="float-end" align="end">
            <Dropdown.Toggle
                data-testid="certificate-details-3-dots"
                variant="link"
                tag="a"
                className="card-drop arrow-none cursor-pointer p-0 shadow-none">
                <i className="dripicons-dots-3"></i>
            </Dropdown.Toggle>
            <Dropdown.Menu>
                <Dropdown.Item data-testid="copy-signature" onClick={handleCopySignatureClick}>
                    <i className="mdi mdi-content-copy me-1"></i>Copy signature
                </Dropdown.Item>
                <Dropdown.Item data-testid="copy-serial-number" onClick={handleCopySerialNumberClick}>
                    <i className="mdi mdi-content-copy me-1"></i>Copy serial number
                </Dropdown.Item>
                <Dropdown.Item onClick={handleViewDetailsClick}>
                    <i className="mdi mdi-information-outline me-1"></i>View details
                </Dropdown.Item>
            </Dropdown.Menu>
        </Dropdown>
    )
}

IdentityCertificateDropDown.propTypes = {
    handleCopySignatureClick: PropTypes.func.isRequired,
    handleCopySerialNumberClick: PropTypes.func.isRequired,
}

function CertificateDetails({ data }) {
    const getExpireColorStyle = (notAfter) => {
        const _notAfter = dayjs(notAfter).unix()
        const threeMonthsFromNow = dayjs().add(3, 'month').unix()

        if (dayjs().unix() > _notAfter) {
            return 'text-danger'
        }

        if (_notAfter < threeMonthsFromNow) {
            return 'text-warning'
        }

        return 'text-success'
    }

    const getBadgeClassName = () => {
        const threeMonthsFromNow = dayjs().add(3, 'month').unix()
        const _notAfter = dayjs(data?.not_after).unix()

        if (data?.revoked || (dayjs().unix() > _notAfter)) {
            return 'bg-danger'
        }

        if (_notAfter < threeMonthsFromNow) {
            return 'bg-warning'
        }

        return 'bg-primary'
    }

    const getBadgeLabel = () => {
        const threeMonthsFromNow = dayjs().add(3, 'month').unix()
        const _notAfter = dayjs(data?.not_after).unix()

        if (data?.revoked) {
            return 'revoked'
        }

        if (dayjs().unix() > _notAfter) {
            return 'expired'
        }

        if (_notAfter < threeMonthsFromNow) {
            return 'expiring soon'
        }

        return 'valid'
    }


    const handleCopySignatureClick = async (signature) => {
        await copyToClipboard(signature)
    }

    const handleCopySerialNumberClick = async (serial) => {
        await copyToClipboard(serial)
    }

    return (
        <Card>
            {
                !data ? <Card.Body className='fst-italic text-muted'>Identity certificate not available</Card.Body> : (
                    <Card.Body>
                        <IdentityCertificateDropDown
                            handleCopySerialNumberClick={() => handleCopySerialNumberClick(data?.serial_number)}
                            handleCopySignatureClick={() => handleCopySignatureClick(data?.signature)}
                        />
                        <h4 className="m-0 text-black">Identity Certificate</h4>
                        <span data-testid="certificate-state" className={`badge rounded-pill px-1 ${getBadgeClassName()}`}>{getBadgeLabel()}</span>

                        <p className="fw-bold mb-1 mt-3">Serial number: <span className="fw-normal">{formatDisplayedData(data?.serial_number)}</span></p>
                        <p data-testid="certificate-expiration-date" className={`mb-1 ${getExpireColorStyle(data?.not_after)}`}><span className='fw-bold'>Expires:</span>  {new Date(data?.not_after).toUTCString()}</p>
                        <p className='mb-1'><span className='fw-bold'>Issuer: </span>{data?.issuer?.common_name}</p>
                        <p className='mb-1'><span className='fw-bold'>Subject: </span>{data?.subject?.common_name}</p>

                        <Row>
                            <Col>
                                <FileInformationCard file={data?.data} name="Public identity key" ext=".PEM" />
                                <FileInformationCard file={data?.chain} name="TRISA trust chain (CA)" ext=".GZ" />
                            </Col>
                        </Row>
                    </Card.Body>
                )
            }
        </Card>
    )
}

export default React.memo(CertificateDetails)
