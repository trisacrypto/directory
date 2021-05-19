import React from 'react';
import TRIXO from './TRIXO';
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
      entity: {},
      contacts: {
        technical: {name: "", email: "", phone: ""},
        legal: {name: "", email: "", phone: ""},
        administrative: {name: "", email: "", phone: ""},
        billing: {name: "", email: "", phone: ""},
      },
      trisa_endpoint: "",
      common_name: "",
      website: "",
      business_category: "",
      established_on: "",
      trixo: {
        primary_national_jurisdiction: "",
        primary_regulator: "",
        other_jurisdictions: [],
        financial_transfers_permitted: "",
        has_required_regulatory_program: "",
        conducts_customer_kyc: false,
        kyc_threshold: 0.0,
        kyc_threshold_currency: "USD",
        must_comply_travel_rule: false,
        applicable_regulations: ["FATF Recommendation 16"],
        compliance_threshold: 0.0,
        compliance_threshold_currency: "USD",
        must_safeguard_pii: false,
        safegaurds_pii: false,
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
    const summaryFormData = JSON.stringify(this.state.formData, null, "  ");

    return (
      <Form noValidate validated={this.state.validated} onSubmit={this.handleSubmit}>
        <Tab.Container id="registration-form" activeKey={this.state.tabKey} onSelect={(k) => this.setState({tabKey: k})}>
          <Row className="pt-3 pb-5">
            <Col sm={3}>
              <Nav variant="pills" className="flex-column">
                <Nav.Item>
                  <Nav.Link eventKey="introduction">Introduction</Nav.Link>
                  <Nav.Link eventKey="entity-details">Entity Details</Nav.Link>
                  <Nav.Link eventKey="trisa-implementation">TRISA Implementation</Nav.Link>
                  <Nav.Link eventKey="contacts">Contacts</Nav.Link>
                  <Nav.Link eventKey="trixo">TRIXO Questionnaire</Nav.Link>
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
                    <Form.Group>
                      <Button type="reset" disabled variant="secondary">Reset</Button>{' '}
                      <Button variant="primary" onClick={(e) => this.setState({tabKey:"entity-details"})}>Next</Button>
                    </Form.Group>
                  </fieldset>
                </Tab.Pane>
                <Tab.Pane eventKey="entity-details">
                  <fieldset>
                    <legend>Entity Details</legend>
                    <p></p>
                    <Form.Group>
                      <Button variant="secondary" onClick={(e) => this.setState({tabKey:"introduction"})}>Back</Button>{' '}
                      <Button variant="primary" onClick={(e) => this.setState({tabKey:"trisa-implementation"})}>Next</Button>
                    </Form.Group>
                  </fieldset>
                </Tab.Pane>
                <Tab.Pane eventKey="trisa-implementation">
                  <fieldset>
                    <legend>TRISA Implementation</legend>
                    <p></p>
                    <Form.Group>
                      <Button variant="secondary" onClick={(e) => this.setState({tabKey:"entity-details"})}>Back</Button>{' '}
                      <Button variant="primary" onClick={(e) => this.setState({tabKey:"contacts"})}>Next</Button>
                    </Form.Group>
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

                    <Accordion defaultActiveKey="technical" className="pb-3">
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
                            <p>
                              Compliance officer or legal contact for requests about the compliance
                              requirements and legal status of your organization. (Strongly recommended).
                            </p>
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
                            <p>
                              Administrative or executive contact for your organization to field
                              high-level requests or queries. (Optional).
                            </p>
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
                            <p>
                              Billing contact for your organization to handle account and invoice
                              requests or queries relating to the operation of the TRISA network. (Optional).
                            </p>
                            <Contact
                              contact={this.state.formData.contacts.billing}
                              onChange={this.createChangeHandler("billing", "contacts")}
                            />
                          </Card.Body>
                        </Accordion.Collapse>
                      </Card>
                    </Accordion>

                    <Form.Group>
                      <Button variant="secondary" onClick={(e) => this.setState({tabKey:"trisa-implementation"})}>Back</Button>{' '}
                      <Button variant="primary" onClick={(e) => this.setState({tabKey:"trixo"})}>Next</Button>
                    </Form.Group>
                  </fieldset>
                </Tab.Pane>
                <Tab.Pane eventKey="trixo">
                  <fieldset>
                    <legend>TRIXO Questionnaire</legend>
                    <p>
                      This questionnaire is designed to help the TRISA working group and other
                      TRISA members understand the regulatory regime of your organization to ensure
                      that required compliance information exchanges are conducted correctly and
                      safely. All verified TRISA members will have access to this information.
                    </p>
                    <TRIXO
                      data={this.state.formData.trixo}
                      onChange={this.createChangeHandler("trixo")}
                    />
                    <Form.Group>
                      <Button variant="secondary" onClick={(e) => this.setState({tabKey:"contacts"})}>Back</Button>{' '}
                      <Button variant="primary" onClick={(e) => this.setState({tabKey:"summary"})}>Next</Button>
                    </Form.Group>
                  </fieldset>
                </Tab.Pane>
                <Tab.Pane eventKey="summary">
                  <fieldset>
                    <legend>Summary</legend>
                    <p><pre>{summaryFormData}</pre></p>
                    <Form.Group>
                      <Button type="reset" disabled variant="secondary">Reset</Button>{' '}
                      <Button type="submit">Download</Button>
                    </Form.Group>
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