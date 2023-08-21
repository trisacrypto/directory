import { useState } from 'react';

import { isValidUuid } from 'utils/utils';
import { lookup } from './service';

const useFetchLookup = () => {
  const [error, setError] = useState<any>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [data, setData] = useState<any>(null);
  const [searchString, setSearchString] = useState<string>('');

  const resetData = () => {
    setData(null);
    setSearchString('');
    setError(null);
  };
  const handleSearch = (searchQuery: string) => {
    const query = isValidUuid(searchQuery) ? `uuid=${searchQuery}` : `common_name=${searchQuery}`;

    const fetchLookup = async () => {
      setIsLoading(true);
      try {
        const response = await lookup(query);
        if (!response.mainnet || !response.testnet) setError('No data found');
        setData(data);
        setSearchString(searchQuery);
      } catch (e: any) {
        if (!e?.response?.data?.success) {
          setError(e?.response?.data?.error);
        } else {
          setError('Something went wrong');
        }
      } finally {
        setIsLoading(false);
      }
    };
    fetchLookup();
  };

  return { data, isLoading, error, handleSearch, searchString, resetData };
};

export default useFetchLookup;
