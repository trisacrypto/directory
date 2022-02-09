import React from 'react';
import { Button } from 'react-bootstrap';
import ButtonSpinner from './ButtonSpinner';
import warningImg from 'assets/images/warning.svg'
import { ModalCloseButton } from './Modal';
import PropTypes from 'prop-types'
import { useSelector } from 'react-redux';
import { getVaspDetails } from 'redux/selectors';

function DeleteModal({ onDelete, isLoading = false }) {
    const vasp = useSelector(getVaspDetails)

    return <div className='text-center'>
        <img src={warningImg} alt="Warning" />
        <h5 className='fw-normal' style={{ lineHeight: 1.5 }}>Are you sure you want to delete registration <span className='fw-bold'>{vasp?.name}</span>?</h5>
        <p>This action cannot be undone.</p>
        <div className='d-flex justify-content-around mt-3'>
            <ModalCloseButton>
                <Button variant="outline-primary">Cancel</Button>
            </ModalCloseButton>
            <ButtonSpinner
                isLoading={isLoading}
                label='Delete'
                loadingMessage="Deleting..."
                onClick={onDelete}
                variant='danger'
            />
        </div>
    </div>;
}

DeleteModal.propTypes = {
    onDelete: PropTypes.func,
    isLoading: PropTypes.bool
}

export default DeleteModal;
