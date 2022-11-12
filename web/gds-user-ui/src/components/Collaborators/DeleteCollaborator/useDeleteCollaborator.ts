import { useMutation } from '@tanstack/react-query';
import { deleteCollaborator as deleteCollaboratorService } from 'modules/dashboard/collaborator/CollaboratorService';
import type { DeleteCollaboratorMutation } from 'components/Collaborators/DeleteCollaborator/DeleteCollaboratorType';

export function useDeleteCollaborator(): DeleteCollaboratorMutation {
    const mutation = useMutation(['deleteCollaborator'], deleteCollaboratorService);
    return {
        deleteCollaborator: mutation.mutate,
        reset: mutation.reset,
        collaborator: mutation.data,
        hasCollaboratorFailed: mutation.isError,
        wasCollaboratorDeleted: mutation.isSuccess,
        isDeleting: mutation.isLoading,
        errorMessage: mutation.error,
    };
}

