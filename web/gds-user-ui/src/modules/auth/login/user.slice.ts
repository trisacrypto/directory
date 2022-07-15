import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { loadStepperFromLocalStorage } from 'utils/localStorageHelper';
import { setCookie, getCookie, removeCookie } from 'utils/cookies';

const hasSession = getCookie('session');
export const initialValue: TUser = hasSession
  ? {
      isLoggedIn: true,
      isFetching: false,
      isSuccess: false,
      isError: false,
      errorMessage: '',
      user: hasSession
    }
  : { isLoggedIn: false, isFetching: false, isSuccess: false, isError: false, user: null };

const userSlice: any = createSlice({
  name: 'user',
  initialState: initialValue,
  reducers: {
    login: (state: any, { payload }: any) => {
      state.user = payload.user;
      state.isLoggedIn = true;
    },
    logout: (state: any) => {
      state.user = null;
      state.isLoggedIn = false;
    }
    // isloading: (state: any, { payload }: any) => {
  }
});

export const userReducer = userSlice.reducer;
export const { login, logout, isUserAuthenticated } = userSlice.actions;
// selectors
export const userSelector = (state: any) => state.user.user;
export const isLoggedInSelector = (state: any) => state.user.isLoggedIn;
