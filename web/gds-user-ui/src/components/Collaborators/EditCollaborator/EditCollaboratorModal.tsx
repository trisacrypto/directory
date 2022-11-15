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
  VStack,
  useToast
} from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import SelectFormControl from 'components/ui/SelectFormControl';
import { FiEdit } from 'react-icons/fi';
import type { Collaborator } from 'components/Collaborators/CollaboratorType';
import { useUpdateCollaborator } from 'components/Collaborators/EditCollaborator/useUpdateCollaborator';
import { useFetchCollaborators } from 'components/Collaborators/useFetchCollaborator';
import React, { useEffect, useState } from 'react';
import { t } from '@lingui/macro';

interface Props {
  collaboratorId: string;
  roles?: string[];
}

function EditCollaboratorModal(props: Props) {
  const { collaboratorId, roles } = props;
  const { isOpen, onOpen, onClose } = useDisclosure();
  const toast = useToast();
  const { collaborators, getAllCollaborators } = useFetchCollaborators();
  const [collaborator, setCollaborator] = useState<Collaborator>();
  const [selectedRole, setSelectedRole] = useState('');

  const {
    isUpdating,
    wasCollaboratorUpdated,
    updateCollaborator,
    hasCollaboratorFailed,
    errorMessage
  } = useUpdateCollaborator();

  const updateHandler = () => {
    // update user role
    const collaboratorData = {
      data: {
        roles: new Array(selectedRole)
      },
      id: collaboratorId
    };
    updateCollaborator(collaboratorData);

    onClose();
  };

  // roles options from userRoles {label: string , value: string}
  const rolesOptions = roles?.map((v: string) => ({
    label: v,
    value: v
  }));

  console.log('[roles options]', rolesOptions);

  useEffect(() => {
    let once = false;
    const col = collaborators?.data?.collaborators.find(
      (c: Collaborator) => c.id === collaboratorId
    );
    if (col) {
      if (!once) {
        setCollaborator(col);
      }
    }
    return () => {
      once = true;
    };
  }, [collaboratorId, collaborators]);

  useEffect(() => {
    if (wasCollaboratorUpdated) {
      getAllCollaborators();
      // display success toast
      toast({
        title: 'Collaborator updated',
        description: 'The collaborator has been updated',
        status: 'success',
        duration: 9000,
        isClosable: true,
        position: 'top-right'
      });
    }
  }, [wasCollaboratorUpdated, getAllCollaborators, toast]);

  useEffect(() => {
    if (hasCollaboratorFailed && !wasCollaboratorUpdated) {
      const hasErrored =
        errorMessage &&
        t`An error occurred while updating the collaborator, please try again or contact support at support@trisa.io`;
      toast({
        title: t`Collaborator is not updated`,
        description: hasErrored || t`The collaborator has not been updated`,
        status: 'error',
        duration: 9000,
        isClosable: true,
        position: 'top-right'
      });
    }
  }, [hasCollaboratorFailed, wasCollaboratorUpdated, errorMessage, toast]);

  return (
    <Link color="blue" onClick={onOpen}>
      <FiEdit fontSize="24px" />
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader textTransform="capitalize" textAlign="center" fontWeight={700}>
            <Trans id="Edit collaborator Role">Edit collaborator Role</Trans>
          </ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Text>
              <Trans id="Select the collaborator’s role and save.">
                Select the collaborator’s role and save.
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
              <SelectFormControl
                label={
                  <>
                    <chakra.span fontWeight={700}>
                      <Trans id="Change Role">Change Role</Trans>
                    </chakra.span>
                  </>
                }
                options={rolesOptions}
                controlId="role"
                onChange={(s: any) => {
                  console.log('[selected role]', s.value);
                  setSelectedRole(s.value);
                }}
                value={rolesOptions?.find((v) => v.value === selectedRole)}
                name="role"
                placeholder="Select Role"
              />
            </VStack>
          </ModalBody>

          <ModalFooter display="flex" flexDir="column" gap={3}>
            <Button
              bg="orange"
              minW="150px"
              _hover={{ bg: 'orange' }}
              data-testid="update-collaborator-button"
              onClick={updateHandler}
              isDisabled={!selectedRole}
              isLoading={isUpdating}>
              <Trans id="Save">Save</Trans>
            </Button>
            <Button
              variant="ghost"
              minW="150px"
              color="link"
              onClick={onClose}
              isLoading={isUpdating}>
              <Trans id="Close">Close</Trans>
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Link>
  );
}

export default EditCollaboratorModal;
