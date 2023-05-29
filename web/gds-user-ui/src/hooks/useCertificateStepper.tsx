/* eslint-disable @typescript-eslint/no-use-before-define */
import Store from 'application/store';
import { getCurrentStep } from 'application/store/selectors/stepper';
import { useDispatch, useSelector } from 'react-redux';
import {
  setCurrentStep,
  setStepStatus,
  clearStepper,
  setHasReachSubmitStep,
  setInitialValue,
  setTestnetSubmitted,
  setMainnetSubmitted,
  setCertificateValue,
  setStepMissingFields,
  incrementStep,
  decrementStep,
  setIsDirty,
  addStep
} from 'application/store/stepper.slice';
// import { getFieldNames } from 'utils/getFieldNames';
import { setRegistrationDefaultValue } from 'modules/dashboard/registration/utils';

import { LSTATUS } from 'components/RegistrationForm/CertificateStepLabel';

// 'TODO:' this hook should be improved to be more generic
const useCertificateStepper = () => {
  const dispatch = useDispatch();
  const currentStep: number = useSelector(getCurrentStep);
  console.log('[useCertificateStepper] currentStep', currentStep);

  const removeMissingFields = (steps: TStep[]) => {
    return steps.map((step: TStep) => {
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const { missingFields, isDirty, ...rest } = step;
      return rest;
    });
  };

  const currentState = () => {
    // log store state
    const updatedState = Store.getState().stepper;
    const formatState = {
      current: updatedState.currentStep,
      steps: removeMissingFields(updatedState.steps),
      ready_to_submit: updatedState.hasReachSubmitStep
    };
    return formatState;
  };

  const nextStep = (data?: any) => {
    const errorFields = data?.errors;
    console.log('[useCertificateStepper 1] errorFields', errorFields);
    if (data && errorFields && Object.keys(errorFields).length > 0) {
      console.log('[useCertificateStepper 2] errorFields', errorFields);
      dispatch(setStepStatus({ step: currentStep, status: LSTATUS.ERROR }));
    } else {
      // setInitialState(data?.form);
      dispatch(setStepStatus({ step: currentStep, status: LSTATUS.COMPLETE }));
    }

    dispatch(incrementStep());
  };

  const previousStep = (data?: any) => {
    const errorFields = data?.errors;
    if (data && errorFields) {
      dispatch(setStepMissingFields({ step: currentStep, errors: errorFields }));
      dispatch(setStepStatus({ step: currentStep, status: LSTATUS.ERROR }));
    } else {
      // setInitialState(data?.form);
      dispatch(setStepStatus({ step: currentStep, status: LSTATUS.COMPLETE }));
    }
    dispatch(decrementStep());
  };

  const jumpToStep = (step: number) => {
    dispatch(setHasReachSubmitStep({ hasReachSubmitStep: false }));
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

  // save the registration value to the store
  const setRegistrationValue = (value: any) => {
    dispatch(setCertificateValue({ value }));
  };

  const setInitialState = (value: any) => {
    const state: TPayload = {
      currentStep: value?.state?.current,
      steps: value?.state?.steps,
      lastStep: 6,
      hasReachSubmitStep: value?.state?.ready_to_submit || false,
      testnetSubmitted: value?.state?.testnetSubmitted || false,
      mainnetSubmitted: value?.state?.mainnetSubmitted || false
    };
    dispatch(setInitialValue(state));
  };

  const hasStepErrors = () => {
    const steps = Store.getState().stepper.steps;
    const found = steps.filter((step: any) => step.status === 'error');
    return found.length > 0;
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

  const clearCertificateStepper = () => {
    dispatch(clearStepper());
  };

  const updateHasReachSubmitStep = (hasReachSubmitStep: boolean) => {
    dispatch(setHasReachSubmitStep({ hasReachSubmitStep }));
  };

  const updateIsDirty = (isDirty: boolean, step: number) => {
    dispatch(setIsDirty({ step, isDirty: !!isDirty }));
    // if (isDirty) {
    //   // eslint-disable-next-line @typescript-eslint/no-use-before-define
    //   addDefaultStep();
    // }
  };
  const getIsDirtyState = (step = currentStep) => {
    const steps = Store.getState().stepper.steps;
    const found = steps.filter((s: TStep) => s.key === step && s.isDirty === true);
    return found.length > 0;
  };

  const addDefaultStep = (step?: number) => {
    const payload = {
      step: step || currentStep,
      status: LSTATUS.PROGRESS
    };
    dispatch(addStep(payload));
  };

  const updateStepStatusToIncomplete = () => {
    // check if the current step has progress status and is not first step
    if (currentStep !== 1) {
      const steps = Store.getState().stepper.steps;
      const found = steps.filter(
        (s: TStep) => s.key === currentStep && s.status === LSTATUS.PROGRESS
      );
      if (found.length > 1) {
        dispatch(setStepStatus({ step: currentStep, status: LSTATUS.INCOMPLETE }));
      }
    }
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
    updateStateFromFormValues,
    setRegistrationValue,
    clearCertificateStepper,
    hasStepErrors,
    updateHasReachSubmitStep,
    updateIsDirty,
    getIsDirtyState,
    addDefaultStep,
    updateStepStatusToIncomplete
  };
};

export default useCertificateStepper;
