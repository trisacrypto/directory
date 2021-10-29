import React from 'react';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import withTracker from '../lib/analytics';

const NotFound = () => {
  return (
    <Row>
      <Col md={{span: 6, offset: 3}} className="text-center">
        <p className="big-number">404</p>
        <h4>PAGE NOT FOUND</h4>
        <p className="text-muted">
          The page you are looking for might have been removed, had its name changed, or is temporarily unavailable.
        </p>
        <a href="/" className="btn btn-secondary mt-2">Directory Home</a>
      </Col>
    </Row>
  );
};

export default withTracker(NotFound);