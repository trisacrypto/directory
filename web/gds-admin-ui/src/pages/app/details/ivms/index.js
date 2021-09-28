
import React from 'react'
import { Card, Row } from 'react-bootstrap';
import { formatDisplayedData } from '../../../../utils';
import Geographic from './Geographic';
import Name from './Name';
import NationalIdentification from './NationalIdentification';

function Ivms({ data }) {

    return (
        <Card>
            <Card.Body>
                <h4 className="mt-0 mb-3">IVMS 1010</h4>
                <p className="fw-bold mb-2">Country Of Registration: <span className="fw-normal">{formatDisplayedData(data?.country_of_registration)}</span></p>
                <p className="fw-bold mb-2">Customer Number: <span className="fw-normal">{formatDisplayedData(data?.customer_number)}</span></p>
                <Row>
                    <Geographic data={data?.geographic_addresses} />
                    <NationalIdentification data={data?.national_identification} />
                    <Name data={data?.name} />
                </Row>

            </Card.Body>
        </Card>
    )
}

export default Ivms
