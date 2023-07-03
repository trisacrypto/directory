import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
  Button,
  VStack,
  ModalHeader,
  Text,
  Box
} from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
type InvalidFormPromptProps = {
  isOpen: boolean;
  onClose: () => void;
  handleContinueClick: () => void;
  nextStepBtnContent?: string;
};

function InvalidFormPrompt({ isOpen, onClose, handleContinueClick }: InvalidFormPromptProps) {
  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader textAlign={'center'}>
          <Trans>Unsaved changes alert</Trans>
        </ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <VStack>
            <Text fontWeight="semibold">
              <Trans>
                You have unsaved changes. Are you sure you want to continue without saving?
              </Trans>
            </Text>
            <Box>
              {' '}
              <Text>
                <Trans>
                  If you continue, your changes will be lost. To save your changes, click Cancel and
                  then click on the
                </Trans>
              </Text>
              <Text as="span" fontWeight="bold">
                <Trans>"Save & Next" or "Save & Previous" button.</Trans>
              </Text>
            </Box>
          </VStack>
        </ModalBody>
        <ModalFooter>
          <Button variant="ghost" mr={3} onClick={onClose} data-cy="cancel-bttn">
            Cancel
          </Button>
          <Button colorScheme="blue" onClick={handleContinueClick} data-cy="continue-bttn">
            Continue
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}

export default InvalidFormPrompt;
