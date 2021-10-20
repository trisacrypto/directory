
import React from 'react'
import { Row } from 'react-bootstrap';
import { formatDisplayedData } from '../../../../utils';

function Contact({ data, type }) {
    return (
        <>
            <p className="fw-bold mb-1 mt-2">{type} contact:</p>
            <hr className='my-1' />
            <Row>
                <p className="">{formatDisplayedData(data?.name)}</p>
                <p className="">{formatDisplayedData(data?.phone)}</p>
                <p className="">{formatDisplayedData(data?.email)}</p>
                <p className="">{data?.person ? 'Has IVMS101 Record' : 'No IVMS101 Data'}</p>
            </Row>
        </>
    )
}

export default Contact
