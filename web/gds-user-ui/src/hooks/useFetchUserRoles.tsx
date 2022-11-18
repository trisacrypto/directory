import { useQuery } from '@tanstack/react-query';

import axiosInstance from 'utils/axios';
export function useFetchUserRoles(): any {
  const query = useQuery(
    ['fetch-userRoles'],
    async () => {
      return await axiosInstance.get('/users/roles');
    },
    {
      refetchOnWindowFocus: true,
      refetchOnMount: true,
      staleTime: 0
    }
  );
  return {
    roles: query.data,
    hasUserRolesFailed: query.isError,
    wasUserRolesFetched: query.isSuccess,
    isFetchingUserRoles: query.isLoading,
    errorMessage: query.error,
    refetchUserRoles: query.refetch
  };
}
