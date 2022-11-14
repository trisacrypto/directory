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
  useDisclosure
} from '@chakra-ui/react';
import FormLayout from 'layouts/FormLayout';
import React, { useState } from 'react';
import { Trans } from '@lingui/macro';
import { FiMail } from 'react-icons/fi';
import EditCollaboratorModal from './EditCollaborator/EditCollaboratorModal';
import DeleteCollaboratorModal from './DeleteCollaborator/DeleteCollaboratorModal';
import AddCollaboratorModal from 'components/AddCollaboratorModal';
// import { getCollaborators, setCollaborators } from 'application/store/selectors/collaborator';
import { useFetchCollaborators } from './useFetchCollaborator';
// import { useDispatch } from 'react-redux';
import type { Collaborator } from './CollaboratorType';
import { formatIsoDate } from 'utils/formate-date';
import { sortCollaboratorsByRecentDate } from './lib';
import Loader from 'components/Loader';
import { useFetchUserRoles } from 'hooks/useFetchUserRoles';
// const rows: Row[] = [
//   {
//     id: '18001',
//     username: 'Jones Ferdinand',
//     email: 'jones.ferdinand@gmail.com',
//     role: 'Admin',
//     joined: '14/01/2022',
//     organization: 'Cypertrace, Inc'
//   },
//   {
//     id: '18001',
//     username: 'Eason Yang',
//     email: 'eason.yang@gmail.com',
//     role: 'Member',
//     joined: '14/01/2022',
//     organization: 'VASPnet, LLC'
//   },
//   {
//     id: '18001',
//     username: 'Anusha Aggarwal',
//     email: 'anusha.aggarwal@gmail.com',
//     role: 'Member',
//     joined: '14/01/2022',
//     organization: 'VASPnet, LLC'
//   }
// ];

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
          <Td>{formatIsoDate(row?.created_at)}</Td>
          <Td textTransform="capitalize"></Td>
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
  console.log('[collaborators]', collaborators);
  return (
    <>
      {sortCollaboratorsByRecentDate(collaborators?.data?.collaborators).map(
        (collaborator: Collaborator) => (
          <TableRow key={collaborator.id} row={collaborator} />
        )
      )}
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
          <Button minW="170px" onClick={modalHandler}>
            <Trans>Add Contact</Trans>
          </Button>
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
              <Trans>Name & Email</Trans>
            </Th>
            <Th>
              <Trans>Role</Trans>
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
