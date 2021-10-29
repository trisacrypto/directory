
import classNames from 'classnames';
import React from 'react'
import { Row } from 'react-bootstrap';
import { formatDisplayedData, isVerifiedContact } from '../../../../utils';

function Contact({ data, type, verifiedContact }) {
    const isVerified = isVerifiedContact(data, verifiedContact);

    return (
        <div className={classNames({ 'text-danger': !isVerified })}>
            <p className="fw-bold mb-1 mt-2">{type} contact:</p>
            <hr className='my-1' />
            <Row>
                <ul className='list-unstyled'>
                    {data?.name ? <li className="">{formatDisplayedData(data?.name)}</li> : null}
                    {data?.phone ? <li className="">{formatDisplayedData(data?.phone)}</li> : null}
                    {data?.email ? <li className="">{formatDisplayedData(data?.email)} {isVerified ? <small style={{ fontStyle: 'italic' }}>verified</small> : null}</li> : null}
                    <li><small style={{ fontStyle: 'italic' }}>{data?.person ? 'Has IVMS101 Record' : 'No IVMS101 Data'}</small></li>
                </ul>
            </Row>
        </div>
    )
}

export default Contact
