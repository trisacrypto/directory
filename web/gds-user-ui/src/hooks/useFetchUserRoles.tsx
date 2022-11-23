import { useQuery } from '@tanstack/react-query';

import axiosInstance from 'utils/axios';
export function useFetchUserRoles(): any {
  const query = useQuery(
    ['fetch-userRoles'],
    async () => {
      return await axiosInstance.get('/users/roles');
    },
    {
      refetchOnWindowFocus: false,
      refetchOnMount: true,
      staleTime: 1000 * 60 * 60 * 24
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
