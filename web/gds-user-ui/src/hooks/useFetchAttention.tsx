import { useState, useEffect } from 'react';
import { AxiosError } from 'axios';
import axiosInstance, { setAuthorization } from 'utils/axios';
import { handleError } from 'utils/utils';
const useFetchAttention = () => {
  const [attentionResponse, setAttentionResponse] = useState<any>();
  const [attentionError, setAttentionError] = useState<AxiosError | any>();
  const [attentionLoading, setAttentionLoading] = useState(true);

  const fetchAttentionData = async () => {
    try {
      setAuthorization();
      const result = await axiosInstance.get('/attention');
      if (result?.status === 200) {
        setAttentionResponse(result?.data);
      } else {
        setAttentionResponse([]);
      }
    } catch (err: any) {
      setAttentionError(err);
      handleError(err, '[useFetchAttention] fetch Attention Data failed');
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
