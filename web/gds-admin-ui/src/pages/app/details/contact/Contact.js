
import classNames from 'classnames';
import React from 'react'
import { Row } from 'react-bootstrap';
import { VERIFIED_CONTACT_STATUS, VERIFIED_CONTACT_STATUS_LABEL } from '../../../../constants';
import { formatDisplayedData, verifiedContactStatus } from '../../../../utils';
import PropTypes from 'prop-types';

function Contact({ data, type, verifiedContact }) {
    const status = verifiedContactStatus({ data, type, verifiedContact })

    const getColor = (state) => {
        if (state === VERIFIED_CONTACT_STATUS.ALTERNATE_VERIFIED) {
            return 'text-warning'
        }
        if (state === VERIFIED_CONTACT_STATUS.UNVERIFIED) {
            return 'text-danger'
        }
        return undefined
    }

    return (
        <div data-testid="contact-node" className={classNames(getColor(status))}>
            <p className="fw-bold mb-1 mt-2">{type} contact:</p>
            <hr className='my-1' />
            <Row>
                <ul className='list-unstyled'>
                    {data?.name ? <li>{formatDisplayedData(data?.name)}</li> : null}
                    {data?.phone ? <li>{formatDisplayedData(data?.phone)}</li> : null}
                    {data?.email ? <li>{formatDisplayedData(data?.email)} <small data-testid="verifiedContactStatus" style={{ fontStyle: 'italic' }}>{VERIFIED_CONTACT_STATUS_LABEL[status]}</small> </li> : null}
                    <li><small style={{ fontStyle: 'italic' }}>{data?.person ? 'Has IVMS101 Record' : 'No IVMS101 Data'}</small></li>
                </ul>
            </Row>
        </div>
    )
}

Contact.propTypes = {
    type: PropTypes.oneOf(['Technical', 'Administrative', 'Billing', 'Legal']).isRequired,
    verifiedContact: PropTypes.objectOf(PropTypes.string).isRequired,
    data: PropTypes.object.isRequired
}

export default Contact
