import { Dispatch } from '@reduxjs/toolkit';
import { useSelector, useDispatch } from 'react-redux';
import { setMemberNetwork, setDefaultMemberNetwork, memberSelector } from '../member.slice';
const useMemberState = () => {
  const dispatch: Dispatch = useDispatch<Dispatch<any>>();

  const { network } = useSelector(memberSelector).members;

  const setNetwork = (value: string) => {
    dispatch(setMemberNetwork(value));
  };

  const getNetwork = () => {
    return network;
  };

  const setDefaultNetwork = () => {
    dispatch(setDefaultMemberNetwork());
  };

  return {
    setNetwork,
    getNetwork,
    setDefaultNetwork
  };
};

export default useMemberState;
