import React, { FC, useEffect } from 'react';
import { CertificateStepsProvider, useCertificateSteps } from 'contexts/certificateStepsContext';
import { useDispatch, useSelector, RootStateOrAny } from 'react-redux';
import { addStep, setCurrentStep, setStepStatus, TStep } from 'application/store/stepper.slice';
interface StepsProps {
  currentStep?: number;
  children: JSX.Element[];
}

const CertificateSteps: FC<StepsProps> = (props: any): any => {
  const CurrentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);

  return React.Children.map(props.children, (child: any, index: any) => {
    const isCurrentStep = +child.key === CurrentStep;
    return React.cloneElement(child, {
      ...child.props,
      isCurrentStep
    });
  });
};

export default CertificateSteps;
