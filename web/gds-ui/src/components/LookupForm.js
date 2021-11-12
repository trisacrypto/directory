import React from 'react';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import { Trans } from "@lingui/macro";
import { t } from "@lingui/macro";


class Lookup extends React.Component {
  state = { query: '', inputType: ''};

  onFormSubmit = (event) => {
    event.preventDefault();
    this.props.onSubmit(this.state.query, this.state.inputType);
  }

  uuidRE = /^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/
  cnameRE = /^[0-9a-zA-Z.-]+$/

  onTextInput = (event) => {
    let inputType = '';
    if (this.uuidRE.test(event.target.value)) {
      inputType = 'uuid';
    } else if (this.cnameRE.test(event.target.value)) {
      inputType = 'common name'
    }
    this.setState({ query: event.target.value, inputType: inputType });
  }

  render() {
    const detectedType = this.state.inputType !== '' ? t`Detected input type: ${this.state.inputType}` : '';

    return (
      <Form className="justify-content-center" onSubmit={this.onFormSubmit}>
        <Form.Row className="align-items-top">
          <Col>
            <Form.Label htmlFor="lookupInput" srOnly>
              <Trans>Common Name or VASP ID</Trans>
            </Form.Label>
            <Form.Control
              id="lookupInput"
              placeholder={t`Common Name or VASP ID`}
              onChange={this.onTextInput}
            />
            <Form.Text id="passwordHelpBlock" muted>
              {detectedType}
            </Form.Text>
          </Col>
          <Col xs="auto">
            <Button type="submit"><Trans>Submit</Trans></Button>
          </Col>
        </Form.Row>
      </Form>
    );
  }
}

export default Lookup;