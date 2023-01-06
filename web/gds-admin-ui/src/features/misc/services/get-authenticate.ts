import { APICore, setCookie } from '@/helpers/api/apiCore';
import { UseScriptStatus } from '@/hooks/useScript';
import { ExtractFnReturnType, QueryConfig } from '@/lib/react-query';
import { getCookie } from '@/utils';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';

const apiCore = new APICore();

export async function getAuthenticate() {
    const response = await apiCore.get(`/authenticate`);

    return response.data;
}

type QueryFnType = typeof getAuthenticate;

type UseGetAuthenticateOptions = {
    config?: QueryConfig<QueryFnType>;
    loadScriptStatus: UseScriptStatus;
};

export const useGetAuthenticate = ({ config, loadScriptStatus }: UseGetAuthenticateOptions) => {
    return useQuery<ExtractFnReturnType<QueryFnType>>({
        ...config,
        onSuccess: () => {
            const csrfToken = getCookie('csrf_token');
            setCookie(csrfToken);
        },
        onError: (error: any) => {
            toast.error(error);
        },
        queryFn: getAuthenticate,
        enabled: loadScriptStatus === 'ready',
        queryKey: [],
        staleTime: 0,
        cacheTime: 0,
    });
};
