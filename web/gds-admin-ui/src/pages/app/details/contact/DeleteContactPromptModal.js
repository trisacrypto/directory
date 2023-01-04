import PropTypes from 'prop-types';
import { Button, Col, Row } from 'react-bootstrap';
import { useSelector } from 'react-redux';

import { ModalCloseButton } from '@/components/Modal';
import { getVaspDetailsLoadingState } from '@/redux/selectors';

function DeleteContactPromptModal({ type, onDelete }) {
  const isLoading = useSelector(getVaspDetailsLoadingState);

  return (
    <>
      <p className="text-center">
        Are you sure you want to delete the <span className="fw-bold">{type}</span> contact?
      </p>
      <p className="text-center">This action cannot be undone.</p>
      <Row>
        <Col className="text-center">
          <ModalCloseButton>
            <Button variant="outline-primary">Cancel</Button>
          </ModalCloseButton>
        </Col>
        <Col className="text-center">
          <Button onClick={onDelete} variant="danger" disabled={isLoading}>
            Delete
          </Button>
        </Col>
      </Row>
    </>
  );
}

DeleteContactPromptModal.propTypes = {
  type: PropTypes.oneOf(['administrative', 'legal', 'billing', 'technical']),
  onDelete: PropTypes.func.isRequired,
};

export default DeleteContactPromptModal;
