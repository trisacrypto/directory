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
import { useForm } from 'react-hook-form';

function ConfirmIdentityCertificate() {
  const { onClose, onOpen, isOpen } = useDisclosure();
  const { register } = useForm();

  return (
    <>
      <Button bg="#55ACD8" color="#fff" onClick={onOpen}>
        Request New Identity Certificate
      </Button>
      <form>
        <Modal isOpen={isOpen} onClose={onClose}>
          <ModalOverlay />
          <ModalContent border="1px solid">
            <ModalHeader mt={3} pb={1}>
              New X.509 Identity Certificate Request
            </ModalHeader>
            <ModalBody display="flex" flexDirection="column" gap={[4, 6]}>
              <Text>
                Requesting a new X.509 Identity Certificate will invalidate and revoke your current
                X.509 Identity Certificate.
              </Text>
              <Stack>
                <Checkbox
                  {...register('accept', {
                    required: true
                  })}>
                  I acknowledge that requesting a new X.509 Identity Certificate will invalidate and
                  revoke my organization’s current X.509 Identity Certificate.
                </Checkbox>
              </Stack>
              <Text>
                You are required to re-confirm your organization’s profile with TRISA. Click next to
                proceed. You can cancel later.
              </Text>
            </ModalBody>

            <ModalFooter color="#fff" justifyContent="space-evenly">
              <Button bg="#55ACD8" type="submit">
                Next
              </Button>
              <Button bg="#555151D4" _hover={{ boxShadow: "#555151D4" }} onClick={onClose}>
                Cancel
              </Button>
            </ModalFooter>
          </ModalContent>
        </Modal>
      </form>
    </>
  );
}

export default ConfirmIdentityCertificate;
