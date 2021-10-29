import React, { useState } from 'react'
import gds from '../api/gds';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import LookupForm from './LookupForm';
import LookupResults from './LookupResults';
import withTracker from '../lib/analytics';

const Lookup = ({ onAlert }) => {
  const [results, setResults] = useState({});

  const onLookup = async (query, inputType) => {
    try {
      const response = await gds.lookup(query, inputType);
      if (response.error) {
        onAlert('warning', response.error.message);
      } else {
        setResults(response);
      }
    } catch(err) {
      onAlert('danger', err.message);
      console.warn(err);
    }
  }

  return (
    <>
    <Row className="py-3">
      <Col md={{span: 8, offset: 2}}>
        <LookupForm onSubmit={onLookup} />
      </Col>
    </Row>
    <Row>
      <Col>
        <LookupResults results={results} />
      </Col>
    </Row>
    </>
  );
}

export default withTracker(Lookup);