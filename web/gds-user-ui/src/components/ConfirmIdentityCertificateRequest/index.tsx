import {
  Button,
  Checkbox,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Stack,
  Text,
  useDisclosure
} from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { useForm } from 'react-hook-form';

function ConfirmIdentityCertificate() {
  const { onClose, onOpen, isOpen } = useDisclosure();
  const { register } = useForm();

  return (
    <>
      <Button bg="#55ACD8" color="#fff" onClick={onOpen}>
        <Trans id="Request New Identity Certificate">Request New Identity Certificate</Trans>
      </Button>
      <form>
        <Modal isOpen={isOpen} onClose={onClose}>
          <ModalOverlay />
          <ModalContent border="1px solid">
            <ModalHeader mt={3} pb={1}>
              <Trans id="New X.509 Identity Certificate Request">
                New X.509 Identity Certificate Request
              </Trans>
            </ModalHeader>
            <ModalBody display="flex" flexDirection="column" gap={[4, 6]}>
              <Text>
                <Trans id="Requesting a new X.509 Identity Certificate will invalidate and revoke your current X.509 Identity Certificate.">
                  Requesting a new X.509 Identity Certificate will invalidate and revoke your
                  current X.509 Identity Certificate.
                </Trans>
              </Text>
              <Stack>
                <Checkbox
                  {...register('accept', {
                    required: true
                  })}>
                  <Trans id="I acknowledge that requesting a new X.509 Identity Certificate will invalidate and revoke my organization’s current X.509 Identity Certificate.">
                    I acknowledge that requesting a new X.509 Identity Certificate will invalidate
                    and revoke my organization’s current X.509 Identity Certificate.
                  </Trans>
                </Checkbox>
              </Stack>
              <Text>
                You are required to re-confirm your organization’s profile with TRISA. Click next to
                proceed. You can cancel later.
              </Text>
            </ModalBody>

            <ModalFooter color="#fff" justifyContent="space-evenly">
              <Button bg="#55ACD8" type="submit">
                <Trans id="Next">Next</Trans>
              </Button>
              <Button bg="#555151D4" _hover={{ boxShadow: '#555151D4' }} onClick={onClose}>
                <Trans id="Cancel">Cancel</Trans>
              </Button>
            </ModalFooter>
          </ModalContent>
        </Modal>
      </form>
    </>
  );
}

export default ConfirmIdentityCertificate;
