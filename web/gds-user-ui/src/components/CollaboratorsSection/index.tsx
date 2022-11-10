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
  chakra,
  Link,
  useDisclosure
} from '@chakra-ui/react';
import FormLayout from 'layouts/FormLayout';
import React from 'react';
import { Trans } from '@lingui/macro';
import { FiMail } from 'react-icons/fi';
import EditCollaboratorModal from './EditCollaboratorModal';
import DeleteCollaboratorModal from './DeleteCollaboratorModal';
import AddCollaboratorModal from 'components/AddCollaboratorModal';

type Row = {
  id: string;
  username: string;
  role: string;
  joined: string;
  organization: string;
  email: string;
};

const rows: Row[] = [
  {
    id: '18001',
    username: 'Jones Ferdinand',
    email: 'jones.ferdinand@gmail.com',
    role: 'Admin',
    joined: '14/01/2022',
    organization: 'Cypertrace, Inc'
  },
  {
    id: '18001',
    username: 'Eason Yang',
    email: 'eason.yang@gmail.com',
    role: 'Member',
    joined: '14/01/2022',
    organization: 'VASPnet, LLC'
  },
  {
    id: '18001',
    username: 'Anusha Aggarwal',
    email: 'anusha.aggarwal@gmail.com',
    role: 'Member',
    joined: '14/01/2022',
    organization: 'VASPnet, LLC'
  }
];

const RowItem: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return <Tr>{children}</Tr>;
};

const TableRow: React.FC<{ row: Row }> = ({ row }) => {
  return (
    <>
      <RowItem>
        <>
          <Td display="flex" flexDir="column">
            <chakra.span display="block" textTransform="capitalize">
              {row.username}
            </chakra.span>
            <chakra.span display="block" fontSize="sm" color="gray.700">
              {row.email}
            </chakra.span>
          </Td>
          <Td textTransform="capitalize">{row.role}</Td>
          <Td>{row.joined}</Td>
          <Td textTransform="capitalize">{row.organization}</Td>
          <Td paddingY={0}>
            <HStack width="100%" justifyContent="center" alignItems="center" spacing={5}>
              <Link color="blue" href={`mailto:${row.email}`}>
                <FiMail fontSize="26px" />
              </Link>
              <EditCollaboratorModal />
              <DeleteCollaboratorModal />
            </HStack>
          </Td>
        </>
      </RowItem>
    </>
  );
};

const TableRows: React.FC = () => {
  return (
    <>
      {rows.map((row) => (
        <TableRow key={row.id} row={row} />
      ))}
    </>
  );
};

const CollaboratorsSection: React.FC = () => {
  const { onOpen, isOpen, onClose } = useDisclosure();
  const [isAddCollaboratorModalOpen, setIsAddCollaboratorModalOpen] = React.useState(false);
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
