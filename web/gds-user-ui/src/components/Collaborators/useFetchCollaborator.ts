import { useQuery } from '@tanstack/react-query';

import { getAllCollaborators } from 'modules/dashboard/collaborator/CollaboratorService';
import type { getCollaborators } from 'modules/dashboard/collaborator/getCollaboratorType';

export function useFetchCollaborators(): getCollaborators {
    const query = useQuery(['fetch-collaborators'], getAllCollaborators, {
        refetchOnWindowFocus: true,
        refetchOnMount: true,
        staleTime: 0,
    });
    return {
        getAllCollaborators: query.refetch,
        collaborators: query.data,
        hasCollaboratorsFailed: query.isError,
        wasCollaboratorsFetched: query.isSuccess,
        isFetchingCollaborators: query.isLoading,
        errorMessage: query.error,
    };
}
