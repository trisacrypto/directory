import React from 'react';
import TRIXO from './TRIXO';
import Contact from './Contact';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Tab from 'react-bootstrap/Tab';
import Nav from 'react-bootstrap/Nav';
import Form from 'react-bootstrap/Form';
import Card from 'react-bootstrap/Card';
import Modal from 'react-bootstrap/Modal';
import Button from 'react-bootstrap/Button';
import Spinner from 'react-bootstrap/Spinner';
import Accordion from 'react-bootstrap/Accordion';
import update from 'immutability-helper';
import LegalPerson from './ivms101/LegalPerson';
import gds from '../api/gds';
import { isTestNet } from '../lib/testnet';
import { Trans } from "@lingui/macro"


const testNet = isTestNet();
const registrationFormVersion = "v1beta1";
const submitButtonText = testNet ? "Submit TestNet Registration" : "Submit Production Registration"

// Returns a legal person object with default values populated.
const makeLegalPerson = () => {
  return {
    name: {
      name_identifiers: [{legal_person_name: "", legal_person_name_identifier_type: 1}],
      local_name_identifiers: [],
      phonetic_name_identifiers: [],
    },
    geographic_addresses: [{
      address_type: 2,
      address_line: ["", "", ""],
      country: "",
    }],
    customer_number: "",
    national_identification: {
      national_identifier: "",
      national_identifier_type: 9,
      country_of_issue: "",
      registration_authority: "",
    },
    country_of_registration: "",
  };
}

// Returns a TRIXO form with default values populated
const makeTRIXOForm = () => {
  return {
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
    safeguards_pii: false,
  }
}

const makeContacts = () => {
  const reducer = (contacts, contactType) => {
    contacts[contactType] = {
      name: "",
      email: "",
      phone: "",
    };
    return contacts;
  }
  return ["technical", "legal", "administrative", "billing"].reduce(reducer, {})
}

// Returns an empty form data state with default values populated
const makeFormData = () => {
  return {
    entity: makeLegalPerson(),
      contacts: makeContacts(),
      trisa_endpoint: "",
      common_name: "",
      website: "",
      business_category: 3,
      vasp_categories: [],
      established_on: "",
      trixo: makeTRIXOForm(),
  }
}

class Registration extends React.Component {
  state = {
    tabKey: "introduction",
    validated: false,
    formData: makeFormData(),
    formDownloadURL: "",
    submitting: false,
    showSubmittedModal: false,
    submissionResponse: {},
  }

  handleSubmit = (event) => {
    event.preventDefault();

    // Collect form and validate it
    const form = event.currentTarget;
    if (form.checkValidity() === false) {
      event.stopPropagation();
      this.setState({validated: true});
      this.props.onAlert("warning", "some fields could not be validated, please review form");
      return;
    }

    this.setState({submitting: true})
    // Submit the registration request to the server
    try {
      gds.register(this.state.formData)
        .then(rep => {
          this.setState({
            submitting: false,
            showSubmittedModal: true,
            submissionResponse: rep
          });
        })
        .catch(err => {
          this.setState({
            submitting: false,
            validated: false,
            showSubmittedModal: false,
            submissionResponse: {}
          });
          this.props.onAlert("danger", err.message);
          console.warn(err);
        });
    } catch (err) {
      this.setState({
        submitting: false,
        validated: false,
        showSubmittedModal: false,
        submissionResponse: {}
      });
      this.props.onAlert("danger", "There was a significant internal problem processing this form, please contact the admins for assistance.");
      console.error(err);
    }
  }

  handleReset = (event) => {
    event.preventDefault();
    this.setState({validated: false, formData: makeFormData()});
    this.deleteLocalStorage();
  }

  handleDownload = (event) => {
    const blob = new Blob([JSON.stringify({version: registrationFormVersion, registrationForm: this.state.formData}, null, "  ")]);
    const fileDownloadURL = URL.createObjectURL(blob);
    this.setState({fileDownloadURL: fileDownloadURL},
      () => {
        this.dofileDownload.click();
        URL.revokeObjectURL(fileDownloadURL);
        this.setState({fileDownloadURL: ""});
    });
  }

  handleModalClose = (event) => {
    this.setState({showSubmittedModal: false, submissionResponse: {}});
  }

  upload = (e) => {
    e.preventDefault();
    this.dofileUpload.click();
  }

  openFile = (e) => {
    const fileObj = e.target.files[0];
    const reader = new FileReader();

    let fileloaded = e => {
      // e.target.result is the file's content as text
      const fileContents = e.target.result;
      console.log(`File name: ${fileObj.name}, Length: ${fileContents.length} bytes.`);

      const data = JSON.parse(fileContents);
      if (data.version !== registrationFormVersion) {
        console.warn(`current form version is ${registrationFormVersion} cannot load version ${data.version}`);
        this.props.onAlert("danger", "Could not load data: invalid version");
        return
      }

      // TODO: validate the form data better
      this.setState({formData: data.registrationForm});
    }

    fileloaded = fileloaded.bind(this);
    reader.onload = fileloaded;
    reader.readAsText(fileObj);
  }

  saveToLocalStorage = (formData) => {
    try {
      const serialized = JSON.stringify(formData);
      localStorage.setItem("registrationForm", serialized);
    } catch (e) {
      console.warn(e);
    }
  }

  loadFromLocalStorage = () => {
    try {
      const serialized = localStorage.getItem("registrationForm");
      if (serialized === null) return undefined;
      console.log("data loaded from local storage");
      return JSON.parse(serialized);
    } catch (e) {
      console.warn(e);
      return undefined;
    }
  }

  deleteLocalStorage = () => {
    try {
      localStorage.removeItem("registrationForm");
      console.log("data deleted from local storage");
    } catch (e) {
      console.warn(e);
    }
  }

  // React Life Cycle
  componentDidMount() {
    const formData = this.loadFromLocalStorage();
    if (formData) {
      this.setState({formData: formData})
    }
  }

  componentDidUpdate(prevProps, prevState) {
    this.saveToLocalStorage(this.state.formData);
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
      const value = Array.from(event.target.selectedOptions, option => option.value);
      onChange(null, value);
    }
  }

  render() {
    const summaryFormData = JSON.stringify(this.state.formData, null, "  ");
    let submitBtn = <Button type="submit" variant="primary">{submitButtonText}</Button>;
    if (this.state.submitting) {
      submitBtn = (
        <Button variant="primary" disabled>
          <Spinner
            as="span"
            animation="border"
            size="sm"
            role="status"
            aria-hidden="true"
          />
          <span className="pl-2"><Trans>Submitting Registration &hellip;</Trans></span>
        </Button>
      );
    }

    return (
      <>
      <Form noValidate validated={this.state.validated} onSubmit={this.handleSubmit}>
        <Tab.Container id="registration-form" activeKey={this.state.tabKey} onSelect={(k) => this.setState({tabKey: k})}>
          <Row className="pt-3 pb-5">
            <Col sm={3}>
              <Nav variant="pills" className="flex-column">
                <Nav.Item>
                  <Nav.Link eventKey="introduction"><Trans>Introduction</Trans></Nav.Link>
                  <Nav.Link eventKey="basic-details"><Trans>Basic Details</Trans></Nav.Link>
                  <Nav.Link eventKey="legal-person"><Trans>Legal Person</Trans></Nav.Link>
                  <Nav.Link eventKey="contacts"><Trans>Contacts</Trans></Nav.Link>
                  <Nav.Link eventKey="trisa-implementation"><Trans>TRISA Implementation</Trans></Nav.Link>
                  <Nav.Link eventKey="trixo"><Trans>TRIXO Questionnaire</Trans></Nav.Link>
                  <Nav.Link eventKey="summary"><Trans>Summary</Trans></Nav.Link>
                </Nav.Item>
              </Nav>
            </Col>
            <Col sm={9}>
              <Tab.Content>
                <Tab.Pane eventKey="introduction">
                  <fieldset>
                    <legend><Trans>Introduction</Trans></legend>
                    <p>
                      <Trans>Thank you for your interest in the TRISA network for Travel Rule Compliance.
                      This multi-part form is the first step in the registration and certificate
                      issuance process. The information you provide will be used to verify the
                      legal entity that you represent and, where appropriate, will be available to
                      verified TRISA members to facilitate compliance decisions.</Trans>
                    </p>
                    <p>
                      <Trans>To assist in completing the registration form, which is somewhat lengthy, the
                      form is broken into multiple sections, with information stored in your <em>local
                      browser storage</em> so that you can come back and complete the process. <strong>No
                      information is sent until you submit the form in the summary section</strong>.</Trans>
                    </p>
                    <p>
                      <Trans>This registration form is currently in its beta implementation. On the summary page
                      you are able to download the form to save offline. You may also load a saved form below.</Trans>
                    </p>
                    <Form.Group>
                      <Button variant="primary" onClick={(e) => this.setState({tabKey:"basic-details"})}><Trans>Next</Trans></Button>{' '}
                      <Button type="button" onClick={this.upload} variant="info"><Trans>Load</Trans></Button>
                      <input type="file" className="d-none"
                        multiple={false}
                        accept=".json,application/json"
                        onChange={e=>this.openFile(e)}
                        ref={e=>this.dofileUpload=e}
                      />
                    </Form.Group>
                  </fieldset>
                </Tab.Pane>
                <Tab.Pane eventKey="basic-details">
                  <fieldset>
                    <legend><Trans>Basic Details</Trans></legend>
                    <p>
                      <Trans>To get started, please tell us a bit about your organization.</Trans>
                    </p>
                    <Form.Group>
                      <Form.Label><Trans>Website</Trans></Form.Label>
                      <Form.Control
                        type="url"
                        value={this.state.formData.website}
                        onChange={this.createFlatChangeHandler("website")}
                        required={true}
                      />
                    </Form.Group>
                    <Form.Group>
                      <Form.Label><Trans>Date of Incorporation/Establishment</Trans></Form.Label>
                      <Form.Control
                        type="date"
                        value={this.state.formData.established_on}
                        onChange={this.createFlatChangeHandler("established_on")}
                      />
                    </Form.Group>
                    <Form.Group>
                      <Form.Label><Trans>Business Category</Trans></Form.Label>
                      <Form.Control
                        as="select" custom
                        value={this.state.formData.business_category}
                        onChange={this.createIntChangeHandler("business_category")}
                      >
                        <option value={0}></option>
                        <option value={1}><Trans>Private Organization</Trans></option>
                        <option value={2}><Trans>Government Entity</Trans></option>
                        <option value={3}><Trans>Business Entity</Trans></option>
                        <option value={4}><Trans>Non-Commercial Entity</Trans></option>
                      </Form.Control>
                      <Form.Text className="text-muted">
                        <Trans>Please select the entity category that most closely matches your organization.</Trans>
                      </Form.Text>
                    </Form.Group>
                    <Form.Group>
                      <Form.Label><Trans>VASP Category</Trans></Form.Label>
                      <Form.Control
                        as="select" custom multiple
                        value={this.state.formData.vasp_categories}
                        onChange={this.createMultiselectChangeHandler("vasp_categories")}
                      >
                        <option value="Exchange"><Trans>Centralized Exchange</Trans></option>
                        <option value="DEX"><Trans>Decentralized Exchange</Trans></option>
                        <option value="P2P"><Trans>Person-to-Person Exchange</Trans></option>
                        <option value="Kiosk"><Trans>Kiosk / Crypto ATM Operator</Trans></option>
                        <option value="Custodian"><Trans>Custody Provider</Trans></option>
                        <option value="OTC"><Trans>Over-The-Counter Trading Desk</Trans></option>
                        <option value="Fund"><Trans>Investment Fund - hedge funds, ETFs, and family offices</Trans></option>
                        <option value="Project"><Trans>Token Project</Trans></option>
                        <option value="Gambling"><Trans>Gambling or Gaming Site</Trans></option>
                        <option value="Miner"><Trans>Mining Pool</Trans></option>
                        <option value="Mixer"><Trans>Mixing Service</Trans></option>
                        <option value="Individual"><Trans>Legal person</Trans></option>
                        <option value="Other"><Trans>Other</Trans></option>
                      </Form.Control>
                      <Form.Text className="text-muted">
                        <Trans>Please select as many categories needed to represent the types of virtual asset services your organization provides.</Trans>
                      </Form.Text>
                    </Form.Group>
                  </fieldset>
                  <Form.Group>
                    <Button variant="secondary" onClick={(e) => this.setState({tabKey:"introduction"})}>Back</Button>{' '}
                    <Button variant="primary" onClick={(e) => this.setState({tabKey:"legal-person"})}>Next</Button>
                  </Form.Group>
                </Tab.Pane>
                <Tab.Pane eventKey="legal-person">
                  <fieldset>
                    <legend className="legend">Legal Person</legend>
                    <p>
                      <Trans>Please enter the information that identify your organization as a
                      Legal Person. This form represents the IVMS 101 data structure for
                      legal persons and is strongly suggested for use as KYC information
                      exchanged in TRISA transfers.</Trans>
                    </p>
                    <LegalPerson
                      person={this.state.formData.entity}
                      onChange={this.createChangeHandler("entity")}
                    />
                  </fieldset>
                  <Form.Group>
                    <Button variant="secondary" onClick={(e) => this.setState({tabKey:"basic-details"})}><Trans>Back</Trans></Button>{' '}
                    <Button variant="primary" onClick={(e) => this.setState({tabKey:"contacts"})}><Trans>Next</Trans></Button>
                  </Form.Group>
                </Tab.Pane>
                <Tab.Pane eventKey="contacts">
                  <fieldset>
                    <legend><Trans>Contacts</Trans></legend>
                    <p>
                      <Trans>Please supply contact information for representatives of your organization.
                      All contacts will receive an email verification token and the contact email
                      must be verified before the registration can proceed.</Trans>
                    </p>

                    <Accordion defaultActiveKey="technical" className="pb-3">
                      <Card>
                        <Accordion.Toggle as={Card.Header} eventKey="technical">
                          <Trans>Technical Contact</Trans>
                        </Accordion.Toggle>
                        <Accordion.Collapse eventKey="technical">
                          <Card.Body>
                            <p>
                              <Trans>Primary contact for handling technical queries about the operation
                              and status of your service participating in the TRISA network.
                              Can be a group or admin email. (Required).</Trans>
                            </p>
                            <Contact
                              contact={this.state.formData.contacts.technical}
                              onChange={this.createChangeHandler("technical", "contacts")}
                              required={true}
                            />
                          </Card.Body>
                        </Accordion.Collapse>
                      </Card>
                      <Card>
                        <Accordion.Toggle as={Card.Header} eventKey="legal">
                          <Trans>Legal/Compliance Contact</Trans>
                        </Accordion.Toggle>
                        <Accordion.Collapse eventKey="legal">
                          <Card.Body>
                            <p>
                              <Trans>Compliance officer or legal contact for requests about the compliance
                              requirements and legal status of your organization. (Strongly recommended).</Trans>
                            </p>
                            <Contact
                              contact={this.state.formData.contacts.legal}
                              onChange={this.createChangeHandler("legal", "contacts")}
                              required={false}
                            />
                          </Card.Body>
                        </Accordion.Collapse>
                      </Card>
                      <Card>
                        <Accordion.Toggle as={Card.Header} eventKey="administrative">
                          <Trans>Administrative Contact</Trans>
                        </Accordion.Toggle>
                        <Accordion.Collapse eventKey="administrative">
                          <Card.Body>
                            <p>
                              <Trans>Administrative or executive contact for your organization to field
                              high-level requests or queries. (Required).</Trans>
                            </p>
                            <Contact
                              contact={this.state.formData.contacts.administrative}
                              onChange={this.createChangeHandler("administrative", "contacts")}
                              required={true}
                            />
                          </Card.Body>
                        </Accordion.Collapse>
                      </Card>
                      <Card>
                        <Accordion.Toggle as={Card.Header} eventKey="billing">
                          <Trans>Billing Contact</Trans>
                        </Accordion.Toggle>
                        <Accordion.Collapse eventKey="billing">
                          <Card.Body>
                            <p>
                              <Trans>Billing contact for your organization to handle account and invoice
                              requests or queries relating to the operation of the TRISA network. (Optional).</Trans>
                            </p>
                            <Contact
                              contact={this.state.formData.contacts.billing}
                              onChange={this.createChangeHandler("billing", "contacts")}
                              requried={false}
                            />
                          </Card.Body>
                        </Accordion.Collapse>
                      </Card>
                    </Accordion>

                    <Form.Group>
                      <Button variant="secondary" onClick={(e) => this.setState({tabKey:"legal-person"})}>Back</Button>{' '}
                      <Button variant="primary" onClick={(e) => this.setState({tabKey:"trisa-implementation"})}>Next</Button>
                    </Form.Group>
                  </fieldset>
                </Tab.Pane>
                <Tab.Pane eventKey="trisa-implementation">
                  <fieldset>
                    <legend><Trans>TRISA Implementation</Trans></legend>
                    <p>
                      <Trans>Each VASP is required to establish a TRISA endpoint for inter-VASP
                      communication. Please specify the details of your endpoint for
                      certificate issuance.</Trans>
                    </p>
                    <Form.Group>
                      <Form.Label><Trans>TRISA Endpoint</Trans></Form.Label>
                      <Form.Control
                        type="url"
                        value={this.state.formData.trisa_endpoint}
                        onChange={this.createFlatChangeHandler("trisa_endpoint")}
                        placeholder="trisa.example.com:443"
                      />
                      <Form.Text className="text-muted">
                        <Trans>The address and port of the TRISA endpoint for partner VASPs to connect on via gRPC.</Trans>
                      </Form.Text>
                    </Form.Group>
                    <Form.Group>
                      <Form.Label><Trans>Certificate Common Name</Trans></Form.Label>
                      <Form.Control
                        type="text"
                        value={this.state.formData.common_name}
                        onChange={this.createFlatChangeHandler("common_name")}
                        placeholder="trisa.example.com"
                      />
                      <Form.Text className="text-muted">
                        <Trans>The common name for the mTLS certificate. This should match the TRISA endpoint without the port in most cases.</Trans>
                      </Form.Text>
                    </Form.Group>
                    <Form.Group>
                      <Button variant="secondary" onClick={(e) => this.setState({tabKey:"contacts"})}><Trans>Back</Trans></Button>{' '}
                      <Button variant="primary" onClick={(e) => this.setState({tabKey:"trixo"})}><Trans>Next</Trans></Button>
                    </Form.Group>
                  </fieldset>
                </Tab.Pane>
                <Tab.Pane eventKey="trixo">
                  <fieldset>
                    <legend><Trans>TRIXO Questionnaire</Trans></legend>
                    <p>
                      <Trans>This questionnaire is designed to help the TRISA working group and TRISA members
                      understand the regulatory regime of your organization. The information you provide
                      will help ensure that required compliance information exchanges are conducted
                      correctly and safely. All verified TRISA members will have access to this information.</Trans>
                    </p>
                    <TRIXO
                      data={this.state.formData.trixo}
                      onChange={this.createChangeHandler("trixo")}
                    />
                    <Form.Group>
                      <Button variant="secondary" onClick={(e) => this.setState({tabKey:"trisa-implementation"})}>Back</Button>{' '}
                      <Button variant="primary" onClick={(e) => this.setState({tabKey:"summary"})}>Next</Button>
                    </Form.Group>
                  </fieldset>
                </Tab.Pane>
                <Tab.Pane eventKey="summary">
                  <fieldset>
                    <legend><Trans>Summary</Trans></legend>
                    <div className="mt-2 mb-5"><pre className="summaryJSON">{summaryFormData}</pre></div>
                    <Form.Group className="mt-2">
                      <Row>
                        <Col xs={6}>
                          {submitBtn}
                        </Col>
                        <Col xs={6} className="text-right">
                          <Button type="button" variant="info" onClick={this.handleDownload}><Trans>Download</Trans></Button>{' '}
                          <Button type="reset" variant="secondary" onClick={this.handleReset}><Trans>Reset</Trans></Button>
                          <a className="d-none"
                            download="trisa_registration.json"
                            href={this.state.fileDownloadURL}
                            ref={e=>this.dofileDownload = e}
                          >
                          <Trans>download data</Trans>
                          </a>
                        </Col>
                      </Row>
                    </Form.Group>
                  </fieldset>
                </Tab.Pane>
              </Tab.Content>
            </Col>
          </Row>
        </Tab.Container>
      </Form>
      <Modal
        show={this.state.showSubmittedModal}
        onHide={this.handleModalClose}
        backdrop="static"
        keyboard={false}
        centered
        size="lg"
      >
        <Modal.Header closeButton>
          <Modal.Title><Trans>TRISA Registration Request Submitted!</Trans></Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <p>
            <Trans>Your registration request has been successfully received by the Directory Service.
            Verification emails have been sent to all contacts listed. Once your contact
            information has been verified, the registration form will be sent to the
            TRISA review board to verify your membership in the TRISA network.</Trans>
          </p>
          <p>
            <Trans>When you are verified you will be issued PKCS12 encrypted identity certificates
            for use in mTLS authentication between TRISA members. The password to decrypt
            those certificates is shown below:</Trans>
          </p>
          <p className="text-center mark"><strong>
            {this.state.submissionResponse.pkcs12password}
          </strong></p>
          <p className="text-center text-danger">
            <Trans>This is the only time the PKCS12 password is shown during the registration process.
            <br />
            Please copy and paste this password and store somewhere safe!</Trans>
          </p>
          <p className="text-muted text-center">
            <Trans>ID: {this.state.submissionResponse.id}
            <br />
            Verification Status: "{this.state.submissionResponse.status}"</Trans>
          </p>
          <p className="text-muted small">
            <Trans>Message from server: "{this.state.submissionResponse.message}"</Trans>
          </p>
        </Modal.Body>
        <Modal.Footer className="text-center">
          <Button variant="danger" onClick={this.handleModalClose}>
            <Trans>Understood</Trans>
          </Button>
        </Modal.Footer>
      </Modal>
      </>
    );
  }
}

export default Registration;