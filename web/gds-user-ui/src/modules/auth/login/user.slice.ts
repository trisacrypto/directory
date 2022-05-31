import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { loadStepperFromLocalStorage } from 'utils/localStorageHelper';
import { setCookie, getCookie, removeCookie } from 'utils/cookies';

const hasSession = getCookie('session');
export const initialValue: TUser = hasSession
  ? { isLoggedIn: true, user: hasSession }
  : { isLoggedIn: false, user: null };

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
  }
});

export const userReducer = userSlice.reducer;
export const { login, logout, isUserAuthenticated } = userSlice.actions;
// selectors
export const userSelector = (state: any) => state.user.user;
export const isLoggedInSelector = (state: any) => state.user.isLoggedIn;
