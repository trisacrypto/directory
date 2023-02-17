import { APICore } from '@/helpers/api/apiCore';
import { ExtractFnReturnType, QueryConfig } from '@/lib/react-query';
import { useQuery } from '@tanstack/react-query';

const api = new APICore();

export async function getAppVersion() {
  const response = await api.get('/status');

  return response.data
}

type QueryFnType = typeof getAppVersion;

type UseGetAppVersionOptions = {
  config?: QueryConfig<QueryFnType>;
};

export const useGetAppVersion = ({ config }: UseGetAppVersionOptions = {}) => {

    return useQuery<ExtractFnReturnType<QueryFnType>>({
        ...config,
        queryFn: getAppVersion,
        queryKey: ['get-version']
    })
}