// @flow
import React from 'react';
import { Col, Row } from 'react-bootstrap';
import { Link, useLocation } from 'react-router-dom';

const ErrorPageNotFoundAlt = (): React$Element<React$FragmentType> => {
  const { state } = useLocation();

  return (
    <div className="account-pages pt-2 pt-sm-5 pb-4 pb-sm-5">
      <div className="container">
        <Row className="justify-content-center">
          <Col lg={4}>
            <div className="text-center">
              <h1 className="text-error mt-4">404</h1>
              <h4 className="text-uppercase text-danger mt-3">Page Not Found</h4>
              <p className="text-muted mt-3">
                {state ? state.error : "It's looking like you may have taken a wrong turn"}
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
};

export default ErrorPageNotFoundAlt;
