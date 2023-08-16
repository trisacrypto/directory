import axiosInstance from 'utils/axios';
import { getCookie } from 'utils/cookies';
export const getOverview = async () => {
  const response = await axiosInstance.get(`/overview`, {
    headers: {
      Authorization: `Bearer ${getCookie('access_token')}`
    }
  });
  return response;
};
export const getAnnouncementsData = async () => {
  const response = await axiosInstance.get(`/announcements`, {
    headers: {
      Authorization: `Bearer ${getCookie('access_token')}`
    }
  });
  return response;
};
