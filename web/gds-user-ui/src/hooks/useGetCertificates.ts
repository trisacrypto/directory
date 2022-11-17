import { useQuery } from '@tanstack/react-query';
import axiosInstance from 'utils/axios';

const getCertificates = async () => {
  const response = await axiosInstance.get('/certificates');
  return response.data;
};

const useGetCertificates = () =>
  useQuery({
    queryKey: ['certificates'],
    queryFn: getCertificates,
    staleTime: 30 * 1000
  });

export default useGetCertificates;
