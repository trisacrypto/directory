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
  Tooltip
} from '@chakra-ui/react';
import FormLayout from 'layouts/FormLayout';
import React from 'react';
import { Trans, t } from '@lingui/macro';

const TableRow: React.FC = () => {
  return (
    <Tr>
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
            }}></Button>
        </HStack>
      </Td>
    </Tr>
  );
};

const TableRows: React.FC = () => {
  return (
    <>
      <TableRow />
    </>
  );
};

const MemberTable: React.FC = () => {
  const modalHandler = () => {
    console.log('modalHandler');
  };

  return (
    <FormLayout overflowX={'scroll'}>
      <Table variant="simple">
        <TableCaption placement="top" textAlign="end" p={0} m={0} mb={3} fontSize={20}>
          <Tooltip label={t`you do not have permission to invite a collaborator`}>
            <Button minW="170px" onClick={modalHandler} bg={'black'}>
              <Trans>Export</Trans>
            </Button>
          </Tooltip>
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
              <Trans>Metwork</Trans>
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
