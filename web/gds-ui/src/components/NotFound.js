import React from 'react';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import withTracker from '../lib/analytics';
import { Trans } from "@lingui/macro";


const NotFound = () => {
  return (
    <Row>
      <Col md={{span: 6, offset: 3}} className="text-center">
        <p className="big-number">404</p>
        <h4><Trans>PAGE NOT FOUND</Trans></h4>
        <p className="text-muted">
          <Trans>The page you are looking for might have been removed, had its name changed, or is temporarily unavailable.</Trans>
        </p>
        <a href="/" className="btn btn-secondary mt-2"><Trans>Directory Home</Trans></a>
      </Col>
    </Row>
  );
};

export default withTracker(NotFound);