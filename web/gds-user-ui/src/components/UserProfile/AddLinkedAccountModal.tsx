import {
  Button,
  chakra,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  useDisclosure
} from '@chakra-ui/react';
import CkLazyLoadImage from 'components/LazyImage';
import AddIcon from 'assets/carbon_add-alt.svg';
import { Trans } from '@lingui/macro';

function AddLinkedAccountModal() {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <>
      <Button variant="unstyled" display="flex" justifyContent="center" onClick={onOpen}>
        <chakra.span color="blue">
          <Trans>Add</Trans>
        </chakra.span>
        <CkLazyLoadImage src={AddIcon} width="25px" ml="3px" />
      </Button>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent border="1px solid black" px={5}>
          <ModalHeader textAlign="center" fontWeight={700} fontSize="md">
            <Trans>Add Linked Account</Trans>
          </ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Trans>
              If you have additional accounts with the TRISA Global Directory Service, you can link
              them here. You are required to log in to the linked account to verify account
              ownership.
            </Trans>
          </ModalBody>

          <ModalFooter display="flex" flexDir="column" rowGap={2}>
            <Button bg="orange" _hover={{ bg: 'orange' }}>
              <Trans>Link Account</Trans>
            </Button>
            <Button variant="ghost" onClick={onClose}>
              <Trans>Cancel</Trans>
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
}

export default AddLinkedAccountModal;
