// @flow
import React from 'react';
import { Col, Row } from 'react-bootstrap';
import { Link } from 'react-router-dom';

const ErrorPage = (): React$Element<React$FragmentType> => (
  <div className="account-pages pt-2 pt-sm-5 pb-4 pb-sm-5">
    <div className="container">
      <Row className="justify-content-center">
        <Col lg={4}>
          <div className="text-center">
            <img src="/oops-icon.png" alt="" height="200" />
            <h4 className="text-uppercase mt-3">Sorry, Something went wrong </h4>
            <p className="text-muted mt-3">
              An internal error occured. If the problem persists, contact trisa administrator
            </p>

            <Link className="btn btn-info mt-3" to="/">
              <i className="mdi mdi-reply" /> Return Home
            </Link>
          </div>
        </Col>
      </Row>
    </div>
  </div>
);

export default ErrorPage;
