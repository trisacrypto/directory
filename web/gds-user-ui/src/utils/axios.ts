import axios from 'axios';
import { getCookie } from 'utils/cookies';
import { auth0CheckSession } from 'utils/auth0.helper';
import { setCookie, clearCookies } from './cookies';
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
    if (response.status === 200) {
    }

    return response;
  },
  async (error) => {
    // let _retry = 0;
    const originalRequest = error.config;
    originalRequest._retry = originalRequest._retry || 0;
    //

    if (error && !error.response) {
      return Promise.reject<any>(new Error('Network connection error'));
    }
    // handle 403/401 error by regenerating a new token and retrying the request

    if (error?.response?.status === 403 || error?.response?.status === 401) {
      // retry the request 1 time

      if (originalRequest._retry < 1) {
        const tokenPayload: any = await auth0CheckSession();
        const token = tokenPayload?.accessToken;
        if (token) {
          setCookie('access_token', tokenPayload.accessToken);
          setCookie('user_locale', tokenPayload?.idTokenPayload?.locale || 'en');
          const csrfToken = getCookie('csrf_token');
          axiosInstance.defaults.headers.common.Authorization = `Bearer ${token}`;
          axiosInstance.defaults.headers.common['X-CSRF-Token'] = csrfToken;
          originalRequest._retry += 1;
          return axiosInstance(originalRequest);
        }
      } else {
        clearCookies();
        switch (error.response.status) {
          case '401':
            window.location.href = `/auth/login?q=token_expired`;
            break;
          case '403':
            window.location.href = `/auth/login?q=unauthorized`;
            break;
          default:
            window.location.href = `/auth/login?error_description=${error.response.data.error}`;
        }
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
