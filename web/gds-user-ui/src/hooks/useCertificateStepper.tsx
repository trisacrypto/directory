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
  incrementStep,
  decrementStep,
  setIsDirty,
  addStep,
  setDeletedSteps,
  setDeletedStepValue
} from 'application/store/stepper.slice';
// import { getFieldNames } from 'utils/getFieldNames';
import { setRegistrationDefaultValue } from 'modules/dashboard/registration/utils';

import { LSTATUS } from 'components/RegistrationForm/CertificateStepLabel';
import { getStepNumber } from 'components/BasicDetailsForm/util';

// 'TODO:' this hook should be improved to be more generic
const useCertificateStepper = () => {
  const dispatch = useDispatch();
  const currentStep: number = useSelector(getCurrentStep);

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
    console.log('updatedState', updatedState);
    const formatState = {
      current: updatedState.currentStep,
      steps: removeMissingFields(updatedState.steps),
      ready_to_submit: updatedState.hasReachSubmitStep
    };
    return formatState;
  };

  const nextStep = (data?: any) => {
    const errorFields = data?.errors;
    const stepNumber = getStepNumber(data?.step) || currentStep;

    if (data && errorFields && Object.keys(errorFields).length > 0) {
      dispatch(setStepStatus({ step: stepNumber, status: LSTATUS.ERROR }));
    } else {
      // setInitialState(data?.form);
      dispatch(setStepStatus({ step: stepNumber, status: LSTATUS.COMPLETE }));
    }
    if (currentStep === 5) {
      localStorage.setItem('isFirstRender', 'true');
    }

    dispatch(incrementStep());
  };

  const previousStep = (data?: any) => {
    const errorFields = data?.errors;
    console.log('previous errorFields', errorFields);
    const stepNumber = getStepNumber(data?.step) || currentStep; // get step number from step name
    if (data && errorFields && Object.keys(errorFields).length > 0) {
      dispatch(setStepStatus({ step: stepNumber, status: LSTATUS.ERROR }));
    } else {
      // setInitialState(data?.form);
      dispatch(setStepStatus({ step: stepNumber, status: LSTATUS.COMPLETE }));
    }
    dispatch(decrementStep());
  };

  const jumpToStep = (step: number) => {
    if (currentStep === 5) {
      localStorage.setItem('isFirstRender', 'true');
    }
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
      mainnetSubmitted: value?.state?.mainnetSubmitted || false,
      deletedSteps: []
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
      mainnetSubmitted: false,
      deletedSteps: []
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
    // check if the current step is added to the steps array
    const steps = Store.getState().stepper.steps;
    const found = steps.filter((s: TStep) => s.key === step);
    if (!found.length) {
      addDefaultStep(step);
    }
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

  const addDefaultStep = (step?: number, status?: string) => {
    const payload = {
      step: step || currentStep,
      status: status || LSTATUS.PROGRESS
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

  const updateStepStatusToError = (step: number) => {
    dispatch(setStepStatus({ step, status: LSTATUS.ERROR }));
  };

  const updateStepStatusState = (payload: any) => {
    dispatch(setStepStatus(payload));
  };

  const updateCurrentStepState = (step: string) => {
    const s = getStepNumber(step);
    dispatch(setCurrentStep({ currentStep: s }));
  };

  const clearStepperState = () => {
    dispatch(clearStepper());
  };

  const getIsDirtyStateByStep = (stepName: string) => {
    const stepNumber = getStepNumber(stepName);
    const steps = Store.getState().stepper.steps;
    const found = steps.filter((s: TStep) => s.key === stepNumber && s.isDirty === true);
    return found.length > 0;
  };

  const updateDeleteStepState = (values: TDeleteStep) => {
    const { step, isDeleted } = values;
    const elm = Store.getState().stepper.deletedSteps;
    const found = elm.filter((s: TDeleteStep) => s.step === step);
    if (!found.length) {
      const payload = {
        step,
        isDeleted: true
      };
      dispatch(setDeletedSteps(payload));
    } else {
      // change the status to true
      dispatch(setDeletedStepValue({ step, isDeleted }));
    }
  };
  // get deleted step state
  const getDeletedStepState = (step: string) => {
    const elm = Store.getState().stepper.deletedSteps;
    return elm.filter((s: TDeleteStep) => s.step === step) || [];
  };

  const isStepDeleted = (step: string) => {
    const elm = Store.getState()?.stepper?.deletedSteps;
    if (!elm) return false;
    const found = elm?.filter((s: TDeleteStep) => s.step === step && s.isDeleted === true);
    return found?.length > 0;
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
    updateStepStatusToIncomplete,
    updateStepStatusState,
    updateCurrentStepState,
    clearStepperState,
    getIsDirtyStateByStep,
    updateDeleteStepState,
    getDeletedStepState,
    isStepDeleted,
    updateStepStatusToError,
    currentStep
  };
};

export default useCertificateStepper;
