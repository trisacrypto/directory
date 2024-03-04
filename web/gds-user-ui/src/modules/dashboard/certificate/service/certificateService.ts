import axiosInstance from 'utils/axios';
import type { PayloadDTO } from 'modules/dashboard/certificate/types';

export const getCertificateStepService = async (payload: PayloadDTO) => {
  const { key } = payload;
  const response = await axiosInstance.get(`/register?step=${key}`);
  return response?.data;
};

export const postCertificateStepService = async (payload: any) => {
  if (!payload) return;
  const response = await axiosInstance('/register', {
    method: 'PUT',
    data: {
      ...payload
    }
  });
  return response?.data;
};
export const deleteCertificateStepService = async (payload: any) => {
  if (!payload) return;
  const url = payload?.step ? `/register?step=${payload?.step}` : '/register';
  const response = await axiosInstance.delete(url);
  return response?.data;
};
