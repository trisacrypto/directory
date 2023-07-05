import React from 'react';
import MemberTableRows from './MemberTableRows';
import { useFetchMembers } from '../hooks/useFetchMembers';

const MemberTableList = () => {
  const { error, members, isFetchingMembers } = useFetchMembers();
  const isUnverified = error && error?.response?.status === 451;

  return (
    <MemberTableRows rows={members?.vasps} hasError={isUnverified} isLoading={isFetchingMembers} />
  );
};

export default MemberTableList;
