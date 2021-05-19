import React from 'react';
import Contact from './Contact';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Tab from 'react-bootstrap/Tab';
import Nav from 'react-bootstrap/Nav';
import Form from 'react-bootstrap/Form';
import Card from 'react-bootstrap/Card';
import Button from 'react-bootstrap/Button';
import Accordion from 'react-bootstrap/Accordion';
import update from 'immutability-helper';


class Registration extends React.Component {
  state = {
    tabKey: "introduction",
    validated: false,
    formData: {
      contacts: {
        technical: {name: "", email: "", phone: ""},
        legal: {name: "", email: "", phone: ""},
        administrative: {name: "", email: "", phone: ""},
        billing: {name: "", email: "", phone: ""},
      }
    },
  }

  handleSubmit = (event) => {
    event.preventDefault();

    const form = event.currentTarget;
    if (form.checkValidity() === false) {
      event.stopPropagation();
    }

    this.setState({validated: true});
    console.log(this.state.formData);
  }

  createChangeHandler = (field, ...parents) => (event, value, selectedKey) => {
    // Identify the particular change that we're going to be making
    const reducer = (acc, item) => ({[item]: acc});
    const changes = parents.length === 0 ? {[field]: {$set: value}} : parents.reduce(reducer, {[field]: {$set: value}});

    // Create an immutable copy of the data first to avoid changes by reference.
    const formData = update(this.state.formData, changes);
    this.setState({formData: formData})
  }

  render() {
    return (
      <Form noValidate validated={this.state.validated} onSubmit={this.handleSubmit}>
        <Tab.Container id="registration-form" activeKey={this.state.tabKey} onSelect={(k) => this.setState({tabKey: k})}>
          <Row className="pt-3 pb-5">
            <Col sm={3}>
              <Nav variant="pills" className="flex-column">
                <Nav.Item>
                  <Nav.Link eventKey="introduction">Introduction</Nav.Link>
                  <Nav.Link eventKey="entity-details">Entity Details</Nav.Link>
                  <Nav.Link eventKey="contacts">Contacts</Nav.Link>
                  <Nav.Link eventKey="summary">Summary</Nav.Link>
                </Nav.Item>
              </Nav>
            </Col>
            <Col sm={9}>
              <Tab.Content>
                <Tab.Pane eventKey="introduction">
                  <fieldset>
                    <legend>Introduction</legend>
                    <p>
                      Thank you for your interest in the TRISA network for Travel Rule Compliance.
                      This multi-part form is the first step in the registration and certificate
                      issuance process. The information you provide will be used to verify the
                      legal entity that you represent and, where appropriate, will be available to
                      verified TRISA members to facilitate compliance-decisions.
                    </p>
                    <p>
                      To assist in completing the registration form, which is somewhat lengthy, the
                      form is broken into multiple sections, with information stored in your <em>local
                      browser storage</em> so that you can come back and complete the process. <strong>No
                      information is sent until you submit the form in the summary section</strong>.
                    </p>
                    <p className="text-danger">
                      The registration form is currently in its beta implementation. You will not be
                      able to submit the form now, but can download the form to submit it via the CLI
                      registration interface. The registration form will be completed when the TRISA
                      v1beta1 protocol is released.
                    </p>
                    <Button disabled variant="secondary">Reset</Button>{' '}
                    <Button variant="primary" onClick={(e) => this.setState({tabKey:"entity-details"})}>Next</Button>
                  </fieldset>
                </Tab.Pane>
                <Tab.Pane eventKey="entity-details">
                  <fieldset>
                    <legend>Entity Details</legend>
                    <p></p>
                  </fieldset>
                </Tab.Pane>
                <Tab.Pane eventKey="contacts">
                  <fieldset>
                    <legend>Contacts</legend>
                    <p>
                      Please supply contact information for representatives of your organization.
                      All contacts will receive an email verification token and the contact email
                      must be verified before the registration can proceed.
                    </p>

                    <Accordion defaultActiveKey="technical">
                      <Card>
                        <Accordion.Toggle as={Card.Header} eventKey="technical">
                          Technical Contact
                        </Accordion.Toggle>
                        <Accordion.Collapse eventKey="technical">
                          <Card.Body>
                            <p>
                              Primary contact for handling technical queries about the operation
                              and status of your service participating in the TRISA network.
                              Can be a group or admin email. (Required).
                            </p>
                            <Contact
                              contact={this.state.formData.contacts.technical}
                              onChange={this.createChangeHandler("technical", "contacts")}
                            />
                          </Card.Body>
                        </Accordion.Collapse>
                      </Card>
                      <Card>
                        <Accordion.Toggle as={Card.Header} eventKey="legal">
                          Legal/Compliance Contact
                        </Accordion.Toggle>
                        <Accordion.Collapse eventKey="legal">
                          <Card.Body>
                            <Contact
                              contact={this.state.formData.contacts.legal}
                              onChange={this.createChangeHandler("legal", "contacts")}
                            />
                          </Card.Body>
                        </Accordion.Collapse>
                      </Card>
                      <Card>
                        <Accordion.Toggle as={Card.Header} eventKey="administrative">
                          Administrative Contact
                        </Accordion.Toggle>
                        <Accordion.Collapse eventKey="administrative">
                          <Card.Body>
                            <Contact
                              contact={this.state.formData.contacts.administrative}
                              onChange={this.createChangeHandler("administrative", "contacts")}
                            />
                          </Card.Body>
                        </Accordion.Collapse>
                      </Card>
                      <Card>
                        <Accordion.Toggle as={Card.Header} eventKey="billing">
                          Billing Contact
                        </Accordion.Toggle>
                        <Accordion.Collapse eventKey="billing">
                          <Card.Body>
                            <Contact
                              contact={this.state.formData.contacts.billing}
                              onChange={this.createChangeHandler("billing", "contacts")}
                            />
                          </Card.Body>
                        </Accordion.Collapse>
                      </Card>
                    </Accordion>

                  </fieldset>
                </Tab.Pane>
                <Tab.Pane eventKey="summary">
                  <fieldset>
                    <legend>Summary</legend>
                    <Button type="submit">Submit</Button>
                  </fieldset>
                </Tab.Pane>
              </Tab.Content>
            </Col>
          </Row>
        </Tab.Container>
      </Form>
    );
  }
}

export default Registration;