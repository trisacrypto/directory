import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { setCookie } from 'utils/cookies';
import { logUserInBff, getUserRoles } from 'modules/auth/login/auth.service';
import { t } from '@lingui/macro';
import {
  auth0SignIn,
  auth0SignUp,
  auth0SignWithSocial,
  auth0Hash,
  auth0CheckSession
} from 'utils/auth0.helper';
import { handleError } from 'utils/utils';

export const userLoginWithSocial = (social: string) => {
  if (social === 'google') {
    auth0SignWithSocial('google-oauth2');
  }
};
export const userLoginWithEmail = createAsyncThunk(
  'users/userLoginWithEmail',
  async (data: any, thunkAPI) => {
    try {
      await auth0SignIn({
        username: data.username,
        password: data.password,
        responseType: 'token id_token',
        realm: 'Username-Password-Authentication'
      });
    } catch (err: any) {
      // handleError(e);
      thunkAPI.rejectWithValue(err.response.data);
    }
  }
);
export const userSignupWithEmail = createAsyncThunk(
  'users/userSignUpWithEmail',
  async (data: any, thunkAPI) => {
    try {
      const response: any = await auth0SignUp({
        email: data.username,
        password: data.password,
        connection: 'Username-Password-Authentication'
      });
      if (response) {
        if (!response.emailVerified) {
          return thunkAPI.rejectWithValue(response);
        }
      }
    } catch (err: any) {
      return thunkAPI.rejectWithValue(err.response.data);
    }
  }
);

export const getAuth0User: any = createAsyncThunk(
  'users/getuser',
  async (hasToken: boolean, thunkAPI) => {
    try {
      // then login with auth0
      const getUserInfo: any = hasToken && (await auth0Hash());
      // console.log('[getUserInfo]', getUserInfo);
      const updatedTime = new Date(getUserInfo?.idTokenPayload?.updated_at).getTime() / 1000;
      const expiresTime = updatedTime + getUserInfo.expiresIn;
      setCookie('access_token', getUserInfo?.accessToken);
      setCookie('user_locale', getUserInfo?.idTokenPayload?.locale || 'en');
      setCookie('expires_in', expiresTime);
      if (getUserInfo && getUserInfo?.idTokenPayload?.email_verified) {
        const getUser = await logUserInBff();
        // get user roles
        const getRoles = await getUserRoles() as any;
        // console.log('[getUser]', getUser);
        // check if user response contains refresh_token flag
        if (getUser?.data?.refresh_token) {
          // refresh token
          const newUserPayload: any = await auth0CheckSession();
          // get user info data
          // console.log('[newUserPayload]', newUserPayload);
          setCookie('access_token', newUserPayload?.accessToken);
          setCookie('user_locale', newUserPayload?.idTokenPayload?.locale || 'en');
          // set expired time
          const updated = new Date(newUserPayload?.idTokenPayload?.updated_at).getTime() / 1000;
          const expires = updated + getUserInfo.expiresIn;
          setCookie('expires_in', expires);
          const userInfo: TUser = {
            isLoggedIn: true,
            user: {
              name: newUserPayload?.idTokenPayload?.name,
              pictureUrl: newUserPayload?.idTokenPayload?.picture,
              email: newUserPayload?.idTokenPayload?.email,
              roles: getRoles?.data?.roles
            }
          };
          return userInfo;
        }
        // return;
        if (getUser.status === 204) {
          const userInfo: TUser = {
            isLoggedIn: true,
            user: {
              name: getUserInfo?.idTokenPayload?.name,
              pictureUrl: getUserInfo?.idTokenPayload?.picture,
              email: getUserInfo?.idTokenPayload?.email,
              roles: getRoles?.data?.roles
            }
          };
          return userInfo;
        } else {
          return thunkAPI.rejectWithValue(t`Something went wrong. Please try again later.`);
        }
      } else {
        return thunkAPI.rejectWithValue(
          t`Your account has not been verified. Please check your email to verify your account.`
        );
      }
    } catch (err: any) {
      handleError(err, '[getAuth0User] failed to get user');
      return thunkAPI.rejectWithValue(err.response.data);
    }
  }
);

export const initialValue: TUser = {
  isLoggedIn: false,
  isFetching: false,
  isError: false,
  errorMessage: '',
  user: null
};

const userSlice: any = createSlice({
  name: 'user',
  initialState: initialValue,
  reducers: {
    login: (state: any, { payload }: any) => {
      state.user = payload.user;
      state.isLoggedIn = true;
    },
    logout: (state: any) => {
      state.isError = false;
      state.isLoggedIn = false;
      state.isFetching = false;
      state.user = null;

      return state;
    }
    // isloading: (state: any, { payload }: any) => {
  },
  extraReducers: {
    [getAuth0User.fulfilled]: (state, { payload }) => {
      // console.log('[getAuth0User.fulfilled]', payload);
      state.isFetching = false;
      state.isLoggedIn = true;
      state.user = payload.user;
    },
    [getAuth0User.pending]: (state) => {
      // console.log('[getAuth0User.pending]', state);
      state.isFetching = true;
    },
    [getAuth0User.rejected]: (state, { payload }) => {
      // console.log('[getAuth0User.rejected]', payload);
      state.isFetching = false;
      state.isError = true;
      state.errorMessage = payload?.error ? payload.error : payload;
    }
  }
});

export const userReducer = userSlice.reducer;
export const { login, logout } = userSlice.actions;
// selectors
export const userSelector = (state: any) => state.user;
export const isLoggedInSelector = (state: any) => state.user.isLoggedIn;
