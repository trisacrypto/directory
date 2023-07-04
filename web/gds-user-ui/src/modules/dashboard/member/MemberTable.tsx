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
  useColorModeValue
} from '@chakra-ui/react';
import UnverifiedMember from './UnverifiedMember';
import FormLayout from 'layouts/FormLayout';

import React from 'react';
import { Trans } from '@lingui/macro';
import { BsEye } from 'react-icons/bs';

interface MemberTableProps {
  data: any;
}

interface TableRowsProps {
  rows: any;
}

interface TableRowProps {
  row: any;
}

const TableRow: React.FC<TableRowProps> = (row: any) => {
  return (
    <Tr>
      {row.length > 0 ? (
        <>
          <Td>
            <chakra.span display="block"></chakra.span>
            <chakra.span display="block" fontSize="sm" color="gray.700"></chakra.span>
          </Td>
          <Td></Td>
          <Td></Td>
          <Td></Td>
          <Td></Td>
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
        </>
      ) : (
        <Td colSpan={6}>
          <UnverifiedMember />
        </Td>
      )}
    </Tr>
  );
};

const TableRows: React.FC<TableRowsProps> = (data) => {
  return (
    <>
      <TableRow row={data} />
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
          <TableRows rows={data} />
        </Tbody>
      </Table>
    </FormLayout>
  );
};
export default MemberTable;
