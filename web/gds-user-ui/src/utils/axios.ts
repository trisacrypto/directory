import axios from 'axios';
import { getRefreshToken } from 'utils/utils';

import Cookies from 'universal-cookie';

const cookies = new Cookies();
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
  (error) => {
    const originalRequest = error.config;
    // retry 3 time if request failed

    // Reject promise for now if usual error
    if (error.response.status !== 401) {
      return Promise.reject(error);
    }

    if (error.response.status === 401 && error.response.data.error === 'Unauthorized') {
      const token = getRefreshToken();
      if (token) {
        cookies.set('token', token, { path: '/' });
        axiosInstance.defaults.headers.common.Authorization = `Bearer ${token}`;
        return axiosInstance.request(error.config);
      }
    }
  }
);

export default axiosInstance;
