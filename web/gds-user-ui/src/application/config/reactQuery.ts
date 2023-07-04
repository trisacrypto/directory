import { MutationCache, QueryCache, QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient({
  // define default options for all queries
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 15,
      refetchInterval: false,
      refetchIntervalInBackground: false,
      refetchOnMount: true,
      refetchOnWindowFocus: false
    }
  }
});

const queryCache = new QueryCache({
  onError: (error) => {
    console.log(error);
  },
  onSuccess: (data) => {
    console.log(data);
  }
});
const mutationCache = new MutationCache({
  onError: (error) => {
    console.log(error);
  },
  onSuccess: (data) => {
    console.log(data);
  }
});

export { mutationCache, queryCache, queryClient, QueryClientProvider };
