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
  decrementStep
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
      const { missingFields, ...rest } = step;
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

  const nextStep = (errorFields?: any) => {
    console.log('[nextStep] errorFields', errorFields);
    if (errorFields && Object.keys(errorFields).length > 0) {
      console.log('[nextStep] errorFields 2', errorFields);
      // dispatch(setStepMissingFields({ step: currentStep, errors: errorFields }));
      console.log('[nextStep] errorFields currentStep', currentStep);
      dispatch(setStepStatus({ step: currentStep, status: LSTATUS.ERROR }));
    } else {
      dispatch(setStepStatus({ step: currentStep, status: LSTATUS.COMPLETE }));
    }
    dispatch(incrementStep());
  };

  const previousStep = (errorFields?: any) => {
    if (errorFields) {
      dispatch(setStepMissingFields({ step: currentStep, errors: errorFields }));
      dispatch(setStepStatus({ step: currentStep, status: LSTATUS.ERROR }));
    } else {
      dispatch(setStepStatus({ step: currentStep, status: LSTATUS.COMPLETE }));
    }
    dispatch(decrementStep());
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

  // save the registration value to the store
  const setRegistrationValue = (value: any) => {
    dispatch(setCertificateValue({ value }));
  };

  const setInitialState = (value: any) => {
    const state: TPayload = {
      currentStep: value.currentStep,
      steps: value.steps,
      lastStep: 6,
      hasReachSubmitStep: value.hasReachSubmitStep,
      testnetSubmitted: value.testnetSubmitted,
      mainnetSubmitted: value.mainnetSubmitted,
      data: value.data
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

  const clearCertificateStepper = () => {
    dispatch(clearStepper());
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
    clearCertificateStepper
  };
};

export default useCertificateStepper;
