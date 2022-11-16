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
  Link,
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
// import { getCollaborators, setCollaborators } from 'application/store/selectors/collaborator';
import { useFetchCollaborators } from './useFetchCollaborator';
// import { useDispatch } from 'react-redux';
import type { Collaborator } from './CollaboratorType';
import { formatIsoDate } from 'utils/formate-date';
import { sortCollaboratorsByRecentDate } from './lib';
import Loader from 'components/Loader';
import { useFetchUserRoles } from 'hooks/useFetchUserRoles';
import { USER_PERMISSION } from 'types/enums';
import Store from 'application/store';
// const rows: any[] = [
//   {
//     id: '18002',
//     username: 'Jones Ferdinand',
//     email: 'jones.ferdinand@gmail.com',
//     role: 'Admin',
//     joined: '14/01/2022',
//     status: 'Completed',
//     verified_at: '14/01/2022',
//     organization: 'Cypertrace, Inc'
//   },
//   {
//     id: '18003',
//     username: 'Eason Yang',
//     email: 'eason.yang@gmail.com',
//     role: 'Member',
//     joined: '14/01/2022',
//     status: 'Completed',
//     organization: 'VASPnet, LLC'
//   },
//   {
//     id: '18001',
//     username: 'Anusha Aggarwal',
//     email: 'anusha.aggarwal@gmail.com',
//     role: 'Member',
//     joined: '14/01/2022',
//     status: 'Pending',
//     organization: 'VASPnet, LLC'
//   }
// ];

const getStatusBgColor = (status: string) => {
  switch (status && status.toLowerCase()) {
    case 'completed':
      return 'green.400';
    case 'pending':
      return 'yellow.400';
    case 'inactive':
      return 'red.400';
  }
};

const getCollaboratorActivatedDate = (verifiedAt: any) => {
  if (verifiedAt) {
    return formatIsoDate(verifiedAt);
  }
  return '-';
};

const isAuthorizedToInvite = () => {
  const userPermission = Store.getState().user?.user?.permissions;
  return userPermission?.includes(USER_PERMISSION.UPDATE_COLLABORATOR);
};

const RowItem: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return <Tr>{children}</Tr>;
};

const TableRow: React.FC<{ row: Collaborator }> = ({ row }) => {
  const { roles: userRoles } = useFetchUserRoles();
  return (
    <>
      <RowItem>
        <>
          <Td display="flex" flexDir="column">
            <chakra.span display="block" textTransform="capitalize">
              {row?.name}
            </chakra.span>
            <chakra.span display="block" fontSize="sm" color="gray.700">
              {row?.email}
            </chakra.span>
          </Td>
          <Td textTransform="capitalize">{row?.roles}</Td>
          <Td textTransform="capitalize">
            <Tag
              bg={row?.status ? getStatusBgColor(row?.status) : 'transparent'}
              color={'white'}
              size={'md'}>
              {row?.status}
            </Tag>
          </Td>
          <Td>{getCollaboratorActivatedDate(row?.verified_at)}</Td>
          <Td>{formatIsoDate(row?.created_at)}</Td>
          <Td textTransform="capitalize">{row?.organization}</Td>
          <Td paddingY={0}>
            <HStack width="100%" justifyContent="center" alignItems="center" spacing={5}>
              <Link color="blue" href={`mailto:${row?.email}`}>
                <FiMail fontSize="26px" />
              </Link>
              <EditCollaboratorModal collaboratorId={row?.id} roles={userRoles?.data} />
              <DeleteCollaboratorModal collaboratorId={row?.id} />
            </HStack>
          </Td>
        </>
      </RowItem>
    </>
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
            isDisabled={isAuthorizedToInvite()}>
            <Button minW="170px" onClick={modalHandler} isDisabled={!isAuthorizedToInvite()}>
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
              <Trans>Contact Detail</Trans>
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
            <Th>
              <Trans>Organization</Trans>
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
