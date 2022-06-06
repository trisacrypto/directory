import axios from 'axios';
import { getRefreshToken } from 'utils/utils';
import { getCookie, setCookie } from 'utils/cookies';
const axiosInstance = axios.create({
  baseURL: process.env.REACT_APP_TRISA_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
});
axiosInstance.defaults.withCredentials = true;
axiosInstance.interceptors.request.use(
  (response) => {
    return response;
  },
  async (error) => {
    const originalRequest = error.config;
    // retry 3 time if request failed

    // Reject promise for now if usual error
    if (error.response.status !== 401) {
      return Promise.reject(error);
    }

    if (error.response.status === 403) {
      console.log('403');
    }

    if (error.response.status === 401 && error.response.data.error === 'Unauthorized') {
      const token = await getRefreshToken();
      if (token) {
        setCookie('access_token', token);
        axiosInstance.defaults.headers.common.Authorization = `Bearer ${token}`;
        return axiosInstance.request(originalRequest);
      }
    }
  }
);

export default axiosInstance;
