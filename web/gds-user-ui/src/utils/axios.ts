import axios from 'axios';
import { getRefreshToken } from 'utils/utils';
import { getCookie, setCookie, removeCookie, clearCookies } from 'utils/cookies';
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
    // if (error.response.status === 401 || error.response.status === 403) {
    //   clearCookies();
    //   // get origin url
    //   const origin = window.location.origin;
    //   if (error.response.status === 401) {
    //     window.location.href = `/auth/login?from=${origin}&q=unauthorized`;
    //   }
    //   if (error.response.status === 403) {
    //     window.location.href = `/auth/login?from=${origin}&q=token_expired`;
    //   }
    //   return Promise.reject<any>(new Error('Unauthorized'));
    // }

    // // if 403 detected [unauthorize issue],clear cookies , reauthenticate and retry the request once again
    // // if (error?.response?.status === 403 || error?.response?.status === 401) {
    // //   if (_retry === 0) {
    // //     console.log('[AxiosError]', error.response.status);
    // //     removeCookie('access_token');
    // //     clearCookies();
    // //     const token: any = await getRefreshToken();
    // //     if (token) {
    // //       // set token to axios header
    // //       const getToken = token.accessToken;
    // //       axiosInstance.defaults.headers.common.Authorization = `Bearer ${getToken}`;
    // //       // set token to cookie
    // //       setCookie('access_token', getToken);
    // //       // retry the request
    // //       _retry++;
    // //       return axiosInstance(originalRequest);
    // //     }
    // //   } else {
    // //     // clean and redirect to login page
    // //     clearCookies();
    // //     window.location.href = 'auth/login';
    // //   }
    // // }

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
