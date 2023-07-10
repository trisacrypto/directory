import React, { useEffect } from 'react';
import MemberTableRows from './MemberTableRows';
import { useFetchMembers } from '../hooks/useFetchMembers';
import { memberSelector } from '../member.slice';
import { useSelector } from 'react-redux';
const MemberTableList = () => {
  const { network } = useSelector(memberSelector);

  const { error, members, isFetchingMembers, getMembers } = useFetchMembers(network);
  const isUnverified = error && error?.response?.status === 451;

  // if network changes, we need to refetch members
  useEffect(() => {
    getMembers();
  }, [network, getMembers]);

  return (
    <MemberTableRows rows={members?.vasps} hasError={isUnverified} isLoading={isFetchingMembers} />
  );
};

export default MemberTableList;
