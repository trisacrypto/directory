import { useState, useEffect } from 'react';
import axios, { AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios';

axios.defaults.baseURL = process.env.REACT_APP_TRISA_API_URL;

const useAxios = (axiosParams: AxiosRequestConfig) => {
  const [response, setResponse] = useState<AxiosResponse>();
  const [error, setError] = useState<AxiosError | any>();
  const [loading, setLoading] = useState(true);

  const fetchData = async (params: AxiosRequestConfig) => {
    try {
      const result = await axios.request(params);
      if (result.data.success) {
        setResponse(result.data);
      } else {
        setError({ status: result.data.data.status });
      }
    } catch (err: any) {
      setError(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    (async () => {
      await fetchData(axiosParams);
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return { response, error, loading };
};
export default useAxios;
