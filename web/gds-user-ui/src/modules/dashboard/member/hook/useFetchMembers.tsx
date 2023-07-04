import { useQuery } from '@tanstack/react-query';

import { getMembersService } from '../service';
import type { membersQuery } from '../memberType';

export function useFetchMembers(): membersQuery {
  const query = useQuery(['fetch-members'], getMembersService, {
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    // set state time to 15 minutes
    staleTime: 1000 * 60 * 15
  });
  return {
    getMembers: query.refetch,
    members: query.data?.data?.members,
    hasMembersFailed: query.isError,
    wasMembersFetched: query.isSuccess,
    isFetchingMembers: query.isLoading,
    errorMessage: query.error
  };
}
