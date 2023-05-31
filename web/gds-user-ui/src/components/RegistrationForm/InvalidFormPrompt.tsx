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
import { useEffect, useState } from 'react';
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
  const [btnContent, setIsBtnContent] = useState("");

  useEffect(() => {
    if (isNextStep) {
      setIsBtnContent("Save & Next");
    } else {
      setIsBtnContent("Save & Previous");
    }
  }, [isNextStep]);

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
                  then click on the
                </Trans>
              </Text>
              <Text fontWeight={'bold'} as="span">
                <Trans>{btnContent} button.</Trans>
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
