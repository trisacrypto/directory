import {
  Table,
  TableCaption,
  Tbody,
  Td,
  Th,
  Thead,
  Heading,
  Tr,
  Button,
  HStack,
  chakra,
  useColorModeValue,
  Tag
} from '@chakra-ui/react';

import FormLayout from 'layouts/FormLayout';

import React from 'react';
import { Trans } from '@lingui/macro';
import { BsEye } from 'react-icons/bs';
import { mainnetMembersMockValue } from './__mocks__';
import { formatIsoDate } from 'utils/formate-date';

interface MemberTableProps {
  data: any;
}

const vasps = mainnetMembersMockValue.vasps;

const TableRow: React.FC = () => {
  return (
    vasps.map((member: any) => (
      <Tr key={member.id}>
        <Td>
          <chakra.span display="block">{member.name}</chakra.span>
        </Td>
        <Td>{formatIsoDate(member.first_listed)}</Td>
        <Td>{formatIsoDate(member.last_updated)}</Td>
        <Td>
          {member.registered_directory === 'vaspdirectory.net' && <span>MainNet</span>}
          {member.registered_directory === 'trisatest.net' && <span>TestNet</span>}
        </Td>
        <Td>
          <Tag bg="green.400" color="white">{member.status}</Tag>
        </Td>
        <Td paddingY={0}>
          <HStack width="100%" justifyContent="center" alignItems="center">
            <Button
              color="blue"
              as={'a'}
              href={``}
              bg={'transparent'}
              _hover={{
                bg: 'transparent'
              }}
              _focus={{
                bg: 'transparent'
              }}>
              <BsEye fontSize="24px" />
            </Button>
          </HStack>
        </Td>
      </Tr>
    ))
  );
};

const TableRows: React.FC = () => {
  return (
    <>
      <TableRow />
    </>
  );
};

const MemberTable: React.FC<MemberTableProps> = (data) => {
  console.log('data', data);
  const modalHandler = () => {
    console.log('modalHandler');
  };

  return (
    <FormLayout overflowX={'scroll'}>
      <Table variant="simple">
        <TableCaption placement="top" textAlign="end" p={0} m={0} mb={3} fontSize={20}>
          <HStack justify={'space-between'} mb={'10'}>
            <Heading size="md" color={'black'}>
              <Trans>Member List</Trans>
            </Heading>
            <Button
              minW="100px"
              onClick={modalHandler}
              bg={useColorModeValue('black', 'white')}
              _hover={{
                bg: useColorModeValue('black', 'white')
              }}
              color={useColorModeValue('white', 'black')}>
              <Trans>Export</Trans>
            </Button>
          </HStack>
        </TableCaption>
        <Thead>
          <Tr>
            <Th>
              <Trans>Member Name</Trans>
            </Th>
            <Th>
              <Trans>Joined</Trans>
            </Th>
            <Th>
              <Trans>Last Updated</Trans>
            </Th>
            <Th>
              <Trans>Network</Trans>
            </Th>
            <Th>
              <Trans>Status</Trans>
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
export default MemberTable;
