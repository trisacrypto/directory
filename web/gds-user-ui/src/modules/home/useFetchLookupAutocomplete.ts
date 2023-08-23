/* eslint-disable @typescript-eslint/no-shadow */
import { useEffect, useState } from 'react';
import { lookupAutocomplete } from './service';

export const useFetchLookupAutocomplete = () => {
  const [autocomplete, setAutocomplete] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<any>(null);

  useEffect(() => {
    const fetchAutocomplete = async () => {
      setIsLoading(true);
      try {
        const data = await lookupAutocomplete();
        setAutocomplete(data);
        // eslint-disable-next-line no-catch-shadow
      } catch (error: any) {
        setError(error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchAutocomplete();
  }, []);

  return { vasps: autocomplete, isLoading, error };
};

export default useFetchLookupAutocomplete;
