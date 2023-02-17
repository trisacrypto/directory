import { APICore } from "@/helpers/api/apiCore";
import { ExtractFnReturnType, QueryConfig } from "@/lib/react-query";
import { useQuery } from "@tanstack/react-query";

const apiCore = new APICore()

export async function getVasps(queryParams?: string){
  const response = await apiCore.get(`/vasps`, queryParams);

  return response.data
}

type QueryFnType = typeof getVasps;

type UseGetVaspsOptions = {
  config?: QueryConfig<QueryFnType>;
  queryParams?: string;
};

export const useGetVasps = ({ config, queryParams }: UseGetVaspsOptions = {}) => {

    return useQuery<ExtractFnReturnType<QueryFnType>>({
        ...config,
        keepPreviousData: true,
        queryFn: () => getVasps(queryParams),
        queryKey: ['get-vasps', queryParams]
    })
}