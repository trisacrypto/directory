import { useQuery } from '@tanstack/react-query';

import { getMemberService } from '../service';
import type { MemberQuery } from '../memberType';

export function useFetchMembers(vaspID: string): MemberQuery {
  const query = useQuery(['fetch-member', vaspID], () => getMemberService, {
    retry: 0,
    enabled: !!vaspID
  });
  return {
    getMember: query.refetch,
    member: query.data,
    hasMemberFailed: query.isError,
    wasMemberFetched: query.isSuccess,
    isFetchingMember: query.isLoading,
    error: query.error
  };
}
