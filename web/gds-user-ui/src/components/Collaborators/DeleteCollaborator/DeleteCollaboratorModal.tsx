import {
  Box,
  Button,
  Checkbox,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
  useDisclosure,
  VStack,
  useToast
} from '@chakra-ui/react';
import { useFetchCollaborators } from 'components/Collaborators/useFetchCollaborator';
import { useDeleteCollaborator } from 'components/Collaborators/DeleteCollaborator/useDeleteCollaborator';
import { Trans } from '@lingui/react';
import { BsTrash } from 'react-icons/bs';
import { useEffect, useState } from 'react';
import type { Collaborator } from 'components/Collaborators/CollaboratorType';
import { t } from '@lingui/macro';
import { USER_PERMISSION } from 'types/enums';
import { useSafeDisableIconButton } from 'components/Collaborators/useSafeDisableIconButton';
interface Props {
  collaboratorId: string;
}
function DeleteCollaboratorModal(props: Props) {
  const { collaboratorId } = props;
  const { isOpen, onOpen, onClose } = useDisclosure();
  const toast = useToast();
  const { collaborators, getAllCollaborators } = useFetchCollaborators();
  const {
    isDeleting,
    wasCollaboratorDeleted,
    deleteCollaborator,
    hasCollaboratorFailed,
    errorMessage
  } = useDeleteCollaborator();

  const [collaborator, setCollaborator] = useState<Collaborator>();

  const [isDeleteChecked, setIsDeleteChecked] = useState(false);

  const { isDisabled: isNotCurrentUserAndHasPermission } = useSafeDisableIconButton(
    USER_PERMISSION.UPDATE_COLLABORATOR,
    collaborator?.email as string
  );

  const deleteHandler = () => {
    // delete collaborator
    deleteCollaborator(collaboratorId);

    setIsDeleteChecked(false);
    onClose();
  };

  useEffect(() => {
    if (wasCollaboratorDeleted) {
      getAllCollaborators();
      // display success toast
      toast({
        title: 'Collaborator deleted',
        description: 'The collaborator has been deleted',
        status: 'success',
        duration: 9000,
        isClosable: true,
        position: 'top-right'
      });
    }
  }, [wasCollaboratorDeleted, getAllCollaborators, toast]);

  useEffect(() => {
    const col = collaborators?.find((c: Collaborator) => c.id === collaboratorId);
    if (col) {
      setCollaborator(col);
    }
  }, [collaboratorId, collaborators]);

  useEffect(() => {
    if (hasCollaboratorFailed && !wasCollaboratorDeleted) {
      // display error toast
      const hasErrored =
        errorMessage &&
        t`An error occurred while deleting the collaborator, please try again or contact support at support@rotational.io`;
      toast({
        title: t`Collaborator not deleted`,
        description: hasErrored || t`The collaborator has not been deleted`,
        status: 'error',
        duration: 9000,
        isClosable: true,
        position: 'top-right'
      });
    }
  }, [hasCollaboratorFailed, toast, wasCollaboratorDeleted, errorMessage]);

  return (
    <Button
      color="blue"
      onClick={onOpen}
      bg={'transparent'}
      data-testid="icon-collaborator-button"
      disabled={!isNotCurrentUserAndHasPermission}
      _hover={{
        bg: 'transparent'
      }}
      _focus={{
        bg: 'transparent'
      }}>
      <BsTrash fontSize="26px" />
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader textTransform="capitalize" textAlign="center" fontWeight={700}>
            <Trans id="Delete Collaborator">Delete Collaborator</Trans>
          </ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Text>
              <Trans id="Delete the following collaborator from the VASP account. The collaborator will no longer have access to the account.">
                Delete the following collaborator from the VASP account. The collaborator will no
                longer have access to the account.
              </Trans>
            </Text>
            <VStack align="start" spacing={4} mt={2}>
              <Box>
                <Text fontWeight={700}>
                  <Trans id="Collaborator Name & Email">Collaborator Name & Email</Trans>
                </Text>
                <Text textTransform="capitalize" data-testid="collaborator-name">
                  {collaborator?.name}
                </Text>
                <Text data-testid="collaborator-email">{collaborator?.email}</Text>
              </Box>
              <Checkbox
                borderColor={'black'}
                isChecked={isDeleteChecked}
                colorScheme="gray"
                onChange={(e) => setIsDeleteChecked(e.target.checked)}>
                <Trans id="Check to delete user account.">Check to delete user account.</Trans>
              </Checkbox>
            </VStack>
          </ModalBody>

          <ModalFooter display="flex" flexDir="column" gap={3}>
            <Button
              bg="orange"
              minW="150px"
              data-testid="delete-collaborator-button"
              _hover={{ bg: 'orange' }}
              onClick={deleteHandler}
              isDisabled={!isDeleteChecked}
              isLoading={isDeleting}>
              <Trans id="Delete">Delete</Trans>
            </Button>
            <Button variant="ghost" minW="150px" color="link" onClick={onClose}>
              <Trans id="Close">Close</Trans>
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Button>
  );
}

export default DeleteCollaboratorModal;
