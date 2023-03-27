import { useQuery } from '@tanstack/react-query';
import { GetAllOrganisations } from './organizationService';
import type { Organization, OrganizationQuery } from './organizationType';
import { FETCH_ORGANIZATION } from 'constants/query-keys';

export function useOrganizationListQuery(): OrganizationQuery {
  const query = useQuery([FETCH_ORGANIZATION], GetAllOrganisations, {
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
