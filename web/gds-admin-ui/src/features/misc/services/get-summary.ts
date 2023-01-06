import { APICore } from "@/helpers/api/apiCore";
import { ExtractFnReturnType, QueryConfig } from "@/lib/react-query";
import { useQuery } from "@tanstack/react-query";

const apiCore = new APICore()

export async function getSummary(){
  const response = await apiCore.get(`/summary`);

  return response.data
}

type QueryFnType = typeof getSummary;

type UseGetSummaryOptions = {
  config?: QueryConfig<QueryFnType>;
};

export const useGetSummary = ({ config }: UseGetSummaryOptions = {}) => {

    return useQuery<ExtractFnReturnType<QueryFnType>>({
        ...config,
        queryFn: getSummary,
        queryKey: ['get-summary']
    })
}