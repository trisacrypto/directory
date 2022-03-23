import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export type TStep = {
  status: string;
  key?: number;
};
export type TPayload = {
  currentStep: number;
  steps: TStep[];
};
export const initialValue: TPayload = {
  currentStep: 1,
  steps: [
    {
      key: 1,
      status: 'progress'
    }
  ]
};

const stepperSlice: any = createSlice({
  name: 'stepper',
  initialState: initialValue,
  reducers: {
    setCurrentStep: (state: any, { payload }: PayloadAction<Partial<TPayload>>) => {
      state.currentStep = payload.currentStep;
    },
    setSteps: (state: any, { payload }: PayloadAction<Partial<TPayload>>) => {
      state.steps = [...state.steps, payload.steps];
    }
  }
});

export const stepperReducer = stepperSlice.reducer;
export const { setSteps, setCurrentStep } = stepperSlice.actions;
