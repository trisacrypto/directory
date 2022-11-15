import { QueryCache } from '@tanstack/react-query';

interface IQueryCache {
  queryKey: string;
}

const useFetchQueryCache = ({ queryKey }: IQueryCache) => {
  const queryCache = new QueryCache();
  const k = [queryKey];
  return queryCache.find(k);
};

export default useFetchQueryCache;
