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
  isNextStep: boolean;
};

function InvalidFormPrompt({
  isOpen,
  onClose,
  handleContinueClick,
  isNextStep
}: InvalidFormPromptProps) {
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
              <Text>
                <Trans>
                  If you continue, your changes will be lost. To save your changes, click Cancel and
                  then click on the{' '}
                {isNextStep ? <Text as="span" fontWeight={700} whiteSpace={'break-spaces'}>Save & Next</Text> : <Text as="span" fontWeight={700} whiteSpace={'break-spaces'}>Save & Previous</Text>} button.
                </Trans>
              </Text>
            </Box>
          </VStack>
        </ModalBody>
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
