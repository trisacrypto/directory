import { useState, useEffect } from 'react';
import { useFetchMembers } from '../hooks/useFetchMembers';
import { useSelector } from 'react-redux';
import { memberSelector } from '../member.slice';

import { downloadMembers2CVS } from '../utils';
const useExportMembers = () => {
  const network = useSelector(memberSelector).members.network;
  const { members, getMembers, error } = useFetchMembers(network);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const LOADING_TIMEOUT = 500;

  const isUnverified = error && error?.status === 451;
  const isMemberEmpty = !members?.vasps?.length;
  const hasError = !isUnverified && error; // if error is not 451, then it's a real error

  const exportHandler = () => {
    try {
      setIsLoading(true);
      setTimeout(() => {
        // simulate loading time for user experience
        downloadMembers2CVS(members?.vasps);
        setIsLoading(false);
      }, LOADING_TIMEOUT);
    } catch (er: any) {
    }
  };

  useEffect(() => {
    if (isLoading) {
      getMembers();
    }
  }, [isLoading, getMembers, error]);

  useEffect(() => {
    // this should avoid the loading state to be stuck
    if (error) {
      setIsLoading(false);
    }
  }, [error]);

  return {
    exportHandler,
    isLoading,
    isDisabled: isUnverified || isMemberEmpty || hasError
  };
};

export { useExportMembers };
