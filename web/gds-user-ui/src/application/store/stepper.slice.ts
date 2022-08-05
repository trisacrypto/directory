import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { clear } from 'console';
import { loadStepperFromLocalStorage } from 'utils/localStorageHelper';
import { loadDefaultStepperSync } from 'modules/dashboard/registration/utils';
export type TStep = {
  status: string;
  key?: number;
};
export type TPayload = {
  currentStep: number | string;
  steps: TStep[];
  lastStep: number | null;
  hasReachSubmitStep?: boolean;
};
export const initialValue: TPayload = {
  currentStep: 1,
  steps: [
    {
      status: 'progress'
    }
  ],
  lastStep: null,
  hasReachSubmitStep: false
};

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
    setHasReachSubmitStep: (state: any, { payload }: any) => {
      state.hasReachSubmitStep = payload.hasReachSubmitStep;
    },
    setLastStep: (state: any, { payload }: any) => {
      state.lastStep = payload.lastStep;
    },
    setStepFormValue: (state: any, { payload }: any) => {
      state.steps.map((step: any) => {
        if (step.key === payload.step && state.currentStep) {
          step.data = payload.formValues;
        }
      });
    },
    getCurrentFormValues: (state: any, { payload }: any | null) => {
      const found = state.steps.filter(
        (step: any) => step.key === payload?.step || state.currentStep
      );
      if (found.length === 1) {
        return found[0].data;
      }
      return null;
    },
    setSubmitStep: (state: any, { payload }: any) => {
      state.hasReachSubmitStep = payload.submitStep;
    },
    // set initial value
    setInitialValue: (state: TPayload, { payload }: any) => {
      state.currentStep = payload.currentStep;
      state.steps = payload.steps;
      state.lastStep = payload.lastStep;
      state.hasReachSubmitStep = payload.hasReachSubmitStep;
    },
    // get current state
    getCurrentState: (state: TPayload) => {
      return state;
    },
    clearStepper: (state: any) => {
      state.steps = [
        {
          key: 1,
          status: 'progress'
        }
      ];
      state.currentStep = 1;
      state.lastStep = null;
      state.hasReachSubmitStep = false;
    }
  }
});

export const stepperReducer = stepperSlice.reducer;
export const {
  addStep,
  setCurrentStep,
  setStepStatus,
  setLastStep,
  setStepFormValue,
  getCurrentFormValues,
  setSubmitStep,
  clearStepper,
  setHasReachSubmitStep,
  setInitialValue,
  getCurrentState
} = stepperSlice.actions;
