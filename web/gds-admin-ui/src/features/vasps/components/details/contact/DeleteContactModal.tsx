import { Button, Col, Row } from 'react-bootstrap';

import { ModalContent, ModalOpenButton } from '@/components/Modal';

import DeleteContactPromptModal from './DeleteContactPromptModal';
import useDeleteContactModal from '../../../services/use-delete-contact-modal';

type DeleteContactPromptModalProps = {
    type: string;
};

function DeleteContactModal({ type }: DeleteContactPromptModalProps) {
    const { handleDeleteClick } = useDeleteContactModal({ contactType: type });

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

export default DeleteContactModal;
