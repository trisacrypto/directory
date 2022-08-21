import React, { useState, useEffect } from 'react';
import { getSteps, getLastStep, resetStepper } from '../application/store/selectors/stepper';
import Store from 'application/store';
import { getCurrentStep } from 'application/store/selectors/stepper';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import {
  addStep,
  setCurrentStep,
  setStepStatus,
  TStep,
  setStepFormValue,
  setSubmitStep,
  clearStepper,
  setHasReachSubmitStep,
  setInitialValue,
  setTestnetSubmitted,
  setMainnetSubmitted
} from 'application/store/stepper.slice';
import {
  setRegistrationDefaultValue,
  postRegistrationValue
} from 'modules/dashboard/registration/utils';
import { findStepKey } from 'utils/utils';
import { LSTATUS } from 'components/TestnetProgress/CertificateStepLabel';
import { hasStepError } from '../utils/utils';

interface TState {
  status?: boolean;
  isMissed?: boolean;
  step?: number;
  errors?: any;
  isFormCompleted?: boolean;
  formValues?: any;
  values?: any;
  registrationValues?: any;
  setRegistrationState?: any;
  isDirty?: boolean;
}

// 'TODO:' this hook should be improve

const useCertificateStepper = () => {
  const dispatch = useDispatch();
  const currentStep: number = useSelector(getCurrentStep);
  const steps: TStep[] = useSelector(getSteps);
  const lastStep: number = useSelector(getLastStep);

  // get store state after dispatch action

  const currentState = () => {
    // log store state
    const updatedState = Store.getState().stepper;
    const formatState = {
      current: updatedState.currentStep,
      steps: updatedState.steps,
      ready_to_submit: updatedState.hasReachSubmitStep
    };
    return formatState;
  };

  // save form data to trtl if field is dirty

  const saveFormValue = async (stateData: any, formValue: any, setState?: any) => {
    // form value
    const _mergedData = { ...stateData, ...formValue };

    const getRegistrationData = () => {
      return {
        ..._mergedData,
        state: {
          ...currentState()
        }
      };
    };
    // save the form value if fields changed

    await postRegistrationValue(getRegistrationData());
    if (setState) {
      setState(getRegistrationData());
    }
  };

  const nextStep = async (state?: TState) => {
    const { values: formValues, registrationValues, setRegistrationState, isDirty } = state || {};

    // only for status update
    if (state?.isFormCompleted || !state?.errors) {
      dispatch(setStepStatus({ status: LSTATUS.COMPLETE, step: currentStep }));
    }
    // if we got an error that means require element are not completed
    if (state?.errors) {
      dispatch(setStepStatus({ status: LSTATUS.ERROR, step: currentStep }));
    }
    // allow manually setting the step status
    if (state?.status) {
      const found = findStepKey(steps, currentStep);
      if (found.length === 1) {
        dispatch(setStepStatus(state));
        dispatch(setCurrentStep({ currentStep: currentStep + 1 }));
        const foundNext = findStepKey(steps, currentStep + 1);
        if (foundNext.length === 0) {
          if (currentStep === lastStep) {
            return;
          }
          dispatch(addStep({ key: currentStep + 1, status: LSTATUS.PROGRESS }));
        }
      } else {
        if (currentStep === lastStep && state.isFormCompleted) {
          // that mean we move to submit step

          dispatch(setSubmitStep({ submitStep: true }));
        }
        dispatch(addStep({ key: currentStep, status: state.status }));
        dispatch(setCurrentStep({ currentStep: currentStep + 1 }));
      }
    }
    // if we reach the last step (here review step) , we need to set the submit step
    if (currentStep === lastStep) {
      // that mean we move to submit step
      if (!hasStepError(steps)) {
        dispatch(setSubmitStep({ submitStep: true }));
        dispatch(setCurrentStep({ currentStep: lastStep }));
      }
    } else {
      const found = findStepKey(steps, currentStep + 1);

      if (found.length === 0) {
        dispatch(setCurrentStep({ currentStep: currentStep + 1 }));
        dispatch(addStep({ key: currentStep + 1, status: LSTATUS.PROGRESS }));
      } else {
        if (found[0].status === LSTATUS.INCOMPLETE) {
          dispatch(setStepStatus({ step: currentStep + 1, status: LSTATUS.PROGRESS }));
        }
        dispatch(setCurrentStep({ currentStep: currentStep + 1 }));
      }
    }
    // save the form value if fields changed
    if (isDirty) {
      await saveFormValue(registrationValues, formValues, setRegistrationState);
    }
  };
  const previousStep = async (state?: TState) => {
    const { values: formValues, registrationValues, isDirty } = state || {};
    // if form value is set then save it to the dedicated step
    if (state?.formValues) {
      dispatch(setStepFormValue({ step: currentStep, formValues: state?.formValues }));
    }
    // do not allow to go back for the first step
    const step = currentStep;
    if (currentStep === 1) {
      return;
    }
    dispatch(setCurrentStep({ currentStep: step - 1 }));

    // if the current status is already completed, do not change the status

    const found = findStepKey(steps, currentStep);
    if (found.length > 0 && found[0].status !== LSTATUS.COMPLETE) {
      dispatch(setStepStatus({ step, status: LSTATUS.PROGRESS }));
    }
    // save the form value if fields changed
    if (isDirty) {
      await saveFormValue(registrationValues, formValues);
    }
  };

  const jumpToStep = (step: number) => {
    dispatch(setCurrentStep({ currentStep: step }));
  };

  const jumpToLastStep = () => {
    dispatch(setHasReachSubmitStep({ hasReachSubmitStep: false }));
  };

  const resetForm = () => {
    setRegistrationDefaultValue();
    dispatch(clearStepper());
  };

  // testnet submission state
  const testnetSubmissionState = () => {
    dispatch(setTestnetSubmitted({ testnetSubmitted: true }));
  };
  // mainnet submission state
  const mainnetSubmissionState = () => {
    dispatch(setMainnetSubmitted({ mainnetSubmitted: true }));
  };

  const setInitialState = (value: any) => {
    const state: TPayload = {
      currentStep: value.currentStep,
      steps: value.steps,
      lastStep: 6,
      hasReachSubmitStep: value.hasReachSubmitStep,
      testnetSubmitted: value.testnetSubmitted,
      mainnetSubmitted: value.mainnetSubmitted
    };
    dispatch(setInitialValue(state));
  };

  // update state from form values
  const updateStateFromFormValues = (values: any) => {
    const state: TPayload = {
      currentStep: values.current,
      steps: values.steps,
      lastStep: 6,
      hasReachSubmitStep: values.ready_to_submit,
      testnetSubmitted: false,
      mainnetSubmitted: false
    };
    dispatch(setInitialValue(state));
  };

  return {
    nextStep,
    previousStep,
    jumpToStep,
    resetForm,
    jumpToLastStep,
    setInitialState,
    currentState,
    testnetSubmissionState,
    mainnetSubmissionState,
    updateStateFromFormValues
  };
};

export default useCertificateStepper;
