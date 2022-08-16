import axios from 'axios';
import { getCookie } from 'utils/cookies';
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
  (error) => {
    // let _retry = 0;
    const originalRequest = error.config;
    console.log('[AxiosError]', error.response.status);

    if (error && !error.response) {
      return Promise.reject<any>(new Error('Network connection error'));
    }

    return Promise.reject(error);

    // }
  }
);

export const setAuthorization = () => {
  const token = getCookie('access_token');
  const csrfToken = getCookie('csrf_token');
  axiosInstance.defaults.headers.common.Authorization = `Bearer ${token}`;
  axiosInstance.defaults.headers.common['X-CSRF-Token'] = csrfToken;
};

export default axiosInstance;
