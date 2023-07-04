import React from 'react';
import { MemberTableRows } from './MemberTableRows';
import { useFetchMembers } from '../hook/useFetchMembers';
import UnverifiedMember from '../UnverifiedMember';
import { Td, Tr } from '@chakra-ui/react';

const MemberTableList = () => {
  const { error, members } = useFetchMembers();
  if (error && error?.response?.status === 451) {
    return (
      <Tr>
        <Td colSpan={6}>
          <UnverifiedMember />
        </Td>
      </Tr>
    );
  }

  return <MemberTableRows rows={members?.vasps} />;
};

export default MemberTableList;
