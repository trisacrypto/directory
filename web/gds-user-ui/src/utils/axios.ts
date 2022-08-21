import axios from 'axios';
import { getCookie } from 'utils/cookies';
import { auth0CheckSession } from 'utils/auth0.helper';
import { setCookie } from './cookies';
const axiosInstance = axios.create({
  baseURL: process.env.REACT_APP_TRISA_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
});

axiosInstance.defaults.withCredentials = true;
// intercept request and check if token has expired or not
axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    // let _retry = 0;
    const originalRequest = error.config;
    //

    if (error && !error.response) {
      return Promise.reject<any>(new Error('Network connection error'));
    }
    // handle 403/401 error by regenerating a new token and retrying the request

    if (
      (error?.response?.status === 403 || error?.response?.status === 401) &&
      !originalRequest._retry
    ) {
      originalRequest._retry = true;
      const tokenPayload: any = await auth0CheckSession();
      console.log('tokenPayload', tokenPayload);
      const token = tokenPayload?.accessToken;
      if (token) {
        setCookie('access_token', tokenPayload.accessToken);
        setCookie('user_locale', tokenPayload?.idTokenPayload?.locale || 'en');
        const csrfToken = getCookie('csrf_token');
        axiosInstance.defaults.headers.common.Authorization = `Bearer ${token}`;
        axiosInstance.defaults.headers.common['X-CSRF-Token'] = csrfToken;

        return axiosInstance(originalRequest);
      }
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
