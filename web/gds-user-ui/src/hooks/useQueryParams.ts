import { useEffect, useState } from 'react';
const getSearchParams = <T extends object>(): Partial<T> => {
  // server side rendering
  if (typeof window === 'undefined') {
    return {};
  }

  const params = new URLSearchParams(window.location.search);

  return new Proxy(params, {
    get(target, prop) {
      return target.get(prop as string) || undefined;
    }
  }) as T;
};

const useSearchParams = <T extends object = any>(): Partial<T> => {
  const [searchParams, setSearchParams] = useState(getSearchParams());
  const dep = typeof window === 'undefined' ? 'once' : window.location.search;
  useEffect(() => {
    setSearchParams(getSearchParams());
  }, [dep]);

  return searchParams;
};

export default useSearchParams;
