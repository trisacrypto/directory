import React from 'react';
import Form from 'react-bootstrap/Form';

const Contact = ({contact, onChange}) => {
  const createChangeHandler = (field) => (event) => {
    let data = {...contact};
    data[[field]] = event.target.value;
    onChange(null, data);
  }

  return (
    <>
    <Form.Group controlId="contactName">
      <Form.Label>Full Name</Form.Label>
      <Form.Control
        type="text"
        value={contact.name}
        onChange={createChangeHandler("name")}
      />
      <Form.Text className="text-muted">
        Preferred name for email communication.
      </Form.Text>
    </Form.Group>
    <Form.Group controlId="contactEmail">
      <Form.Label>Email address</Form.Label>
      <Form.Control
        type="email"
        value={contact.email}
        onChange={createChangeHandler("email")}
      />
      <Form.Text className="text-muted">
        Please use the email address associated with your organization.
      </Form.Text>
      <Form.Control.Feedback type="invalid">
        Please supply a valid email address.
      </Form.Control.Feedback>
    </Form.Group>
    <Form.Group controlId="contactPhone">
      <Form.Label>Phone Number</Form.Label>
      <Form.Control
        type="tel"
        value={contact.phone}
        onChange={createChangeHandler("phone")}
      />
      <Form.Text className="text-muted">
        Optional - if supplied, use full phone number with country code.
      </Form.Text>
      <Form.Control.Feedback type="invalid">
        Please supply a valid phone number or omit entirely.
      </Form.Control.Feedback>
    </Form.Group>

    </>
  );
}

export default Contact;