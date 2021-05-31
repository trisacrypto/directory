import React from 'react';
import Card from 'react-bootstrap/Card';
import countryCodeEmoji from '../lib/country';

const LookupResults = (props) => {
  if (props.results && Object.keys(props.results).length === 0 && props.results.constructor === Object) {
    return null;
  }

  let results = props.results;

  return (
    <div className="lookup-results">
      <Card className="mb-5">
        <Card.Header>Global TRISA Directory Record</Card.Header>
        <Card.Body>
          <Card.Title>{results.name} <small className="text-muted">{results.commonName}</small></Card.Title>
          <dl className="row mt-4">
            <dt class="col-sm-3">Registered Directory</dt>
            <dd class="col-sm-9">{results.registeredDirectory}</dd>

            <dt class="col-sm-3">TRISA Member ID</dt>
            <dd class="col-sm-9">{results.id}</dd>

            <dt class="col-sm-3">TRISA Service Endpoint</dt>
            <dd class="col-sm-9">{results.endpoint}</dd>

            <dt class="col-sm-3">Country</dt>
            <dd class="col-sm-9">{countryCodeEmoji(results.country)} <span className="sr-only">{results.country}</span></dd>

            <dt class="col-sm-3">TRISA Verification</dt>
            <dd class="col-sm-9">VERIFIED on {results.verifiedOn}</dd>

            <dt class="col-sm-3">TRISA Identity Signature</dt>
            <dd class="col-sm-9">{results.identityCertificate.signature}</dd>
          </dl>
        </Card.Body>
      </Card>
    </div>
  );
}

export default LookupResults;