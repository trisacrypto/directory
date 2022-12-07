import {
  IconButton,
  Modal,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  useDisclosure
} from '@chakra-ui/react';
import CkLazyLoadImage from 'components/LazyImage';
import EditIcon from 'assets/edit-input.svg';
import ChangeNameForm from './ChangeNameForm';

function ChangeNameModal() {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <>
      <IconButton
        aria-label="Edit"
        icon={<CkLazyLoadImage src={EditIcon} mx="auto" w="25px" />}
        variant="unstyled"
        marginTop="32px!important"
        onClick={onOpen}
      />
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent border="1px solid black" px={5}>
          <ModalHeader textAlign="center" fontWeight={700} fontSize="md">
            Change Name
          </ModalHeader>
          <ModalCloseButton />
          <ChangeNameForm onCloseModal={onClose} />
        </ModalContent>
      </Modal>
    </>
  );
}

export default ChangeNameModal;
