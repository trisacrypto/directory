import React from 'react';
import Form from 'react-bootstrap/Form';
import { Trans } from "@lingui/macro";


const Contact = ({contact, onChange, required}) => {
  const createChangeHandler = (field) => (event) => {
    let data = {...contact};
    data[[field]] = event.target.value;
    onChange(null, data);
  }

  return (
    <>
    <Form.Group controlId="contactName">
      <Form.Label><Trans>Full Name</Trans></Form.Label>
      <Form.Control
        type="text"
        value={contact.name}
        onChange={createChangeHandler("name")}
        required={required}
      />
      <Form.Text className="text-muted">
        <Trans>Preferred name for email communication.</Trans>
      </Form.Text>
    </Form.Group>
    <Form.Group controlId="contactEmail">
      <Form.Label><Trans>Email address</Trans></Form.Label>
      <Form.Control
        type="email"
        value={contact.email}
        onChange={createChangeHandler("email")}
        required={required}
      />
      <Form.Text className="text-muted">
        <Trans>Please use the email address associated with your organization.</Trans>
      </Form.Text>
      <Form.Control.Feedback type="invalid">
        <Trans>Please supply a valid email address.</Trans>
      </Form.Control.Feedback>
    </Form.Group>
    <Form.Group controlId="contactPhone">
      <Form.Label><Trans>Phone Number</Trans></Form.Label>
      <Form.Control
        type="tel"
        value={contact.phone}
        onChange={createChangeHandler("phone")}
        required={required}
      />
      <Form.Text className="text-muted">
        {required
          ? <Trans>Required - please supply full phone number with country code.</Trans>
          : <Trans>Optional - if supplied, use full phone number with country code.</Trans>
        }
      </Form.Text>
      <Form.Control.Feedback type="invalid">
        <Trans>Please supply a valid phone number or omit entirely if not required.</Trans>
      </Form.Control.Feedback>
    </Form.Group>

    </>
  );
}

export default Contact;