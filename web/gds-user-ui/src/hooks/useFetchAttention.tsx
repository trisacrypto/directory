import { useState, useEffect } from 'react';
import axios, { AxiosError, AxiosResponse } from 'axios';
import axiosInstance, { setAuthorization } from 'utils/axios';
const useFetchAttention = () => {
  const [attentionResponse, setAttentionResponse] = useState<any>();
  const [attentionError, setAttentionError] = useState<AxiosError | any>();
  const [attentionLoading, setAttentionLoading] = useState(true);

  const fetchAttentionData = async () => {
    try {
      setAuthorization();
      const result = await axiosInstance.get('/attention');
      if (result.status === 200) {
        setAttentionResponse(result.data);
      } else {
        setAttentionResponse([]);
      }
    } catch (err: any) {
      setAttentionError(err);
    } finally {
      setAttentionLoading(false);
    }
  };

  useEffect(() => {
    (async () => {
      await fetchAttentionData();
    })();
  }, []);

  return { attentionResponse, attentionError, attentionLoading };
};
export default useFetchAttention;
