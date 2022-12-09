import {
  Table,
  TableCaption,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
  Button,
  HStack,
  Stack,
  chakra,
  Tag,
  Tooltip,
  useDisclosure
} from '@chakra-ui/react';
import FormLayout from 'layouts/FormLayout';
import React, { useState } from 'react';
import { Trans, t } from '@lingui/macro';
import { FiMail } from 'react-icons/fi';
import EditCollaboratorModal from './EditCollaborator/EditCollaboratorModal';
import DeleteCollaboratorModal from './DeleteCollaborator/DeleteCollaboratorModal';
import AddCollaboratorModal from 'components/Collaborators/AddCollaborator';
import { useFetchCollaborators } from './useFetchCollaborator';
import type { Collaborator } from './CollaboratorType';
import { formatIsoDate } from 'utils/formate-date';
import { sortCollaboratorsByRecentDate } from './lib';
import Loader from 'components/Loader';
import { useFetchUserRoles } from 'hooks/useFetchUserRoles';
import { COLLABORATOR_STATUS } from 'types/enums';
import { canInviteCollaborator } from 'utils/permission';
import { isDate } from 'utils/date';

const getStatus = (joinedAt: any): any => {
  if (joinedAt && isDate(joinedAt)) {
    return 'Confirmed';
  }
  return 'Pending';
};

const getStatusBgColor = (joinedAt: string) => {
  const status = getStatus(joinedAt) as TCollaboratorStatus;
  switch (status) {
    case COLLABORATOR_STATUS.Confirmed:
      return 'green.400';
    case COLLABORATOR_STATUS.Pending:
      return 'yellow.400';
  }
};

const TableRow: React.FC<{ row: Collaborator }> = ({ row }) => {
  const { roles: userRoles } = useFetchUserRoles();
  return (
    <Tr>
      <Td>
        <chakra.span display="block">{row?.name}</chakra.span>
        <chakra.span display="block" fontSize="sm" color="gray.700">
          {row?.email}
        </chakra.span>
      </Td>
      <Td>{row?.roles}</Td>
      <Td>
        <Tag bg={getStatusBgColor(row?.joined_at as any)} color={'white'} size={'md'}>
          {getStatus(row?.joined_at)}
        </Tag>
      </Td>
      <Td>{formatIsoDate(row?.created_at)}</Td>
      <Td>{formatIsoDate(row?.joined_at)}</Td>
      <Td paddingY={0}>
        <HStack width="100%" justifyContent="center" alignItems="center">
          <Button
            color="blue"
            as={'a'}
            href={`mailto:${row?.email}`}
            bg={'transparent'}
            _hover={{
              bg: 'transparent'
            }}
            _focus={{
              bg: 'transparent'
            }}>
            <FiMail fontSize="24px" />
          </Button>
          <EditCollaboratorModal collaboratorId={row?.id} roles={userRoles?.data} />
          <DeleteCollaboratorModal collaboratorId={row?.id} />
        </HStack>
      </Td>
    </Tr>
  );
};

const TableRows: React.FC = () => {
  const { collaborators } = useFetchCollaborators();
  return (
    <>
      {sortCollaboratorsByRecentDate(collaborators).map((collaborator: Collaborator) => (
        <TableRow key={collaborator.id} row={collaborator} />
      ))}
    </>
  );
};

const CollaboratorsSection: React.FC = () => {
  const { isFetchingCollaborators, hasCollaboratorsFailed } = useFetchCollaborators();

  const { onOpen, isOpen, onClose } = useDisclosure();
  const [isAddCollaboratorModalOpen, setIsAddCollaboratorModalOpen] = useState<boolean>(false);

  if (isFetchingCollaborators) {
    return <Loader />;
  }

  if (hasCollaboratorsFailed) {
    return <Stack>Failed to fetch collaborators</Stack>;
  }

  const modalHandler = () => {
    setIsAddCollaboratorModalOpen(!isAddCollaboratorModalOpen);
    onOpen();
  };
  const closeModalHandler = () => {
    setIsAddCollaboratorModalOpen(!isAddCollaboratorModalOpen);
    onClose();
  };
  return (
    <FormLayout overflowX={'scroll'}>
      <Table variant="simple">
        <TableCaption placement="top" textAlign="end" p={0} m={0} mb={3} fontSize={20}>
          <Tooltip
            label={t`you do not have permission to invite a collaborator`}
            isDisabled={canInviteCollaborator()}>
            <Button minW="170px" onClick={modalHandler} isDisabled={!canInviteCollaborator()}>
              <Trans>Add Contact</Trans>
            </Button>
          </Tooltip>

          {isAddCollaboratorModalOpen && (
            <AddCollaboratorModal
              onOpen={onOpen}
              isOpen={isOpen}
              onClose={onClose}
              onCloseModal={closeModalHandler}
            />
          )}
        </TableCaption>
        <Thead>
          <Tr>
            <Th>
              <Trans>Contact Details</Trans>
            </Th>
            <Th>
              <Trans>Role</Trans>
            </Th>
            <Th>
              <Trans>Status</Trans>
            </Th>
            <Th>
              <Trans>Invited</Trans>
            </Th>
            <Th>
              <Trans>Joined</Trans>
            </Th>
            <Th textAlign="center">
              <Trans>Actions</Trans>
            </Th>
          </Tr>
        </Thead>
        <Tbody>
          <TableRows />
        </Tbody>
      </Table>
    </FormLayout>
  );
};
export default CollaboratorsSection;
