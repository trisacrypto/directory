import { useQuery } from '@tanstack/react-query';
import { getAllOrganisations } from './organizationService';
import type { OrganizationQuery } from './organizationType';
import { FETCH_ORGANIZATION } from 'constants/query-keys';

export function useOrganizationListQuery(page?: number): OrganizationQuery {
  const query = useQuery([FETCH_ORGANIZATION], () => getAllOrganisations(page), {
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    // set state time to 5 minutes
    staleTime: 1000 * 60 * 5
  });

  return {
    getAllOrganizations: query.refetch,
    organizations: query.data?.data as any,
    hasOrganizationFailed: query.isError,
    wasOrganizationFetched: query.isSuccess,
    isFetching: query.isFetching,
    errorMessage: query.error
  };
}
