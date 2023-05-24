import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
  Button,
  ModalHeader
} from '@chakra-ui/react';

type InvalidFormPromptProps = {
  isOpen: boolean;
  onClose: () => void;
  handleContinueClick: () => void;
};

function InvalidFormPrompt({ isOpen, onClose, handleContinueClick }: InvalidFormPromptProps) {
  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>&nbsp;</ModalHeader>
        <ModalCloseButton />
        <ModalBody>You are about to lose your changes</ModalBody>
        <ModalFooter>
          <Button variant="ghost" mr={3} onClick={onClose}>
            Cancel
          </Button>
          <Button colorScheme="blue" onClick={handleContinueClick}>
            Continue
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}

export default InvalidFormPrompt;
