import React from 'react'
import FileInformationCard from 'components/FileInformationCard';
import { Card, Col, Row } from 'react-bootstrap';
import { copyToClipboard, formatDisplayedData } from 'utils';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime'
import { downloadFile } from 'helpers/api/utils';
import Details from './Details';
import IdentityCertificateDropDown from './IdentityCertificateDropDown'
dayjs.extend(relativeTime)

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

function CertificateDetails({ data }) {


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

    const handlePublicIdentityKeyDownloadClick = (data) => {
        const filename = 'public-identity-key.pem'
        const mimetype = 'application/x-pem-file'
        if (data) {
            downloadFile(data, filename, mimetype)
        }
    }

    const handleTrustChainDownloadClick = (chain) => {
        const filename = 'trust-chain-certificate.gz'
        const mimetype = 'application/x-x509-ca-cert'
        if (chain) {
            downloadFile(chain, filename, mimetype)
        }
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
                            <Col sm={12}>
                                <FileInformationCard file={data?.data} name="Public Identity Key" ext=".PEM" onDownload={() => handlePublicIdentityKeyDownloadClick(data?.data)} />
                                <FileInformationCard file={data?.chain} name="TRISA Trust Chain (CA)" ext=".GZ" onDownload={() => handleTrustChainDownloadClick(data?.chain)} />
                            </Col>
                            <Col>
                                <Details data={data} />
                            </Col>
                        </Row>
                    </Card.Body>
                )
            }
        </Card>
    )
}

export default React.memo(CertificateDetails)
