import React from 'react';
import { Card } from 'react-bootstrap';
import Contact from './Contact';

export default function ContactList({ data }) {

    return (
        <Card>
            <Card.Body>
                <h4 className="mt-0 mb-3">Contacts</h4>
                {data?.technical ? <Contact data={data?.technical} type="Technical" /> : null}
                {data?.administrative ? <Contact data={data?.administrative} type="Administrative" /> : null}
                {data?.legal ? <Contact data={data?.legal} type="Legal" /> : null}
                {data?.billing ? <Contact data={data?.billing} type="Billing" /> : null}
            </Card.Body>
        </Card>
    );
};
