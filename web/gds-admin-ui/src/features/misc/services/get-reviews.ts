import { APICore } from "@/helpers/api/apiCore";
import { ExtractFnReturnType, QueryConfig } from "@/lib/react-query";
import { useQuery } from "@tanstack/react-query";

const apiCore = new APICore()

export async function getReviews(){
  const response = await apiCore.get(`/reviews`);

  return response.data
}

type QueryFnType = typeof getReviews;

type UseGetReviewsOptions = {
  config?: QueryConfig<QueryFnType>;
};

export const useGetReviews = ({ config }: UseGetReviewsOptions = {}) => {

    return useQuery<ExtractFnReturnType<QueryFnType>>({
        ...config,
        queryFn: getReviews,
        queryKey: ['get-reviews']
    })
}