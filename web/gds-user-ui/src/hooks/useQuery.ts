/* eslint-disable @typescript-eslint/no-shadow */
import React from 'react';
import { useLocation } from 'react-router-dom';

function useQuery<QueryParams>() {
  const { search } = useLocation();
  return React.useMemo(() => {
    const params = new URLSearchParams(search);
    const result: any = {} as QueryParams;
    params.forEach((value, key) => {
      result[key] = value as unknown;
    });
    return result;
  }, [search]);
}

export default useQuery;
