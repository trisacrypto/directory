import React from 'react';
import Card from 'react-bootstrap/Card';
import { countryCodeEmoji, getCountryName } from '../lib/country';
import { Trans } from "@lingui/macro"


const LookupResults = (props) => {
  if (props.results && Object.keys(props.results).length === 0 && props.results.constructor === Object) {
    return null;
  }

  let results = props.results;

  return (
    <div className="lookup-results">
      <Card className="mb-5">
        <Card.Header><Trans>Global TRISA Directory Record</Trans></Card.Header>
        <Card.Body>
          <Card.Title>{results.name} <small className="text-muted">{results.commonName}</small></Card.Title>
          <dl className="row mt-4">
            <dt className="col-sm-3"><Trans>Registered Directory</Trans></dt>
            <dd className="col-sm-9">{results.registeredDirectory}</dd>

            <dt className="col-sm-3"><Trans>TRISA Member ID</Trans></dt>
            <dd className="col-sm-9">{results.id}</dd>

            <dt className="col-sm-3"><Trans>TRISA Service Endpoint</Trans></dt>
            <dd className="col-sm-9">{results.endpoint}</dd>

            {results.country &&
              <>
              <dt className="col-sm-3"><Trans>Country</Trans></dt>
              <dd className="col-sm-9"><span className="mr-1">{getCountryName(results.country)}</span> {countryCodeEmoji(results.country)}</dd>
              </>
            }

            {results.verifiedOn &&
              <>
              <dt className="col-sm-3"><Trans>TRISA Verification</Trans></dt>
              <dd className="col-sm-9"><Trans>VERIFIED on {results.verifiedOn}</Trans></dd>
              </>
            }

            {results.identityCertificate && results.identityCertificate.signature &&
              <>
              <dt className="col-sm-3"><Trans>TRISA Identity Signature</Trans></dt>
              <dd className="col-sm-9">{results.identityCertificate.signature}</dd>
              </>
            }
          </dl>
        </Card.Body>
      </Card>
    </div>
  );
}

export default LookupResults;