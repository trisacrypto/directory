// create collaborator hook with axios and react-query
import { useMutation } from 'react-query';
import { createCollaborator } from 'components/AddCollaboratorModal/CollaboratorService';
export function useCreateCollaborator() {
    const { mutate, isLoading, isError, error, isSuccess } = useMutation(createCollaborator);

    return {
        createCollaborator: mutate,
        isLoading,
        isError,
        error,
        isSuccess,

    };
}
