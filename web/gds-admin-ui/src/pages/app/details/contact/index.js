import { Card } from 'react-bootstrap';
import Contact from './Contact';

export default function ContactList({ data, verifiedContact }) {

    return (
        <Card>
            <Card.Body>
                <h4 className="mt-0 mb-3 text-dark">Contacts</h4>
                {
                    Object.entries(data || {}).map(([k]) => (
                        <Contact key={k} data={data[k]} verifiedContact={verifiedContact} type={k} />
                    ))
                }
            </Card.Body>
        </Card >
    );
};
