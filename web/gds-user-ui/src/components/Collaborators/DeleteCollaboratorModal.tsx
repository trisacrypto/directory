import {
  Box,
  Button,
  Checkbox,
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
import { BsTrash } from 'react-icons/bs';

function DeleteCollaboratorModal() {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <Link color="blue" onClick={onOpen}>
      <BsTrash fontSize="26px" />
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader textTransform="capitalize" textAlign="center" fontWeight={700}>
            <Trans>Delete Collaborator</Trans>
          </ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Text>
              <Trans>
                Delete the following collaborator from the VASP account. The collaborator will no
                longer have access to the account.
              </Trans>
            </Text>
            <VStack align="start" spacing={4} mt={2}>
              <Box>
                <Text fontWeight={700}>
                  <Trans>Collaborator Name & Email</Trans>
                </Text>
                <Text textTransform="capitalize">Eason Yang</Text>
                <Text>eyang@vaspnet.co.uk</Text>
              </Box>
              <Checkbox defaultChecked>
                <Trans>Check to delete user account.</Trans>
              </Checkbox>
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

export default DeleteCollaboratorModal;
