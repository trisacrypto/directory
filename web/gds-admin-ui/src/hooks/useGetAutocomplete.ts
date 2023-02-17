import { APICore } from "@/helpers/api/apiCore";
import { ExtractFnReturnType, QueryConfig } from "@/lib/react-query";
import { useQuery } from "@tanstack/react-query";

const apiCore = new APICore()

export async function getAutocompletes(){
  const response = await apiCore.get(`/autocomplete`);

  return response.data.names
}

type QueryFnType = typeof getAutocompletes;

type UseGetAutocompletesOptions = {
  config?: QueryConfig<QueryFnType>;
};

export const useGetAutocompletes = ({ config }: UseGetAutocompletesOptions = {}) => {

    return useQuery<ExtractFnReturnType<QueryFnType>>({
        ...config,
        queryFn: getAutocompletes,
        queryKey: ['get-autocompletes']
    })
}