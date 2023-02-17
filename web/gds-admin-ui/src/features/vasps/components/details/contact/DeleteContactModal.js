import PropTypes from 'prop-types';
import React from 'react';
import { Button, Col, Row } from 'react-bootstrap';
import { useDispatch } from 'react-redux';
import { useParams } from 'react-router-dom';

import { ModalContent, ModalContext, ModalOpenButton } from '@/components/Modal';
import useSafeDispatch from '@/hooks/useSafeDispatch';
import { deleteContactResponse } from '@/redux/vasp-details';

import DeleteContactPromptModal from './DeleteContactPromptModal';

function DeleteContactModal({ type }) {
  const dispatch = useDispatch();
  const safeDispatch = useSafeDispatch(dispatch);
  const params = useParams();
  const [, setIsOpen] = React.useContext(ModalContext);

  const handleDeleteClick = () => {
    if (params && params.id) {
      safeDispatch(deleteContactResponse(params.id, type, setIsOpen));
    }
  };

  return (
    <>
      <ModalOpenButton>
        <Button variant="light" className="btn-circle ms-1" title="Delete">
          <i className=" mdi mdi-delete-circle text-danger" />
        </Button>
      </ModalOpenButton>
      <ModalContent size="sm">
        <Row className="p-4">
          <Col xs={12}>
            <DeleteContactPromptModal onDelete={handleDeleteClick} type={type} />
          </Col>
        </Row>
      </ModalContent>
    </>
  );
}

DeleteContactPromptModal.propTypes = {
  type: PropTypes.oneOf(['administrative', 'legal', 'billing', 'technical']),
};

export default DeleteContactModal;
