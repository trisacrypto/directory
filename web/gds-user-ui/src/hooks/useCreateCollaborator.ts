// create collaborator hook with axios and react-query
import { useMutation } from '@tanstack/react-query';
import { createCollaborator } from 'modules/dashboard/collaborator/CollaboratorService';
import type { CollaboratorMutation } from 'components/AddCollaboratorModal/AddCollaboratorType';
export function useCreateCollaborator(): CollaboratorMutation {
    const mutation = useMutation(['addCollaborator'], createCollaborator);

    return {
        createCollaborator: mutation.mutate,
        reset: mutation.reset,
        collaborator: mutation.data,
        hasCollaboratorFailed: mutation.isError,
        wasCollaboratorCreated: mutation.isSuccess,
        isCreating: mutation.isLoading,
        errorMessage: mutation.error,
    };
}
