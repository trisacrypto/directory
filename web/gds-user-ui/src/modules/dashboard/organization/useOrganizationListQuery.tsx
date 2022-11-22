import { useQuery } from '@tanstack/react-query';
import { GetAllOrganisations } from './organizationService';
import type { Organization, OrganizationQuery } from './organizationType';

export function useOrganizationListQuery(): OrganizationQuery {
  const query = useQuery(['fetch-organization'], GetAllOrganisations, {
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    // set state time to 5 minutes
    staleTime: 1000 * 60 * 5
  });
  return {
    getAllOrganizations: query.refetch,
    organizations: query.data?.data?.organizations as Organization[],
    hasOrganizationFailed: query.isError,
    wasOrganizationFetched: query.isSuccess,
    isFetching: query.isFetching,
    errorMessage: query.error
  };
}
