import { combineReducers } from '@reduxjs/toolkit';
import { stepperReducer } from './stepper.slice';
const rootReducer = combineReducers({
  stepper: stepperReducer
});

export type RootState = ReturnType<typeof rootReducer>;

export default rootReducer;
