/* eslint-disable @typescript-eslint/no-unused-vars */
import React, { useState, useEffect } from 'react';
import {
  Link,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react';
import type { NewCollaborator } from './AddCollaboratorType';
import type { Collaborator } from 'components/Collaborators/CollaboratorType';
import AddCollaboratorForm from './AddCollaboratorForm';

interface Props {
  isOpen: boolean;
  onClose: () => void;
  onOpen: () => void;
  onCloseModal: () => void;
}

function AddCollaboratorModal(props: Props) {
  const { isOpen, onOpen, onClose, onCloseModal } = props;


  return (
    <Link color="blue" onClick={onOpen}>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent w="100%" maxW="600px" px={10}>
          <ModalHeader textTransform="capitalize" textAlign="center" fontWeight={700} pb={1}>
            Add New Contact
          </ModalHeader>

          <ModalCloseButton onClick={onCloseModal} />

          <ModalBody>
            <AddCollaboratorForm onCloseModal={onCloseModal} />
          </ModalBody>
        </ModalContent>
      </Modal>
    </Link>
  );
}

export default AddCollaboratorModal;
