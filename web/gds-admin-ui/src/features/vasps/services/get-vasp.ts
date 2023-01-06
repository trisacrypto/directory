import { APICore } from "@/helpers/api/apiCore";
import { ExtractFnReturnType, QueryConfig } from "@/lib/react-query";
import { useQuery } from "@tanstack/react-query";

const apiCore = new APICore()

export async function getVasp(vaspId: string){
  const response = await apiCore.get(`/vasps/${vaspId}`);

  return response.data
}

type QueryFnType = typeof getVasp;

type UseGetVaspOptions = {
  config?: QueryConfig<QueryFnType>;
  vaspId: string;
};

export const useGetVasp = ({ config, vaspId }: UseGetVaspOptions) => {

    return useQuery<ExtractFnReturnType<QueryFnType>>({
        ...config,
        enabled: !!vaspId,
        queryFn: () => getVasp(vaspId),
        queryKey: ['get-vasps', vaspId],
    })
}