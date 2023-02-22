import { APICore } from '@/helpers/api/apiCore';
import { MutationConfig } from '@/lib/react-query';
import { getCookie } from '@/utils';
import { useMutation } from '@tanstack/react-query';

const api = new APICore();

type UpdateContact = {
    vaspId: string;
    kind: string;
    data: any;
};

function updateContact({ vaspId, kind, data }: UpdateContact) {
    const csrfToken = getCookie('csrf_token');
    return api.update(`/vasps/${vaspId}/contacts/${kind}`, data, {
        headers: {
            'X-CSRF-TOKEN': csrfToken,
        },
    });
}

type UseUpdateReviewNotesOptions = {
    config?: MutationConfig<typeof updateContact>;
};

export const useUpdateContact = ({ config }: UseUpdateReviewNotesOptions = {}) => {
    return useMutation({
        ...config,
        mutationFn: updateContact,
    });
};
