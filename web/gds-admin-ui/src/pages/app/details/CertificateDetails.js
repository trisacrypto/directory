
import FileInformationCard from 'components/FileInformationCard';
import React from 'react'
import { Card, Col, Dropdown, Row } from 'react-bootstrap';
import { formatDisplayedData } from 'utils';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime'
dayjs.extend(relativeTime)


export const IdentityCertificateDropDown = () => {

    return (
        <Dropdown className="float-end" align="end">
            <Dropdown.Toggle
                data-testid="dripicons-dots-3"
                variant="link"
                tag="a"
                className="card-drop arrow-none cursor-pointer p-0 shadow-none">
                <i className="dripicons-dots-3"></i>
            </Dropdown.Toggle>
            <Dropdown.Menu>
                <Dropdown.Item>
                    <i className="mdi mdi-card-search me-1"></i>Copy signature
                </Dropdown.Item>
                <Dropdown.Item>
                    <i className="mdi mdi-card-search me-1"></i>Copy serial number
                </Dropdown.Item>
                <Dropdown.Item>
                    <i className="mdi mdi-card-search me-1"></i>View details
                </Dropdown.Item>
            </Dropdown.Menu>
        </Dropdown>
    )
}

function CertificateDetails({ data }) {
    const getExpireColorStyle = (notAfter) => {
        const _notAfter = dayjs(notAfter).unix()
        const threeMonthsFromNow = dayjs().add(3, 'month').unix()

        // console.log('[+3month]', dayjs().add(3, 'month').format())
        // console.log('[AFTER]', dayjs(notAfter).format())
        // console.log('[NOW]', dayjs().format())
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

        if (_notAfter < threeMonthsFromNow) {
            return 'expiring soon'
        }

        if (dayjs().unix() > _notAfter) {
            return 'expired'
        }

        return 'valid'
    }

    return (
        <Card>
            {
                !data ? <Card.Body className='fst-italic text-muted'>Identity certificate not available</Card.Body> : (
                    <Card.Body>
                        <IdentityCertificateDropDown />
                        <h4 className="m-0 text-black">Identity Certificate</h4>
                        <span className={`badge rounded-pill px-1 ${getBadgeClassName()}`}>{getBadgeLabel()}</span>

                        <p className="fw-bold mb-1 mt-3">Serial number: <span className="fw-normal">{formatDisplayedData(data?.serial_number)}</span></p>
                        <p className={`mb-1 ${getExpireColorStyle(data?.not_after)}`}><span className='fw-bold'>Expires:</span>  {new Date(data?.not_after).toUTCString()}</p>
                        <p className='mb-1'><span className='fw-bold'>Issuer: </span>{data?.issuer?.common_name}</p>
                        <p className='mb-1'><span className='fw-bold'>Subject: </span>{data?.subject?.common_name}</p>

                        <Row>
                            <Col>
                                <FileInformationCard name="Public identity key" ext="Base64" />
                                <FileInformationCard file={data?.chain} name="TRISA trust chain (CA)" />
                            </Col>
                        </Row>
                    </Card.Body>
                )
            }
        </Card>
    )
}

export default React.memo(CertificateDetails)
