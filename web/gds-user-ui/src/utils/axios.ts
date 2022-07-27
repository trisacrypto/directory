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
    console.log('[AxiosError]', error);

    if (error && !error.response) {
      return Promise.reject<any>(new Error('Network connection error'));
    }
    const originalRequest = error.config;
    // retry 3 time if request failed

    if (error.response.status === 403) {
      return Promise.reject(error);
    }

    if (error.response.status === 401 && error.response.data.error === 'Unauthorized') {
      console.log('[TokenError]', error);
      // if (originalRequest.retry < 3) {
      //   originalRequest.retry = originalRequest.retry || 0;
      //   originalRequest.retry += 1;
      const token = await getRefreshToken();
      console.log(token);
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
      // const token = await getRefreshToken();
      if (token) {
        setCookie('access_token', token);
        axiosInstance.defaults.headers.common.Authorization = `Bearer ${token}`;
        return axiosInstance.request(originalRequest);
      }
      // }
    }
  }
);

export const setAuthorization = () => {
  const token = getCookie('access_token');
  const csrfToken = getCookie('csrf_token');
  axiosInstance.defaults.headers.common.Authorization = `Bearer ${token}`;
  axiosInstance.defaults.headers.common['X-CSRF-Token'] = csrfToken;
};

export default axiosInstance;
