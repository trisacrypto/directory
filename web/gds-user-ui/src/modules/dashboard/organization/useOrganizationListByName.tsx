import { useQuery } from '@tanstack/react-query';
import { getOrganizationByName } from './organizationService';
import type { OrganizationQuery, OrganizationResponse } from './organizationType';
import { FETCH_ORGANIZATION_BY_NAME } from 'constants/query-keys';
export function useOrganizationListByName(name = '', page = 1, pageSize = 8): OrganizationQuery {
  const query = useQuery(
    [FETCH_ORGANIZATION_BY_NAME, name],
    () => getOrganizationByName(name, page, pageSize),
    {
      refetchOnWindowFocus: false,
      refetchOnMount: true,
      // set state time to 5 minutes
      staleTime: 1000 * 60 * 5
    }
  );

  return {
    getAllOrganizations: query.refetch,
    organizations: query.data?.data as OrganizationResponse,
    hasOrganizationFailed: query.isError,
    wasOrganizationFetched: query.isSuccess,
    isFetching: query.isFetching,
    errorMessage: query.error
  };
}
