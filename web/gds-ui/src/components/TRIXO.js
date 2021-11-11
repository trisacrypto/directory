import React from 'react';
import Currencies from './select/Currencies';
import Countries from './select/Countries';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import update from 'immutability-helper';
import { Trans } from "@lingui/macro";
import { t } from "@lingui/macro";


const TRIXO = ({data, onChange}) => {

  const createChangeHandler = (field) => (event) => {
    const changes = {[field]: {$set: event.target.value}};
    const updated = update(data, changes);
    onChange(null, updated);
  }

  const createBoolChangeHandler = (field) => (event) => {
    const changes = {[field]: {$set: event.target.checked}};
    const updated = update(data, changes);
    onChange(null, updated);
  }

  const createArrayChangeObjectHandler = (field, idx, key) => (event) => {
    const changes = {[field]: {[idx]: {[key]: {$set: event.target.value}}}};
    const updated = update(data, changes);
    onChange(null, updated);
  }

  const createArrayChangeHandler = (field, idx) => (event) => {
    const changes = {[field]: {[idx]: {$set: event.target.value}}};
    const updated = update(data, changes);
    onChange(null, updated);
  }

  const createArrayRemoveHandler = (field, idx) => (event) => {
    const changes = {[field]: {$splice: [[idx, 1]]}};
    const updated = update(data, changes);
    onChange(null, updated);
  }

  const createArrayPushHandler = (field, value) => (event) => {
    const changes = {[field]: {$push: [value]}};
    const updated = update(data, changes);
    onChange(null, updated);
  }

  const otherJursidictions = data.other_jurisdictions.map((item, idx) => {
    return (
      <Form.Row key={idx}>
        <Form.Group as={Col}>
          <Form.Label><Trans>National Jurisdiction</Trans></Form.Label>
          <Form.Control
            as="select" custom
            value={item.country}
            onChange={createArrayChangeObjectHandler("other_jurisdictions", idx, "country")}
          >
            <Countries />
          </Form.Control>
        </Form.Group>
        <Form.Group as={Col}>
          <Form.Label><Trans>Regulator Name</Trans></Form.Label>
          <Form.Control
            value={item.regulator_name}
            onChange={createArrayChangeObjectHandler("other_jurisdictions", idx, "regulator_name")}
          />
        </Form.Group>
        <Form.Group as={Col} xs={1}>
          <Form.Label>&nbsp;</Form.Label>
          <Button
            className="form-control"
            variant="danger"
            onClick={createArrayRemoveHandler("other_jurisdictions", idx)}
          >
            <i className="fa fa-trash"></i>
          </Button>
        </Form.Group>
      </Form.Row>
    )
  })

  const applicableRegulations = data.applicable_regulations.map((item, idx) => {
    return (
      <Form.Row key={idx}>
        <Form.Group as={Col}>
          <Form.Control
            value={item}
            onChange={createArrayChangeHandler("applicable_regulations", idx)}
          />
        </Form.Group>
        <Form.Group as={Col} xs={1}>
          <Button
            className="form-control"
            variant="danger"
            onClick={createArrayRemoveHandler("applicable_regulations", idx)}
          >
            <i className="fa fa-trash"></i>
          </Button>
        </Form.Group>
      </Form.Row>
    )
  })

  return (
    <>
      <Form.Group controlId="trixoPrimaryNationalJurisdiction">
        <Form.Label><Trans>Primary National Jurisdiction</Trans></Form.Label>
        <Form.Control
          as="select" custom
          value={data.primary_national_jurisdiction}
          onChange={createChangeHandler("primary_national_jurisdiction")}
        >
          <Countries />
        </Form.Control>
      </Form.Group>
      <Form.Group controlId="trixoPrimaryNationalJurisdiction">
        <Form.Label><Trans>Name of Primary Regulator</Trans></Form.Label>
        <Form.Control
          type="text"
          value={data.primary_regulator}
          onChange={createChangeHandler("primary_regulator")}
        />
        <Form.Text className="text-muted">
          <Trans>The name of primary regulator or supervisory authority for your national jurisdiction</Trans>
        </Form.Text>
      </Form.Group>
      <fieldset>
        <legend className="sublegend"><Trans>Other Jursidictions</Trans></legend>
        <p><Trans>Please add any other regulatory jurisdictions your organization complies with.</Trans></p>
        {otherJursidictions}
        <Form.Group>
          <Button size="sm" onClick={createArrayPushHandler('other_jurisdictions', {'country': '', 'regulator_name': ''})}><Trans>Add Jurisdiction</Trans></Button>
        </Form.Group>
      </fieldset>
      <Form.Group>
        <Form.Label><Trans>Is your organization permitted to send and/or receive transfers of virtual assets in the jurisdictions in which it operates?</Trans></Form.Label>
        <Form.Control
          as="select" custom
          value={data.financial_transfers_permitted}
          onChange={createChangeHandler("financial_transfers_permitted")}
        >
          <option value=""></option>
          <option value="yes">{t`Yes`}</option>
          <option value="partial">{t`Partially`}</option>
          <option value="no">{t`No`}</option>
        </Form.Control>
      </Form.Group>
      <fieldset>
        <legend className="sublegend"><Trans>CDD & Travel Rule Policies</Trans></legend>
        <Form.Group>
          <Form.Label><Trans>Does your organization have a programme that sets minimum AML, CFT, KYC/CDD and Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes where it is licensed/approved/registered?</Trans></Form.Label>
          <Form.Control
            as="select" custom
            value={data.has_required_regulatory_program}
            onChange={createChangeHandler("has_required_regulatory_program")}
          >
            <option value=""></option>
            <option value="yes">{t`Yes`}</option>
            <option value="partial">{t`Partially Implemented`}</option>
            <option value="no">{t`No`}</option>
          </Form.Control>
        </Form.Group>
        <Form.Group>
          <Form.Label><Trans>Does your organization conduct KYC/CDD before permitting its customers to send/receive virtual asset transfers?</Trans></Form.Label>
          <Form.Check
            type="switch"
            id="conductsCustomerKYC"
            label={t`Conducts KYC before virtual asset transfers`}
            checked={data.conducts_customer_kyc}
            onChange={createBoolChangeHandler("conducts_customer_kyc")}
          />
        </Form.Group>
        <Form.Group>
          <Form.Label><Trans>At what threshold and currency does your organization conduct KYC?</Trans></Form.Label>
          <Form.Row>
            <Col>
              <Form.Control
                type="number"
                value={data.kyc_threshold}
                onChange={createChangeHandler("kyc_threshold")}
              />
            </Col>
            <Col xs={3}>
              <Form.Control
                as="select" custom
                value={data.kyc_threshold_currency}
                onChange={createChangeHandler("kyc_threshold_currency")}
              >
                <Currencies />
              </Form.Control>
            </Col>
          </Form.Row>
          <Form.Text className="text-muted">
            <Trans>Threshold to conduct KYC before permitting the customer to send/receive virtual asset transfers</Trans>
          </Form.Text>
        </Form.Group>
        <Form.Group>
          <Form.Label><Trans>Is your organization required to comply with the application of the Travel Rule standards in the jurisdiction(s) where it is licensed/approved/registered?</Trans></Form.Label>
          <Form.Check
            type="switch"
            id="mustComplyTravelRule"
            label={t`Must comply with the Travel Rule`}
            checked={data.must_comply_travel_rule}
            onChange={createBoolChangeHandler("must_comply_travel_rule")}
          />
        </Form.Group>
        <fieldset>
          <Form.Label className="mb-0"><Trans>Applicable Regulations</Trans></Form.Label>
          <p className="text-muted mt-0">
            <small><Trans>Please specify the applicable regulation(s) for Travel Rule standards compliance, e.g. "FATF Recommendation 16"</Trans></small>
          </p>
          {applicableRegulations}
          <Form.Group>
            <Button size="sm" onClick={createArrayPushHandler("applicable_regulations", "")}><Trans>Add Regulation</Trans></Button>
          </Form.Group>
        </fieldset>
        <Form.Group>
          <Form.Label><Trans>What is the minimum threshold for Travel Rule compliance?</Trans></Form.Label>
          <Form.Row>
            <Col>
              <Form.Control
                type="number"
                value={data.compliance_threshold}
                onChange={createChangeHandler("compliance_threshold")}
              />
            </Col>
            <Col xs={3}>
              <Form.Control
                as="select" custom
                value={data.compliance_threshold_currency}
                onChange={createChangeHandler("compliance_threshold_currency")}
              >
                <Currencies />
              </Form.Control>
            </Col>
          </Form.Row>
          <Form.Text className="text-muted">
            <Trans>The minimum threshold above which your organization is required to collect/send Travel Rule information.</Trans>
          </Form.Text>
        </Form.Group>
      </fieldset>
      <fieldset>
        <legend className="sublegend"><Trans>Data Protection Policies</Trans></legend>
        <Form.Group>
          <Form.Label><Trans>Is your organization required by law to safeguard PII?</Trans></Form.Label>
          <Form.Check
            type="switch"
            id="mustSafeguardPII"
            label={t`Must Safeguard PII`}
            checked={data.must_safeguard_pii}
            onChange={createBoolChangeHandler("must_safeguard_pii")}
          />
        </Form.Group>
        <Form.Group>
          <Form.Label><Trans>Does your organization secure and protect PII, including PII received from other VASPs under the Travel Rule?</Trans></Form.Label>
          <Form.Check
            type="switch"
            id="safeguardsPII"
            label={t`Safeguards PII`}
            checked={data.safeguards_pii}
            onChange={createBoolChangeHandler("safeguards_pii")}
          />
        </Form.Group>
      </fieldset>
    </>
  )

}

export default TRIXO;