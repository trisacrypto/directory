import React, { FC, useEffect } from 'react';
import { CertificateStepsProvider, useCertificateSteps } from 'contexts/certificateStepsContext';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import {
  addStep,
  setCurrentStep,
  setStepStatus,
  setLastStep,
  TStep
} from 'application/store/stepper.slice';
interface StepsProps {
  currentStep?: number;
  children: JSX.Element[];
}

const CertificateSteps: FC<StepsProps> = (props: any): any => {
  const currentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);

  return React.Children.map(props.children, (child: any, index: any) => {
    const isCurrentStep = +child.key === currentStep;
    if (child.props.isLast) setLastStep(+child.key);
    return React.cloneElement(child, {
      ...child.props,
      isCurrentStep
    });
  });
};

export default CertificateSteps;
