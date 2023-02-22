import { APICore } from '@/helpers/api/apiCore';
import { MutationConfig } from '@/lib/react-query';
import { getCookie } from '@/utils';
import { useMutation } from '@tanstack/react-query';

type DeleteContact = {
    vaspId: string;
    kind: string;
};

const api = new APICore();

export function deleteContact({ vaspId, kind }: DeleteContact) {
    const csrfToken = getCookie('csrf_token');
    return api.delete(`/vasps/${vaspId}/contacts/${kind}`, {
        headers: {
            'X-CSRF-TOKEN': csrfToken,
        },
    });
}

type UseDeleteContactOptions = {
    config?: MutationConfig<typeof deleteContact>;
};

export const useDeleteContact = ({ config }: UseDeleteContactOptions = {}) => {
    return useMutation({
        ...config,
        mutationFn: deleteContact,
    });
};
