import {
  Box,
  Button,
  chakra,
  Link,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
  useDisclosure,
  VStack
} from '@chakra-ui/react';
import { Trans } from '@lingui/macro';
import SelectFormControl from 'components/ui/SelectFormControl';
import { FiEdit } from 'react-icons/fi';

function EditCollaboratorModal() {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <Link color="blue" onClick={onOpen}>
      <FiEdit fontSize="24px" />
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader textTransform="capitalize" textAlign="center" fontWeight={700}>
            <Trans>Edit collaborator Role</Trans>
          </ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Text>
              <Trans>Select the collaborator’s role and save.</Trans>
            </Text>
            <VStack align="start" spacing={4} mt={2}>
              <Box>
                <Text fontWeight={700}>
                  <Trans>Collaborator Name & Email</Trans>
                </Text>
                <Text textTransform="capitalize">Eason Yang</Text>
                <Text>eyang@vaspnet.co.uk</Text>
              </Box>
              <SelectFormControl
                label={
                  <>
                    <chakra.span fontWeight={700}>
                      <Trans>Change Role</Trans>
                    </chakra.span>
                  </>
                }
                controlId="role"
              />
            </VStack>
          </ModalBody>

          <ModalFooter display="flex" flexDir="column" gap={3}>
            <Button bg="orange" minW="150px" _hover={{ bg: 'orange' }} onClick={onClose}>
              <Trans>Save</Trans>
            </Button>
            <Button variant="ghost" minW="150px" color="link">
              <Trans>Close</Trans>
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Link>
  );
}

export default EditCollaboratorModal;
