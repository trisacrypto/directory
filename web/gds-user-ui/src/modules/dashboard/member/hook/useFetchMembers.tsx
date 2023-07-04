import { useQuery } from '@tanstack/react-query';

import { getMembersService } from '../service';
import type { membersQuery, DirectoryType } from '../memberType';

export function useFetchMembers(directory?: DirectoryType): membersQuery {
  const query = useQuery(['fetch-members'], () => getMembersService(directory));
  return {
    getMembers: query.refetch,
    members: query.data?.data,
    hasMembersFailed: query.isError,
    wasMembersFetched: query.isSuccess,
    isFetchingMembers: query.isLoading,
    error: query.error
  };
}
