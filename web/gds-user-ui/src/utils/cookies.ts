import Cookies from 'universal-cookie';

const cookie = new Cookies();

export const setCookie = (key: string, value: any, path = '/') => {
  cookie.set(key, value, { path });
};

export const getCookie = (key: string) => {
  return cookie.get(key);
};

export const removeCookie = (key: string, path?: string) => {
  cookie.remove(key, { path: path || '/' });
};
