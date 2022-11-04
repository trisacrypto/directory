import {
  Table,
  TableCaption,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
  Button,
  chakra,
  HStack,
  Link,
  useDisclosure,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
  Checkbox,
  VStack
} from '@chakra-ui/react';
import FormLayout from 'layouts/FormLayout';
import React from 'react';
import { Trans } from '@lingui/macro';
import { FiMail } from 'react-icons/fi';
import InputFormControl from 'components/ui/InputFormControl';
import EditCollaboratorModal from './EditCollaboratorModal';
import DeleteCollaboratorModal from './DeleteCollaboratorModal';

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
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <FormLayout overflowX={'scroll'}>
      <Table variant="simple">
        <TableCaption placement="top" textAlign="end" p={0} m={0} mb={3} fontSize={20}>
          <Button minW="170px" onClick={onOpen}>
            <Trans>Add Contact</Trans>
          </Button>
          <Modal isOpen={isOpen} onClose={onClose}>
            <ModalOverlay />
            <ModalContent w="100%" maxW="600px" px={10}>
              <ModalHeader textTransform="capitalize" textAlign="center" fontWeight={700} pb={1}>
                Add New Contact
              </ModalHeader>
              <ModalCloseButton />
              <ModalBody>
                <VStack>
                  <Text>
                    Please provide the name of the VASP and email address for the new contact.
                  </Text>
                  <Text>
                    The contact will receive an email to create a TRISA Global Directory Service
                    Account. The invitation to join is valid for 7 calendar days. The contact will
                    be added as a member for the VASP. The contact will have the ability to
                    contribute to certificate requests, check on the status of certificate requests,
                    and complete other actions related to the organization’s TRISA membership.
                  </Text>
                </VStack>
                <Checkbox mt={2} mb={4}>
                  TRISA is a network of trusted members. I acknowledge that the contact is
                  authorized to access the organization’s TRISA account information.
                </Checkbox>

                <InputFormControl
                  controlId="vasp_name"
                  label={
                    <>
                      <chakra.span fontWeight={700}>
                        <Trans>VASP Name</Trans>
                      </chakra.span>{' '}
                      (<Trans>required</Trans>)
                    </>
                  }
                />

                <InputFormControl
                  controlId="email"
                  label={
                    <>
                      <chakra.span fontWeight={700}>
                        <Trans>Email Address</Trans>
                      </chakra.span>{' '}
                      (<Trans>required</Trans>)
                    </>
                  }
                />
              </ModalBody>

              <ModalFooter display="flex" flexDir="column" gap={3}>
                <Button bg="orange" _hover={{ bg: 'orange' }} minW="150px">
                  Invite
                </Button>
                <Button
                  variant="ghost"
                  color="link"
                  fontWeight={400}
                  onClick={onClose}
                  minW="150px">
                  Cancel
                </Button>
              </ModalFooter>
            </ModalContent>
          </Modal>
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
