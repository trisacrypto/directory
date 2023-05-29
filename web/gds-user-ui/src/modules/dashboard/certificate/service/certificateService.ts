import axiosInstance from 'utils/axios';
import type { PayloadDTO } from 'modules/dashboard/certificate/types';

export const getCertificateStepService = async (payload: PayloadDTO) => {
  const { key } = payload;
  const response = await axiosInstance.get(`/register?step=${key}`);
  return response.data;
};

export const postCertificateStepService = async (payload: any) => {
  const response = await axiosInstance.put('/register', {
    ...payload
  });
  return response.data;
};
