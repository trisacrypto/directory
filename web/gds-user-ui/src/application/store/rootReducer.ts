import { combineReducers } from '@reduxjs/toolkit';
import { stepperReducer } from './stepper.slice';
import { userReducer } from 'modules/auth/login/user.slice';
import storageSession from 'reduxjs-toolkit-persist/lib/storage/session';
import { persistReducer, persistStore } from 'redux-persist';

const userPersistConfig = {
  key: 'root',
  storage: storageSession
};

const rootReducer = combineReducers({
  stepper: stepperReducer,
  user: userReducer
});

export type RootState = ReturnType<typeof rootReducer>;

export default rootReducer;
