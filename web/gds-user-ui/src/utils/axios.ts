/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable no-lonely-if */
import axios, { AxiosError } from 'axios';
import { getCookie } from 'utils/cookies';
import { auth0CheckSession } from 'utils/auth0.helper';
import { setCookie, clearCookies } from './cookies';
import { createStandaloneToast } from '@chakra-ui/react';

const toast = createStandaloneToast();

const axiosInstance = axios.create({
  baseURL: process.env.REACT_APP_TRISA_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
});

axiosInstance.defaults.withCredentials = true;
// intercept request and check if token has expired or not
axiosInstance.interceptors.request.use(
  async (config: any) => {
    const token = getCookie('access_token');
    const csrfToken = getCookie('csrf_token');
    if (token) {
      const { exp } = JSON.parse(atob(token.split('.')[1]));
      const isExpired = exp * 1000 < Date.now();
      if (isExpired) {
        const { accessToken } = (await auth0CheckSession()) as any;
        setCookie('token', accessToken);
        config.headers.Authorization = `Bearer ${accessToken}`;
      } else {
        config.headers.Authorization = `Bearer ${token}`;
      }
      if (csrfToken) {
        config.headers['X-CSRF-Token'] = csrfToken;
      }
    }
    return config;
  },
  (error) => {
    Promise.reject(error);
  }
);

axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error: any) => {
    // let _retry = 0;
    const originalRequest = error?.config;
    originalRequest._retry = originalRequest?._retry || 0;
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
        // The user should not be logged out if he can't process a request due to a missing permission
        if (
          error?.response?.status === 401 &&
          (error as AxiosError).config.url !== '/users/login'
        ) {
          toast({
            title: "Sorry, you don't have permission to perform this action",
            status: 'error',
            position: 'top-right',
            isClosable: true
          });
          return;
        } else {
          // remove cookies and local storage and redirect to login page

          clearCookies();
          localStorage.removeItem('trs_stepper');
          localStorage.removeItem('persist:root');
          window.location.href = `/auth/login?q=token_expired`;
        }

        clearCookies();
        switch (error.response.status) {
          // case 401:
          //   console.log('[error] response', error.response);
          //   if ((error as AxiosError).config.url !== '/users/login') {
          //     toast({
          //       title: "Sorry, you don't have permission to perform this action",
          //       status: 'error',
          //       position: 'top-right'
          //     });
          //   } else {
          //     window.location.href = `/auth/login?q=token_expired`;
          //   }
          //   break;
          case 403:
            window.location.href = `/auth/login?q=unauthorized`;
            break;
          case 503:
            window.location.href = `/maintenance`;
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
