import { useQuery } from '@tanstack/react-query';

import { getMembersService } from '../service';
import type { MembersQuery, DirectoryType } from '../memberType';

export function useFetchMembers(directory: DirectoryType): MembersQuery {
  const query = useQuery(['fetch-members', directory], () => getMembersService(directory), {
    retry: 0,
    enabled: !!directory
  });
  return {
    getMembers: query.refetch,
    members: query.data?.data,
    hasMembersFailed: query.isError,
    wasMembersFetched: query.isSuccess,
    isFetchingMembers: query.isLoading,
    error: query.error
  };
}
