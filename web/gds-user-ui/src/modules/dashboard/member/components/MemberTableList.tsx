import React, { useEffect } from 'react';
import MemberTableRows from './MemberTableRows';
import { useFetchMembers } from '../hooks/useFetchMembers';
import { memberSelector } from '../member.slice';
import { useSelector } from 'react-redux';
import { mainnetMembersMockValue } from '../__mocks__';
// import { mainnetMembersMockValue } from '../__mocks__';
const MemberTableList = () => {
  const { network } = useSelector(memberSelector);

  const mainnet = mainnetMembersMockValue;

  const { error, /* members, */ isFetchingMembers, getMembers } = useFetchMembers(network);
  const isUnverified = error && error?.response?.status === 200;

  // if network changes, we need to refetch members
  useEffect(() => {
    getMembers();
  }, [network, getMembers]);

  return (
    <MemberTableRows rows={mainnet?.vasps} hasError={isUnverified} isLoading={isFetchingMembers} />
  );
};

export default MemberTableList;
