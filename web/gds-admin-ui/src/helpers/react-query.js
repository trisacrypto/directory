import { QueryClient } from '@tanstack/react-query';

const QUERY_STALE_TIME = 30 * 1000

const defaultOptions = {
    queries: {
        refetchOnWindowFocus: false,
        staleTime: QUERY_STALE_TIME,
        retry: false
    }
}

const queryClient = new QueryClient({
    defaultOptions
});

export default queryClient;
