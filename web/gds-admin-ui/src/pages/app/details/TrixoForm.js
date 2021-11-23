
import { financialTransfersPermitted, hasRequiredRegulatoryProgram } from 'constants/trixo';
import React from 'react'
import { Card, Col, Row } from 'react-bootstrap';
import { intlFormatter } from "utils"
import countryCodeEmoji, { isoCountries } from 'utils/country';



function TrixoForm({ data }) {
    const getMustComplyRegulations = (status) => status ? "must" : "must not"
    const getConductsCustomerKYC = (status) => status ? "does" : "does not"
    const getMustSafeguardPii = (status) => status ? "must" : "is not required to"
    const getSafeguardPii = (status) => status ? "does" : "does not"

    return (
        <Card>
            <Card.Body>
                <h4 className="mt-0 text-black">TRIXO Questionnaire</h4>
                <p className='lh-lg'>Organization <span className='fw-bold'>{financialTransfersPermitted[data?.financial_transfers_permitted]}</span> permitted to send and/or receive transfers of virtual assets in the jurisdictions in which it operates</p>
                <Row>
                    <Col xs={12} md={6}>
                        <h5 className='text-black'>Juridictions</h5>
                        <hr className='mt-1' />
                        <p><span className='badge bg-primary rounded-pill px-1 rounded-pill'>Primary</span> {countryCodeEmoji(data?.primary_national_jurisdiction)} {isoCountries[data?.primary_national_jurisdiction]} regulated by {data?.primary_regulator}</p>
                        <p>
                            {
                                Array.isArray(data?.other_jurisdictions) && data?.other_jurisdictions.map(juridiction => {
                                    return (
                                        <p>{isoCountries[juridiction?.country]} regulated by {juridiction?.regulator_name}</p>
                                    )
                                })
                            }
                        </p>
                        <p className='lh-lg'>Organization <span className='fw-bold'>{getMustComplyRegulations(data?.must_comply_regulations)}</span> comply with the application of the Travel Rule standards in the jurisdiction(s) where it is licensed/approved/registered.</p>
                    </Col>
                    <Col xs={12} md={6}>
                        <h5 className='text-black'>Applicable Regulations</h5>
                        <hr className='mt-1' />
                        <ul>
                            {
                                data?.applicable_regulations.map(regulation => <li key={regulation}>{regulation}</li>)
                            }
                        </ul>
                        <p><span className='fw-bold'>Minimum Compliance Threshold:</span> {intlFormatter({ currency: data?.compliance_threshold_currency }).format(data?.compliance_threshold)} {data?.compliance_threshold_currency}</p>
                    </Col>
                    <Col xs={12} md={6}>
                        <h5 className='text-black'>CDD & Travel Rule Policies</h5>
                        <hr className='mt-1' />
                        <p className='lh-lg'>
                            Organization <span className='fw-bold'>{hasRequiredRegulatoryProgram[data?.has_required_regulatory_program]}</span> program that sets minimum AML, CFT, KYC/CDD and Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes where it is licensed/approved/registered.
                        </p>
                        <p className='lh-lg'>
                            Organization <span className='fw-bold'>{getConductsCustomerKYC(data?.conducts_customer_kyc)}</span> conduct KYC/CDD before permitting its customers to send/receive virtual asset transfers.
                        </p>
                        <p><span className='fw-bold'>Conducts KYC Threshold:</span> {intlFormatter({ currency: data?.kyc_threshold_currency }).format(data?.kyc_threshold)} {data?.kyc_threshold_currency}</p>
                    </Col>
                    <Col xs={12} md={6}>
                        <h5 className='text-black'>Data Protection Policies</h5>
                        <hr className='mt-1' />
                        <p className='lh-lg'>Organization <span className='fw-bold'>{getMustSafeguardPii(data?.must_safeguard_pii)} to</span> safeguard PII by law.</p>
                        <p className='lh-lg'>Organization <span className='fw-bold'>{getSafeguardPii(data?.safeguard_pii)}</span> secure and protect PII, including PII received from other VASPs under the Travel Rule.</p>
                    </Col>
                </Row>
            </Card.Body>
        </Card>
    )
}

export default React.memo(TrixoForm)
