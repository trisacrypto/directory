import { useQuery } from '@tanstack/react-query';
import { getAllOrganisations } from './organizationService';
import type { OrganizationQuery, OrganizationResponse } from './organizationType';
import { FETCH_ORGANIZATION } from 'constants/query-keys';
export function useOrganizationListQuery({ name = '', page = 1, pageSize = 8 }): OrganizationQuery {
  const query = useQuery(
    [FETCH_ORGANIZATION, page],
    () => getAllOrganisations(name, page, pageSize),
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
