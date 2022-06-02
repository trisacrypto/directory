import { combineReducers } from '@reduxjs/toolkit';
import { stepperReducer } from './stepper.slice';
import { userReducer } from 'modules/auth/login/user.slice';
const rootReducer = combineReducers({
  stepper: stepperReducer,
  user: userReducer
});

export type RootState = ReturnType<typeof rootReducer>;

export default rootReducer;
