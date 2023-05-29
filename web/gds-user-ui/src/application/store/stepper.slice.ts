/* eslint-disable @typescript-eslint/no-use-before-define */
import { createSlice } from '@reduxjs/toolkit';

export type TStep = {
  status: string;
  key?: number;
};

export const initialValue: TPayload = {
  currentStep: 1,
  steps: [
    {
      key: 1,
      status: 'progress',
      isDirty: false
    }
  ],
  lastStep: null,
  hasReachSubmitStep: false,
  testnetSubmitted: false,
  hasReachReviewStep: false,
  mainnetSubmitted: false,
  status: {
    testnet: 'progress',
    mainnet: 'progress'
  },
  data: {}
};

const stepperSlice: any = createSlice({
  name: 'stepper',
  initialState: initialValue,
  reducers: {
    setCurrentStep: (state: any, { payload }: any) => {
      if (payload.currentStep === 6) {
        state.hasReachReviewStep = true;
      }
      state.currentStep = payload.currentStep;
    },
    incrementStep: (state: any) => {
      // always set isDirty to false when incrementing step
      if (state.currentStep) {
        state?.steps?.map((step: any) => {
          if (step.key === state.currentStep) {
            step.isDirty = false;
          }
        });
      }

      if (state.currentStep < 6) {
        state.currentStep += 1;
      }

      // if next step is not in the list, add it
      if (!state?.steps?.find((step: any) => step.key === state.currentStep)) {
        state?.steps?.push({
          key: state.currentStep,
          status: 'progress',
          isDirty: false
        });
      }
    },
    decrementStep: (state: any) => {
      if (state.currentStep) {
        state?.steps?.map((step: any) => {
          if (step.key === state.currentStep) {
            step.isDirty = false;
          }
        });
      }
      state.currentStep -= 1;
      // if current step is 6 then set hasReachSubmitStep to false
      // if (state.currentStep === 6 && state.hasReachSubmitStep) {
      //   state.hasReachSubmitStep = false;
      // }
    },
    addStep: (state: any, { payload }: any) => {
      // if step is not in the list, add it
      const payloadStep = payload?.step || state.currentStep;
      if (!state?.steps?.find((step: any) => step.key === payloadStep)) {
        state?.steps?.push({
          key: payloadStep,
          status: payload?.status || 'progress',
          isDirty: false
        });
      }
    },
    setStepStatus: (state: any, { payload }: any) => {
      state?.steps?.map((step: any) => {
        if (step.key === payload.step) {
          console.log('payload.status', payload.status);
          step.status = payload.status;
        }
      });
    },
    setStepMissingFields: (state: any, { payload }: any) => {
      state?.steps?.map((step: any) => {
        if (step.key === payload.step && state.currentStep) {
          step.missingFields = payload.errors;
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
      state.testnetSubmitted = payload.testnetSubmitted;
      state.mainnetSubmitted = payload.mainnetSubmitted;
      state.hasReachReviewStep = !!(payload.currentStep === 6); // default value
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
      state.testnetSubmitted = false;
      state.mainnetSubmitted = false;
      state.hasReachReviewStep = false;
      state.data = {};
    },
    // set testnet submission
    setTestnetSubmitted: (state: any, { payload }: any) => {
      state.testnetSubmitted = payload.testnetSubmitted;
    },
    // set mainnet submission
    setMainnetSubmitted: (state: any, { payload }: any) => {
      state.mainnetSubmitted = payload.mainnetSubmitted;
    },
    // set certificate data
    setCertificateValue: (state: any, { payload }: any) => {
      state.data = { ...payload.value };
    },

    // get certificate data
    getCertificateData: (state: any) => {
      return state.data;
    },
    setVaspName(state: any, { payload }: any) {
      state.data.organization_name = payload;
    },
    // this should help us open the popup when the user tries to jump to the step from the progress bar
    setIsDirty(state: any, { payload }: any) {
      const payloadStep = payload?.step || state.currentStep;

      state?.steps?.map((step: any) => {
        if (step.key === payloadStep && state.currentStep) {
          step.isDirty = payload.isDirty ?? !step.isDirty;
        }
      });
    },

    getIsDirty(state: any, { payload }: any) {
      const found = state.steps.filter(
        (step: any) => step.key === payload?.step || state.currentStep
      );
      if (found.length === 1) {
        return found[0].isDirty;
      }
      return null;
    }
  }
});

export const stepperReducer = stepperSlice.reducer;
export const {
  addStep,
  incrementStep,
  decrementStep,
  setCurrentStep,
  setStepStatus,
  setLastStep,
  setStepFormValue,
  getCurrentFormValues,
  setSubmitStep,
  clearStepper,
  setHasReachSubmitStep,
  setInitialValue,
  getCurrentState,
  setTestnetSubmitted,
  setMainnetSubmitted,
  setCertificateValue,
  getCertificateData,
  setVaspName,
  setStepMissingFields,
  setIsDirty,
  getIsDirty
} = stepperSlice.actions;
