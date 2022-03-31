import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { loadStepperFromLocalStorage } from 'utils/localStorageHelper';
export type TStep = {
  status: string;
  key?: number;
};
export type TPayload = {
  currentStep: number | string;
  steps: TStep[];
  lastStep: number | null;
};
export const initialValue: TPayload = loadStepperFromLocalStorage();

const stepperSlice: any = createSlice({
  name: 'stepper',
  initialState: initialValue,
  reducers: {
    setCurrentStep: (state: any, { payload }: any) => {
      state.currentStep = payload.currentStep;
    },
    addStep: (state: any, { payload }: any) => {
      state.steps.push(payload);
    },
    setStepStatus: (state: any, { payload }: any) => {
      state.steps.map((step: any) => {
        if (step.key === payload.step && state.currentStep) {
          step.status = payload.status;
        }
      });
    },
    setLastStep: (state: any, { payload }: any) => {
      state.lastStep = payload.lastStep;
    }
  }
});

export const stepperReducer = stepperSlice.reducer;
export const { addStep, setCurrentStep, setStepStatus, setLastStep } = stepperSlice.actions;
