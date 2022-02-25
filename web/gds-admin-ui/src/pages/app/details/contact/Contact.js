
import classNames from 'classnames';
import React from 'react'
import PropTypes from 'prop-types';
import { Button, Col, Row } from 'react-bootstrap';
import { formatPhoneNumberIntl, isValidPhoneNumber } from 'react-phone-number-input';
import { formatDisplayedData, verifiedContactStatus } from 'utils';
import { VERIFIED_CONTACT_STATUS, VERIFIED_CONTACT_STATUS_LABEL } from 'constants/index';
import DeleteContactModal from './DeleteContactModal';
import { Modal, ModalContent, ModalOpenButton } from 'components/Modal';
import EditContactModal from './EditContactModal';


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
            <p className="fw-bold mb-1 mt-2"> <span className='text-capitalize'>{type}</span> contact
                <Modal>
                    <ModalOpenButton>
                        {
                            data ?
                                <Button variant='light' className='btn-circle ms-1' title='Edit'><i className='mdi mdi-square-edit-outline text-warning' /></Button> :
                                <Button variant='light' className='btn-circle ms-1' title='Edit'><i className="mdi mdi-plus-box-outline text-success" /></Button>
                        }
                    </ModalOpenButton>
                    <ModalContent size="lg">
                        <Row className='p-4'>
                            <Col xs={12}>
                                <EditContactModal contactType={type} />
                            </Col>
                        </Row>
                    </ModalContent>
                </Modal >
                <Modal>
                    <DeleteContactModal type={type} />
                </Modal>
            </p>
            <hr className='my-1' />
            <Row>
                <ul className='list-unstyled'>
                    {data?.name ? <li>{formatDisplayedData(data?.name)}</li> : null}
                    {data?.phone && isValidPhoneNumber(data.phone) ? <li>{formatDisplayedData(formatPhoneNumberIntl(data.phone))}</li> : null}
                    {data?.email ? <li>{formatDisplayedData(data?.email)} <small data-testid="verifiedContactStatus" style={{ fontStyle: 'italic' }}>{VERIFIED_CONTACT_STATUS_LABEL[status]}</small> </li> : null}
                    <li><small style={{ fontStyle: 'italic' }}>{data?.person ? 'Has IVMS101 Record' : 'No IVMS101 Data'}</small></li>
                </ul>
            </Row>
        </div>
    )
}

Contact.propTypes = {
    type: PropTypes.oneOf(['technical', 'administrative', 'billing', 'legal']).isRequired,
    verifiedContact: PropTypes.objectOf(PropTypes.string).isRequired,
    data: PropTypes.object
}

export default Contact
