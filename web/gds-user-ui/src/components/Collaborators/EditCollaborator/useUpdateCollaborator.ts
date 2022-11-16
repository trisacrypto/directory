import { useMutation } from '@tanstack/react-query';
import { updateCollaborator as updateCollaboratorService } from 'modules/dashboard/collaborator/CollaboratorService';
import type { UpdateCollaboratorMutation } from 'components/Collaborators/EditCollaborator/UpdateCollaboratorType';

export function useUpdateCollaborator(): UpdateCollaboratorMutation {
    const mutation = useMutation(['update-Collaborator'], updateCollaboratorService, {
        onError: (error: any) => {
            console.log('update-Collaborator-error', error);
            console.log('update-Collaborator-error-response', error?.response.data?.error);
        },
    });
    return {
        updateCollaborator: mutation.mutate,
        reset: mutation.reset,
        collaborator: mutation.data,
        hasCollaboratorFailed: mutation.isError,
        wasCollaboratorUpdated: mutation.isSuccess,
        isUpdating: mutation.isLoading,
        errorMessage: mutation.error?.response.data?.error || mutation.error,
    };
}
