import { useQuery } from '@tanstack/react-query';
import { useToast } from '@chakra-ui/react';

import { getAllCollaborators } from 'modules/dashboard/collaborator/CollaboratorService';
import type { getCollaborators } from 'modules/dashboard/collaborator/getCollaboratorType';

export function useFetchCollaborators(): getCollaborators {
  const toast = useToast();
  const query = useQuery(['fetch-collaborators'], getAllCollaborators, {
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    // set state time to 15 minutes
    staleTime: 1000 * 60 * 15,

    onError: (error: any) => {
      toast({
        title: 'Error',
        description: error?.message,
        status: 'error',
        duration: 5000,
        isClosable: true
      });
    }
  });
  return {
    getAllCollaborators: query.refetch,
    collaborators: query.data?.data?.collaborators,
    hasCollaboratorsFailed: query.isError,
    wasCollaboratorsFetched: query.isSuccess,
    isFetchingCollaborators: query.isLoading,
    error: query.error
  };
}
