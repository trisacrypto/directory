import { combineReducers } from '@reduxjs/toolkit';
import { stepperReducer } from './stepper.slice';
import { userReducer } from 'modules/auth/login/user.slice';
import { collaboratorReducer } from 'modules/dashboard/collaborator/collaborator.slice';

const rootReducer = combineReducers({
  stepper: stepperReducer,
  user: userReducer,
  collaborators: collaboratorReducer
});

export type RootState = ReturnType<typeof rootReducer>;

export default rootReducer;
