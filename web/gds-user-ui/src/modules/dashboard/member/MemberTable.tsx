import {
  Table,
  TableCaption,
  Tbody,
  Th,
  Thead,
  Heading,
  Tr,
  Button,
  HStack,
  useColorModeValue
} from '@chakra-ui/react';
import FormLayout from 'layouts/FormLayout';

import React from 'react';
import { Trans } from '@lingui/macro';
import { MemberTableRows } from './components/MemberTableRows';

import MemberSelectNetwork from './Components/MemberNetworkSelect';

interface MemberTableProps {
  data: any;
}

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
          <MemberSelectNetwork />
        </TableCaption>
        <Thead>
          <Tr>
            <Th data-testid="name-header">
              <Trans>Member Name</Trans>
            </Th>
            <Th data-testid="joined-header">
              <Trans>Joined</Trans>
            </Th>
            <Th data-testid="last-updated-header">
              <Trans>Last Updated</Trans>
            </Th>
            <Th data-testid="network-header">
              <Trans>Network</Trans>
            </Th>
            <Th data-testid="status-header">
              <Trans>Status</Trans>
            </Th>
            <Th textAlign="center" data-testid="actions-header">
              <Trans>Actions</Trans>
            </Th>
          </Tr>
        </Thead>
        <Tbody>
          <MemberTableRows rows={data} />
        </Tbody>
      </Table>
    </FormLayout>
  );
};
export default MemberTable;
