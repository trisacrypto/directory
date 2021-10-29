import React from 'react';
import { Card } from 'react-bootstrap';
import Contact from './Contact';

export default function ContactList({ data, verifiedContact }) {

    return (
        <Card>
            <Card.Body>
                <h4 className="mt-0 mb-3">Contacts</h4>
                {data?.technical ? <Contact data={data?.technical} verifiedContact={verifiedContact} type="Technical" /> : null}
                {data?.administrative ? <Contact data={data?.administrative} verifiedContact={verifiedContact} type="Administrative" /> : null}
                {data?.legal ? <Contact data={data?.legal} verifiedContact={verifiedContact} type="Legal" /> : null}
                {data?.billing ? <Contact data={data?.billing} verifiedContact={verifiedContact} type="Billing" /> : null}
            </Card.Body>
        </Card>
    );
};
