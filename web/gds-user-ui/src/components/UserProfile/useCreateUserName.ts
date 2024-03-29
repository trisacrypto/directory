// create collaborator hook with axios and react-query
import { useMutation } from '@tanstack/react-query';
import { queryStrings } from 'utils/react-query';
import { updateUserFullName } from 'application/api/user';

type UpdateNameMutation = {
    updateName: (name: string) => void;
    isUpdating: boolean;
    hasUpdateFailed: boolean;
    errorMessage?: any;
    wasUpdated: boolean;
};

export function useCreateFullName(): UpdateNameMutation {
    const mutation = useMutation([queryStrings.updateNameKey], updateUserFullName, {
        onError: (error) => {
            console.log('success', error);
        }
    });

    return {
        updateName: mutation.mutate,
        isUpdating: mutation.isLoading,
        hasUpdateFailed: mutation.isError,
        errorMessage: mutation.error,
        wasUpdated: mutation.isSuccess
    };
}
