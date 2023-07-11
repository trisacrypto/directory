import React, { useEffect } from 'react';
import MemberTableRows from './MemberTableRows';
import { useFetchMembers } from '../hooks/useFetchMembers';
import { memberSelector } from '../member.slice';
import { useSelector } from 'react-redux';
// import { mainnetMembersMockValue } from '../__mocks__';
const MemberTableList = () => {
  const network = useSelector(memberSelector).members.network;

  const { error, members, isFetchingMembers, getMembers } = useFetchMembers(network);
  const isUnverified = error && error?.response?.status === 451;

  // if network changes, we need to refetch members
  useEffect(() => {
    // console.log('network changed');
    getMembers();
  }, [network, getMembers]);

  return (
    <MemberTableRows rows={members?.vasps} hasError={isUnverified} isLoading={isFetchingMembers} />
  );
};

export default MemberTableList;
