import { Table, TableCaption, Tbody, Th, Thead, Tr } from '@chakra-ui/react';

import React from 'react';
import { Trans } from '@lingui/macro';

import MemberTableList from './MemberTableList';

// interface MemberTableProps {
//   data: any;
// }

const MemberTable: React.FC = () => {
  return (
    <Table variant="simple">
      <TableCaption placement="top" textAlign="end" p={0} m={0} mb={3} fontSize={20}></TableCaption>
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
        <MemberTableList />
      </Tbody>
    </Table>
  );
};
export default MemberTable;
