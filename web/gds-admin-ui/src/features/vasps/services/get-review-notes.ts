import { APICore } from '@/helpers/api/apiCore';
import { ExtractFnReturnType, QueryConfig } from '@/lib/react-query';
import { useQuery } from '@tanstack/react-query';
import { GET_REVIEW_NOTES } from '../constants';

const apiCore = new APICore();

export async function getReviewNotes(vaspId: string) {
    const response = await apiCore.get(`/vasps/${vaspId}/notes`);

    return response.data?.notes;
}

type QueryFnType = typeof getReviewNotes;

type UseGetReviewNotesOptions = {
    config?: QueryConfig<QueryFnType>;
    vaspId: string;
};

export const useGetReviewNotes = ({ config, vaspId }: UseGetReviewNotesOptions) => {
    return useQuery<ExtractFnReturnType<QueryFnType>>({
        ...config,
        enabled: !!vaspId,
        select(data) {
            const sortedData =
                data && data?.length
                    ? data.sort((a: any, b: any) => {
                          const date1: any = a?.modified ? new Date(a?.modified) : new Date(a?.created);
                          const date2: any = b?.modified ? new Date(b?.modified) : new Date(b?.created);

                          return date2 - date1;
                      })
                    : [];

            return sortedData;
        },
        queryFn: () => getReviewNotes(vaspId),
        queryKey: [GET_REVIEW_NOTES, vaspId],
    });
};
