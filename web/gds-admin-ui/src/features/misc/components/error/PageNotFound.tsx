// @flow
import React from 'react';
import { Card, Col, Row } from 'react-bootstrap';
import { Link, useLocation } from 'react-router-dom';

import { getDirectoryLogo } from '@/utils';

const ErrorPageNotFound = () => {
  const { state } = useLocation<{ error: any }>();
  return (
    <>
      <div className="account-pages pt-2 pt-sm-5 pb-4 pb-sm-5">
        <div className="container">
          <Row className="justify-content-center">
            <Col md={8} lg={6} xl={5} xxl={4}>
              <Card>
                {/* logo */}
                <Card.Header
                  className="pt-4 pb-4 text-center bg-primary"
                  style={{ background: 'linear-gradient(90deg,#24a9df,#1aebb4)' }}
                >
                  <Link to="/">
                    <span>
                      <img src={getDirectoryLogo()} alt="" height="38" />
                    </span>
                  </Link>
                </Card.Header>

                <Card.Body className="p-4">
                  <div className="text-center">
                    <h1 className="text-error">
                      4<i className="mdi mdi-emoticon-sad" />4
                    </h1>
                    <h4 className="text-uppercase text-danger mt-3">Page Not Found</h4>
                    <p className="text-muted mt-3">
                      {state?.error
                              ? state?.error
                              : "It's looking like you may have taken a wrong turn."}
                    </p>

                    <Link className="btn btn-primary mt-3" to="/">
                      <i className="mdi mdi-reply" /> Return Home
                    </Link>
                  </div>
                </Card.Body>
              </Card>
            </Col>
          </Row>
        </div>
      </div>

      <footer className="footer footer-alt">
        Created and maintained by <a href="https://rotational.io">Rotational</a> Labs in partnership
        with <a href="https://ciphertrace.com">CipherTrace</a> on behalf of{' '}
        <a href="https://travelrule.io">TRISA</a>.
      </footer>
    </>
  );
};

export default ErrorPageNotFound;
