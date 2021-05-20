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
import LegalPerson from './ivms101/LegalPerson';


class Registration extends React.Component {
  state = {
    tabKey: "introduction",
    validated: false,
    formData: {
      entity: {
        name: {
          name_identifiers: [{legal_person_name: "", legal_person_name_identifier_type: 0}],
          local_name_identifiers: [],
          phonetic_name_identifiers: [],
        },
        geographic_addresses: [{
          address_type: 1,
          address_line: ["", "", ""],
          country: "",
        }],
        customer_number: "",
        national_identification: {
          national_identifier: "",
          national_identifier_type: 8,
          country_of_issue: "",
          registration_authority: "",
        },
        country_of_registration: "",
      },
      contacts: {
        technical: {name: "", email: "", phone: ""},
        legal: {name: "", email: "", phone: ""},
        administrative: {name: "", email: "", phone: ""},
        billing: {name: "", email: "", phone: ""},
      },
      trisa_endpoint: "",
      common_name: "",
      website: "",
      business_category: 0,
      vasp_categories: [],
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

  createFlatChangeHandler = (field, ...parents) => {
    const onChange = this.createChangeHandler(field, ...parents);
    return (event) => {
      onChange(null, event.target.value);
    }
  }

  createIntChangeHandler = (field, ...parents) => {
    const onChange = this.createChangeHandler(field, ...parents);
    return (event) => {
      onChange(null, parseInt(event.target.value));
    }
  }

  createMultiselectChangeHandler = (field, ...parents) => {
    const onChange = this.createChangeHandler(field, ...parents);
    return (event) => {
      const value = Array.from(event.target.selectedOptions, option => parseInt(option.value));
      onChange(null, value);
    }
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
                    <p>
                      To get started, please tell us a bit about your organization. In addition to
                      some basic organizational details, we'll collect IVMS 101 LegalPerson data,
                      which is required KYC information for TRISA compliance information transfers.
                    </p>
                    <fieldset>
                      <legend className="sublegend">Basic Details</legend>
                      <Form.Group>
                        <Form.Label>Website</Form.Label>
                        <Form.Control
                          type="url"
                          value={this.state.formData.website}
                          onChange={this.createFlatChangeHandler("website")}
                        />
                      </Form.Group>
                      <Form.Group>
                        <Form.Label>Date of Incorporation/Establishment</Form.Label>
                        <Form.Control
                          type="date"
                          value={this.state.formData.established_on}
                          onChange={this.createFlatChangeHandler("established_on")}
                        />
                      </Form.Group>
                      <Form.Group>
                        <Form.Label>Business Category</Form.Label>
                        <Form.Control
                          as="select" custom
                          value={this.state.formData.business_category}
                          onChange={this.createIntChangeHandler("business_category")}
                        >
                          <option value={0}></option>
                          <option value={1}>Private Organization</option>
                          <option value={2}>Government Entity</option>
                          <option value={3}>Business Entity</option>
                          <option value={4}>Non-Commercial Entity</option>
                        </Form.Control>
                        <Form.Text className="text-muted">
                          Please select the entity category that most closely matches your organization.
                        </Form.Text>
                      </Form.Group>
                      <Form.Group>
                        <Form.Label>VASP Category</Form.Label>
                        <Form.Control
                          as="select" custom multiple
                          value={this.state.formData.vasp_categories}
                          onChange={this.createMultiselectChangeHandler("vasp_categories")}
                        >
                          <option value={1}>Exchange</option>
                          <option value={2}>DEX</option>
                          <option value={3}>P2P Vendor</option>
                          <option value={4}>P2P Exchange</option>
                          <option value={5}>Custodial Wallet/Custody Provider</option>
                          <option value={6}>Non-Custodial Wallet</option>
                          <option value={7}>Mixer</option>
                          <option value={8}>Crypto ATM Provider</option>
                          <option value={9}>Crypto ATM Services</option>
                          <option value={10}>Crypto ATM Technology Provider</option>
                          <option value={11}>Crypto ATM Operator</option>
                          <option value={12}>OTC Desk</option>
                          <option value={13}>Payment Processor</option>
                          <option value={14}>P2P Payment Service</option>
                          <option value={15}>Crypto Hedge Fund</option>
                          <option value={16}>Blockchain (DLT) Project</option>
                          <option value={17}>Family Office</option>
                          <option value={18}>Gambling</option>
                        </Form.Control>
                        <Form.Text className="text-muted">
                          Please select as many categories needed to represent the types of virtual asset services your organization provides.
                        </Form.Text>
                      </Form.Group>
                    </fieldset>
                    <fieldset>
                      <legend className="sublegend">Legal Person</legend>
                      <LegalPerson
                        person={this.state.formData.entity}
                        onChange={this.createChangeHandler("entity")}
                      />
                    </fieldset>
                    <Form.Group>
                      <Button variant="secondary" onClick={(e) => this.setState({tabKey:"introduction"})}>Back</Button>{' '}
                      <Button variant="primary" onClick={(e) => this.setState({tabKey:"trisa-implementation"})}>Next</Button>
                    </Form.Group>
                  </fieldset>
                </Tab.Pane>
                <Tab.Pane eventKey="trisa-implementation">
                  <fieldset>
                    <legend>TRISA Implementation</legend>
                    <p>
                      Each VASP is required to establish a TRISA endpoint for inter-VASP
                      communication. Please specify the details of your endpoint for
                      certificate issuance.
                    </p>
                    <Form.Group>
                      <Form.Label>TRISA Endpoint</Form.Label>
                      <Form.Control
                        type="url"
                        value={this.state.formData.trisa_endpoint}
                        onChange={this.createFlatChangeHandler("trisa_endpoint")}
                        placeholder="trisa.example.com:443"
                      />
                      <Form.Text className="text-muted">
                        The address and port of the TRISA endpoint for partner VASPs to connect on via gRPC.
                      </Form.Text>
                    </Form.Group>
                    <Form.Group>
                      <Form.Label>Certificate Common Name</Form.Label>
                      <Form.Control
                        type="text"
                        value={this.state.formData.common_name}
                        onChange={this.createFlatChangeHandler("common_name")}
                        placeholder="trisa.example.com"
                      />
                      <Form.Text className="text-muted">
                        The common name for the mTLS certificate. This should match the TRISA endpoint without the port in most cases.
                      </Form.Text>
                    </Form.Group>
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
                    <div><pre>{summaryFormData}</pre></div>
                    <Form.Group>
                      <Button type="reset" disabled variant="secondary">Reset</Button>{' '}
                      <Button type="submit" disabled variant="primary">Download</Button>
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