import React from 'react';
import { Accordion } from 'react-bootstrap';

import Detail from './Detail';

function Details({ data }) {
  return (
    <Accordion flush>
      <Accordion.Item eventKey="0" className="border rounded-2 mb-1">
        <Accordion.Header className="m-0 p-0 outline-none">
          <h5 className="text-muted font-weight-bold m-0">Issuer details</h5>
        </Accordion.Header>
        <Accordion.Body style={{ marginBottom: 10, lineHeight: 1.75 }}>
          <Detail data={data.issuer} />
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1" className="border rounded-2 mb-1">
        <Accordion.Header className="m-0 p-0 outline-none">
          <h5 className="text-muted font-weight-bold m-0">Subject details</h5>
        </Accordion.Header>
        <Accordion.Body style={{ marginBottom: 10, lineHeight: 1.75 }}>
          <Detail data={data.subject} />
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

export default Details;
