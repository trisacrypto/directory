import { useState } from 'react';
import { useFetchMembers } from '../hooks/useFetchMembers';
import { useSelector } from 'react-redux';
import { memberSelector } from '../member.slice';
import { downloadCSV, convertToCvs } from 'utils/utils';
import { memberTableHeader } from '../utils';
const useExportMembers = () => {
  const { network } = useSelector(memberSelector);
  const { members } = useFetchMembers(network);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const LOADING_TIMEOUT = 500;
  const exportHandler = () => {
    const data = convertToCvs(members?.vasps, memberTableHeader);
    try {
      setIsLoading(true);
      setTimeout(() => {
        downloadCSV(data, 'members');

        setIsLoading(false);
      }, LOADING_TIMEOUT);
    } catch (error) {
      console.log(error);
    }
  };

  return {
    exportHandler,
    isLoading
  };
};

export { useExportMembers };
