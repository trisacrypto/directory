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
    if (error && !error.response) {
      return Promise.reject<any>(
        new Error('Sorry we cannot reach the server, please contact the admin')
      );
    }
    // handle 403 error

    if (error && error.response && error.response.status === 401) {
      console.log('[401]');

      originalRequest.retry = originalRequest.retry || 0;

      if (originalRequest.retry <= 1) {
        console.log('[401]-retry]', originalRequest.retry);
        originalRequest.retry += 1;
        const token = await getRefreshToken();

        if (token) {
          setCookie('access_token', token);
          axiosInstance.defaults.headers.common.Authorization = `Bearer ${token}`;
          return axiosInstance(originalRequest);
        }
      }
    }
    if (error && error?.response?.status === 403) {
      return Promise.reject<any>(new Error('Unauthorize user'));
    }
  }
);

export default axiosInstance;
