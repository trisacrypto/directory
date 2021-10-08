import React from 'react';
import update from 'immutability-helper';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import LegalPersonNameTypeCode from '../select/LegalPersonNameTypeCode';
import { Trans } from "@lingui/macro"


const LegalPersonName = ({name, onChange}) => {
  const createArrayChangeHandler = (field, idx, key) => (event) => {
    const changes = {[field]: {[idx]: {[key]: {$set: event.target.value}}}};
    const updated = update(name, changes);
    onChange(null, updated);
  }

  const createArrayRemoveHandler = (field, idx) => (event) => {
    const changes = {[field]: {$splice: [[idx, 1]]}};
    const updated = update(name, changes);
    onChange(null, updated);
  }

  const createArrayPushHandler = (field) => (event) => {
    const changes = {[field]: {$push: [{"legal_person_name": "", "legal_person_name_identifier_type": 0}]}};
    const updated = update(name, changes);
    onChange(null, updated);
  }

  const nameIdentifiers = name.name_identifiers.map((name, idx) => {
    return (
      <Form.Row key={idx}>
        <Form.Group as={Col}>
          <Form.Control
            type="text"
            value={name.legal_person_name}
            onChange={createArrayChangeHandler("name_identifiers", idx, "legal_person_name")}
          />
        </Form.Group>
        <Form.Group as={Col}>
          <Form.Control
            as="select" custom
            value={name.legal_person_name_identifier_type}
            onChange={createArrayChangeHandler("name_identifiers", idx, "legal_person_name_identifier_type")}
          >
            <LegalPersonNameTypeCode />
          </Form.Control>
        </Form.Group>
        <Form.Group as={Col} xs={1}>
          <Button
            variant="danger"
            onClick={createArrayRemoveHandler("name_identifiers", idx)}
          >
            <i className="fa fa-trash"></i>
          </Button>
        </Form.Group>
      </Form.Row>
    );
  });

  const localNameIdentifiers = name.local_name_identifiers.map((name, idx) => {
    return (
      <Form.Row key={idx}>
        <Form.Group as={Col}>
          <Form.Control
            type="text"
            value={name.legal_person_name}
            onChange={createArrayChangeHandler("local_name_identifiers", idx, "legal_person_name")}
          />
        </Form.Group>
        <Form.Group as={Col}>
          <Form.Control
            as="select" custom
            value={name.legal_person_name_identifier_type}
            onChange={createArrayChangeHandler("local_name_identifiers", idx, "legal_person_name_identifier_type")}
          >
            <LegalPersonNameTypeCode />
          </Form.Control>
        </Form.Group>
        <Form.Group as={Col} xs={1}>
          <Button
            variant="danger"
            onClick={createArrayRemoveHandler("local_name_identifiers", idx)}
          >
            <i className="fa fa-trash"></i>
          </Button>
        </Form.Group>
      </Form.Row>
    );
  });

  const phoneticNameIdentifiers = name.phonetic_name_identifiers.map((name, idx) => {
    return (
      <Form.Row key={idx}>
        <Form.Group as={Col}>
          <Form.Control
            type="text"
            value={name.legal_person_name}
            onChange={createArrayChangeHandler("phonetic_name_identifiers", idx, "legal_person_name")}
          />
        </Form.Group>
        <Form.Group as={Col}>
          <Form.Control
            as="select" custom
            value={name.legal_person_name_identifier_type}
            onChange={createArrayChangeHandler("phonetic_name_identifiers", idx, "legal_person_name_identifier_type")}
          >
            <LegalPersonNameTypeCode />
          </Form.Control>
        </Form.Group>
        <Form.Group as={Col} xs={1}>
          <Button
            variant="danger"
            onClick={createArrayRemoveHandler("phonetic_name_identifiers", idx)}
          >
            <i className="fa fa-trash"></i>
          </Button>
        </Form.Group>
      </Form.Row>
    );
  });

  const nameLabel = () => {
    if (nameIdentifiers.length > 0) {
      return (
        <>
        <Form.Label className="mb-0 pb-0"><Trans>Name Identifiers</Trans></Form.Label>
        <p className="text-muted mt-0 pt-0">
          <small><Trans>The name and type of name by which the legal person is known.</Trans></small>
        </p>
        </>
      );
    }
  }

  const localNameLabel = () => {
    if (localNameIdentifiers.length > 0) {
      return (
        <>
        <Form.Label className="mb-0 pb-0"><Trans>Local Name Identifiers</Trans></Form.Label>
        <p className="text-muted mt-0 pt-0">
          <small><Trans>The name by which the legal person is known using local characters.</Trans></small>
        </p>
        </>
      );
    }
  }

  const phoneticNameLabel = () => {
    if (phoneticNameIdentifiers.length > 0) {
      return (
        <>
        <Form.Label className="mb-0 pb-0"><Trans>Phonetic Name Identifiers</Trans></Form.Label>
        <p className="text-muted mt-0 pt-0">
          <small><Trans>A phonetic representation of the name by which the legal person is known.</Trans></small>
        </p>
        </>
      );
    }
  }

  return (
    <>
    {nameLabel()}
    {nameIdentifiers}

    {localNameLabel()}
    {localNameIdentifiers}

    {phoneticNameLabel()}
    {phoneticNameIdentifiers}

    <Form.Group>
      <Button size="sm" onClick={createArrayPushHandler("name_identifiers")}><Trans>Add Legal Name</Trans></Button>{' '}
      <Button size="sm" onClick={createArrayPushHandler("local_name_identifiers")}><Trans>Add Local Name</Trans></Button>{' '}
      <Button size="sm" onClick={createArrayPushHandler("phonetic_name_identifiers")}><Trans>Add Phonetic Names</Trans></Button>
    </Form.Group>
    </>
  );
}

export default LegalPersonName;