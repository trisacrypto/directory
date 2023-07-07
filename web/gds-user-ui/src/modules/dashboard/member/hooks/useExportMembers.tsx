import { useFetchMembers } from '../hooks/useFetchMembers';
import { useSelector } from 'react-redux';
import { memberSelector } from '../member.slice';
import { downloadCSV, convertToCvs } from 'utils/utils';
import { memberTableHeader } from '../utils';
const useExportMembers = () => {
  const { network } = useSelector(memberSelector);
  const { members, isFetchingMembers } = useFetchMembers(network);
  const exportHandler = () => {
    const data = convertToCvs(members?.vasps, memberTableHeader);
    downloadCSV(data, 'members');
  };
  return {
    exportHandler,
    isLoading: isFetchingMembers
  };
};

export { useExportMembers };
