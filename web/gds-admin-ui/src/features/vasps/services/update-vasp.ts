import { useModalContext } from '@/components/Modal';
import { APICore } from '@/helpers/api/apiCore';
import { MutationConfig, queryClient } from '@/lib/react-query';
import { getCookie } from '@/utils';
import { captureException } from '@sentry/react';
import { useMutation } from '@tanstack/react-query';

const api = new APICore();

type UpdateVasp = {
    vaspId: string;
    data: any;
};

export function updateVasp({ vaspId, data }: UpdateVasp) {
    const csrfToken = getCookie('csrf_token');
    return api.patch(`/vasps/${vaspId}`, data, {
        headers: {
            'X-CSRF-TOKEN': csrfToken,
        },
    });
}

type UseUpdateVaspOptions = {
    config?: MutationConfig<typeof updateVasp>;
};

export const useUpdateVasp = ({ config }: UseUpdateVaspOptions = {}) => {
    const { closeModal } = useModalContext();

    return useMutation({
        ...config,
        onSuccess() {
            closeModal();
            queryClient.invalidateQueries();
        },
        onError(error) {
            captureException(error);
        },
        mutationFn: updateVasp,
    });
};
