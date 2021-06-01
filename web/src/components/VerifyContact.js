import React from 'react';
import gds from '../api/gds';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Alert from 'react-bootstrap/Alert';
import Spinner from 'react-bootstrap/Spinner';

class VerifyContact extends React.Component {
  constructor(props) {
    super(props);
    const params = new URLSearchParams(window.location.search);
    this.state = {
      token: params.get("token"),
      vaspID: params.get("vaspID"),
      loading: false,
      header: "Unable to Verify Contact",
      message: null,
      variant: "danger",
      verificationStatus: null
    }
  }

  componentDidMount() {
    if (!this.state.token || !this.state.vaspID) {
      const message = "Contact verification requires both the registration ID of your organization as well as a unique token that was sent to your email address. Please check the link in your email or copy and paste the complete link into your browser.";
      this.setState({variant: "danger", message: message});
      return
    }

    this.setState({loading: true});
    gds.verifyContact(this.state.vaspID, this.state.token)
      .then(rep => {
        this.setState({loading: false});
        if (rep.error) {
          this.setState({variant: "warning", message: rep.error.message});
        } else {
          this.setState({
            variant: "success",
            message: rep.message,
            header: "Contact Verified",
            verificationStatus: rep.status
          });
        }

      })
      .catch(err => {
        this.setState({loading: false});
        this.setState({variant: "danger", message: err.message});
        console.warn(err);
      });
  }

  render() {
    if (this.state.loading) {
      return (
        <div className="d-flex flex-column align-items-center justify-content-center">
          <div className="row">
            <Spinner animation="border" role="status">
              <span className="sr-only">Loading ...</span>
            </Spinner>
          </div>
          <div className="row pt-1">
            <em>Verifying Contact Information &hellip;</em>
          </div>
        </div>
      );
    }

    const vs = this.state.verificationStatus ? <p>Verification Status: {this.state.verificationStatus}</p> : null;

    return (
      <Row>
        <Col md={{span: 8, offset: 2}} className="text-center">
          <Alert variant={this.state.variant}>
            <h5>{this.state.header}</h5>
            <p>{this.state.message}</p>
            {vs}
            <hr />
            <a href="/" className="btn btn-dark">Return to Directory</a>
          </Alert>
        </Col>
      </Row>
    );
  }
}

export default VerifyContact;