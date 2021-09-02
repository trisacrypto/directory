
import React from 'react'
import { Card, Col, Row } from 'react-bootstrap';
import { formatDisplayedData } from "../../../utils"



function TrixoForm({ data }) {
    return (
        <Card>
            <Card.Body>
                <h4 className="mt-0 mb-3">TRIXO Form</h4>
                <Row>
                    <Col>
                        <p className="fw-bold mb-2">Applicable Regulation : <span className="fw-normal">{formatDisplayedData(data?.applicable_regulations)}</span></p>
                        <p className="fw-bold mb-2">Compliance Threshold : <span className="fw-normal">{formatDisplayedData(data?.compliance_threshold)}</span></p>
                        <p className="fw-bold mb-2">Compliance Threshold Currency : <span className="fw-normal">{formatDisplayedData(data?.compliance_threshold_currency)}</span></p>
                        <p className="fw-bold mb-2">Conducts Customer KYC : <span className="fw-normal">{formatDisplayedData(data?.conducts_customer_kyc)}</span></p>
                        <p className="fw-bold mb-2">Financial Transfers Permitted : <span className="fw-normal">{formatDisplayedData(data?.financial_transfers_permitted)}</span></p>
                    </Col>
                    <Col>
                        <p className="fw-bold mb-2">Has Required Regulatory Program : <span className="fw-normal">{formatDisplayedData(data?.has_required_regulatory_program)}</span></p>
                        <p className="fw-bold mb-2">KYC Threshold : <span className="fw-normal">{formatDisplayedData(data?.kyc_threshold)}</span></p>
                        <p className="fw-bold mb-2">KYC Threshold Currency : <span className="fw-normal">{formatDisplayedData(data?.kyc_threshold_currency)}</span></p>
                        <p className="fw-bold mb-2">Must Comply Travel Rule : <span className="fw-normal">{formatDisplayedData(data?.must_comply_travel_rule)}</span></p>
                        <p className="fw-bold mb-2">Must Save Guard PII : <span className="fw-normal">{formatDisplayedData(data?.must_safeguard_pii)}</span></p>
                    </Col>
                    <Col xl={4}>
                        <p className="fw-bold mb-2">Save Guard PII : <span className="fw-normal">{formatDisplayedData(data?.safeguards_pii)}</span></p>
                        <p className="fw-bold mb-2">Other Juridictions : <span className="fw-normal">{formatDisplayedData(data?.other_jurisdictions)}</span></p>
                        <p className="fw-bold mb-2">Primary National Juridictions : <span className="fw-normal">{formatDisplayedData(data?.primary_national_jurisdiction)}</span></p>
                        <p className="fw-bold mb-2">Primary Regulator : <span className="fw-normal">{formatDisplayedData(data?.primary_regulator)}</span></p>
                    </Col>

                </Row>
            </Card.Body>
        </Card>
    )
}

export default TrixoForm
