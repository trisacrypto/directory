import { useEffect, useState } from 'react';
const getSearchParams = <T extends object>(): Partial<T> => {
  // server side rendering
  if (typeof window === 'undefined') {
    return {};
  }

  const params = new URLSearchParams(window.location.search);

  return new Proxy(params, {
    get(target, prop, receiver) {
      return target.get(prop as string) || undefined;
    }
  }) as T;
};

const useSearchParams = <T extends object = any>(): Partial<T> => {
  const [searchParams, setSearchParams] = useState(getSearchParams());
  const dependancy = typeof window === 'undefined' ? 'once' : window.location.search;

  useEffect(() => {
    setSearchParams(getSearchParams());
  }, [dependancy]);

  return searchParams;
};

export default useSearchParams;
