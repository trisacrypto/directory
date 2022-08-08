import axios from 'axios';
import { getRefreshToken } from 'utils/utils';
import { getCookie, setCookie, removeCookie } from 'utils/cookies';
const axiosInstance = axios.create({
  baseURL: process.env.REACT_APP_TRISA_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
});
let message: any = '';
axiosInstance.defaults.withCredentials = true;
axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    console.log('[AxiosError]', error.response.status);

    if (error && !error.response) {
      return Promise.reject<any>(new Error('Network connection error'));
    }
    const originalRequest = error.config;
    // retry 3 time if request failed

    if (error.response.status === 401 && error.response.data.error === 'Unauthorized') {
      console.log('[TokenError]', error);
      // if (originalRequest.retry < 3) {
      //   originalRequest.retry = originalRequest.retry || 0;
      //   originalRequest.retry += 1;
      const token = await getRefreshToken();
      console.log(token);

      if (token) {
        setCookie('access_token', token);
        axiosInstance.defaults.headers.common.Authorization = `Bearer ${token}`;
        return axiosInstance.request(originalRequest);
      }
      // }
    }
    switch (error.response.status) {
      case 403:
        message = 'Session expired';
        removeCookie('access_token');
        removeCookie('trs_stepper');
        window.sessionStorage.clear();
        window.location.href = '/auth/login';
        break;
      case 404:
        message = error || 'Sorry! the data you are looking for could not be found';
        break;
      case 400:
        message = error;
        break;
      case 500:
        message = error || 'Something went wrong';
        break;
      default: {
        message =
          error.response && error.response.data
            ? error.response.data.message
            : error.message || error;
      }
    }

    return Promise.reject(message || error);

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
