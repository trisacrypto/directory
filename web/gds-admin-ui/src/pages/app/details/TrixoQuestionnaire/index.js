import React from 'react'
import { financialTransfersPermitted, hasRequiredRegulatoryProgram } from 'constants/trixo';
import { Card, Col, Row } from 'react-bootstrap';
import { getConductsCustomerKYC, getMustComplyRegulations, getMustSafeguardPii, getSafeguardPii, intlFormatter, validateIsoCode } from "utils"
import countryCodeEmoji, { isoCountries } from 'utils/country';
import TrixoDropdown from './TrixoDropdown';

function TrixoForm({ data }) {
    // console.log(data);
    return (
        <Card>
            <Card.Body>
                <TrixoDropdown data={data} />
                <h4 className="mt-0 text-black">TRIXO Questionnaire</h4>
                <p className='lh-lg' data-testid="financial-transfers-permitted">Organization <span className='fw-bold'>{financialTransfersPermitted[data?.financial_transfers_permitted]}</span> permitted to send and/or receive transfers of virtual assets in the jurisdiction(s) in which it operates</p>
                <Row>
                    <Col xs={12} md={6}>

                        <h5 className='text-black'>Jurisdictions</h5>
                        <hr className='mt-1' />
                        <ul className='list-unstyled d-flex'>
                            <li>
                                <span className='badge bg-primary rounded-pill px-1 rounded-pill me-1'>Primary</span>
                            </li>
                            <ul className='list-unstyled'>
                                <li data-testid="primary-jurisdiction">
                                    {countryCodeEmoji(validateIsoCode(data?.primary_national_jurisdiction))} {isoCountries[data?.primary_national_jurisdiction]} regulated by {data?.primary_regulator}
                                </li>
                                {
                                    Array.isArray(data?.other_jurisdictions) && data?.other_jurisdictions?.map((jurisdiction, index) => {
                                        return (
                                            <li data-testid='other-jurisdiction' key={index}>{countryCodeEmoji(validateIsoCode(jurisdiction?.country))} {isoCountries[jurisdiction?.country]} regulated by {jurisdiction?.regulator_name}</li>
                                        )
                                    })
                                }
                            </ul>
                        </ul>
                        <p data-testid="must-comply-travel-rule" className='lh-lg'>Organization <span className='fw-bold'>{getMustComplyRegulations(data?.must_comply_travel_rule)}</span> comply with the application of the Travel Rule standards in the jurisdiction(s) where it is licensed/approved/registered.</p>
                    </Col>
                    <Col xs={12} md={6}>
                        <h5 className='text-black'>Applicable Regulations</h5>
                        <hr className='mt-1' />
                        <ul>
                            {
                                data?.applicable_regulations?.map(regulation => <li key={regulation}>{regulation}</li>)
                            }
                        </ul>
                        <p><span className='fw-bold'>Minimum Compliance Threshold:</span> <span data-testid="compliance_threshold_currency">{intlFormatter({ currency: data?.compliance_threshold_currency }).format(data?.compliance_threshold)}</span> {data?.compliance_threshold_currency}</p>
                    </Col>
                    <Col xs={12} md={6}>
                        <h5 className='text-black'>CDD & Travel Rule Policies</h5>
                        <hr className='mt-1' />
                        <p className='lh-lg' data-testid="has-required-regulatory-program">
                            Organization <span className='fw-bold'>{hasRequiredRegulatoryProgram[data?.has_required_regulatory_program]}</span> program that sets minimum AML, CFT, KYC/CDD and Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes where it is licensed/approved/registered.
                        </p>
                        <p className='lh-lg' data-testid="conducts-customer-kyc">
                            Organization <span className='fw-bold'>{getConductsCustomerKYC(data?.conducts_customer_kyc)}</span> conduct KYC/CDD before permitting its customers to send/receive virtual asset transfers.
                        </p>
                        <p><span className='fw-bold'>Conducts KYC Threshold:</span> {intlFormatter({ currency: data?.kyc_threshold_currency }).format(data?.kyc_threshold)} {data?.kyc_threshold_currency}</p>
                    </Col>
                    <Col xs={12} md={6}>
                        <h5 className='text-black'>Data Protection Policies</h5>
                        <hr className='mt-1' />

                        <p data-testid="must-safeguard-pii" className='lh-lg'>Organization <span className='fw-bold'>{getMustSafeguardPii(data?.must_safeguard_pii)}</span> safeguard PII by law.</p>
                        <p data-testid="safeguards-pii" className='lh-lg'>Organization <span className='fw-bold'>{getSafeguardPii(data?.safeguards_pii)}</span> secure and protect PII, including PII received from other VASPs under the Travel Rule.</p>
                    </Col>
                </Row>
            </Card.Body>
        </Card>
    )
}

export default React.memo(TrixoForm)