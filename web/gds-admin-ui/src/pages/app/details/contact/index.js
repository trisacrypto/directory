import { Card } from 'react-bootstrap';

import Contact from './Contact';

export default function ContactList({ data, verifiedContact }) {
  return (
    <Card>
      <Card.Body>
        <h4 className="mt-0 mb-3 text-dark">Contacts</h4>
        <div className="d-flex justify-content-between">
          <div className="">
            <i className="mdi mdi-alert fs-5 text-warning me-1" />
            <small className="fst-italic">Alternate verified</small>
          </div>
          <div>
            <i className="mdi mdi-check-all fs-5 text-success me-1" />
            <small className="fst-italic">Verified</small>
          </div>
          <div>
            <i className="mdi mdi-close-circle fs-5 text-danger me-1" />
            <small className="fst-italic">Unverified</small>
          </div>
        </div>
        {Object.entries(data || {}).map(([k]) => (
          <Contact key={k} data={data[k]} verifiedContact={verifiedContact} type={k} />
        ))}
      </Card.Body>
    </Card>
  );
}
