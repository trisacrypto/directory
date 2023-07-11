import { useState, useEffect } from 'react';
import { useFetchMembers } from '../hooks/useFetchMembers';
import { useSelector } from 'react-redux';
import { memberSelector } from '../member.slice';

import { downloadMembers2CVS } from '../utils';
const useExportMembers = () => {
  const network = useSelector(memberSelector).members.network;
  const { members, getMembers } = useFetchMembers(network);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const LOADING_TIMEOUT = 500;
  const exportHandler = () => {
    try {
      setIsLoading(true);
      setTimeout(() => {
        // simulate loading time for user experience
        downloadMembers2CVS(members?.vasps);
        setIsLoading(false);
      }, LOADING_TIMEOUT);
    } catch (error) {
      console.log('[useExportMembers] error: ', error);
    }
  };

  useEffect(() => {
    if (isLoading) {
      console.log('[ExportButton]');
      getMembers();
    }
  }, [isLoading, getMembers]);

  return {
    exportHandler,
    isLoading
  };
};

export { useExportMembers };
