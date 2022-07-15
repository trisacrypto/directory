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
axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    console.log('[AxiosError detected]', error.response.status);
    const originalRequest = error.config;

    if (error && error.response && error.response.status === 401) {
      console.log('[401]');

      originalRequest.retry = originalRequest.retry || 0;

      if (originalRequest.retry <= 1) {
        console.log('[401]-retry]', originalRequest.retry);
        originalRequest.retry += 1;
        const token = await getRefreshToken();
        // if (token) {
        //   const headers = {
        //     'Content-Type': 'application/json',
        //     Authorization: `Bearer ${token}`
        //   };
        //   const newRequest = {
        //     ...originalRequest,
        //     headers,
        //     url: `${originalRequest.url}?${originalRequest.data}`
        //   };
        //   setCookie('access_token', token);
        //   return axiosInstance.request(newRequest);
        // }
        if (token) {
          setCookie('access_token', token);
          axiosInstance.defaults.headers.common.Authorization = `Bearer ${token}`;
          return axiosInstance.request(originalRequest);
        }
      }
    }
    if (error && !error.response) {
      return Promise.reject<any>(new Error('Network connection error'));
    }
  }
);

export default axiosInstance;
