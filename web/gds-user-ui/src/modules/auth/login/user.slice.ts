import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { setCookie } from 'utils/cookies';
import { logUserInBff } from 'modules/auth/login/auth.service';
import { t } from '@lingui/macro';
import { auth0SignIn, auth0SignUp, auth0SignWithSocial, auth0Hash } from 'utils/auth0.helper';

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
      const getUserInfo: any = hasToken && (await auth0Hash());

      if (getUserInfo && getUserInfo?.idTokenPayload.email_verified) {
        setCookie('access_token', hasToken);
        setCookie('user_locale', getUserInfo?.locale);
        const getUser = await logUserInBff();

        if (getUser.status === 204) {
          const userInfo: TUser = {
            isLoggedIn: true,
            user: {
              name: getUserInfo?.idTokenPayload?.name,
              pictureUrl: getUserInfo?.idTokenPayload?.picture,
              email: getUserInfo?.idTokenPayload?.email
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
  },
  extraReducers: {
    [getAuth0User.fulfilled]: (state, { payload }) => {
      state.isFetching = false;
      state.isLoggedIn = true;
      state.user = payload.user;
    },
    [getAuth0User.pending]: (state) => {
      state.isFetching = true;
    },
    [getAuth0User.rejected]: (state, { payload }) => {
      state.isFetching = false;
      state.isError = true;
      state.errorMessage = payload;
    }
  }
});

export const userReducer = userSlice.reducer;
export const { login, logout } = userSlice.actions;
// selectors
export const userSelector = (state: any) => state.user;
export const isLoggedInSelector = (state: any) => state.user.isLoggedIn;
