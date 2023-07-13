import { useQuery } from '@tanstack/react-query';

import { getMemberService } from '../service';
import type { MemberQuery, MemberDto } from '../memberType';

export function useFetchMember(payload: MemberDto): MemberQuery {
  const query = useQuery(['fetch-member', payload.vaspId], () => getMemberService(payload), {
    retry: 0,
    enabled: !!payload.vaspId
  });
  return {
    getMember: query.refetch,
    member: query.data as any,
    hasMemberFailed: query.isError,
    wasMemberFetched: query.isSuccess,
    isFetchingMember: query.isLoading,
    error: query.error
  };
}
