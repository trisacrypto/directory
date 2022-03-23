import React, { FC, useEffect } from 'react';
import { CertificateStepsProvider, useCertificateSteps } from 'contexts/certificateStepsContext';

interface StepsProps {
  currentStep?: number;
  children: JSX.Element[];
}

const CertificateSteps: FC<StepsProps> = (props: any): any => {
  const [certificateSteps, setCertificateSteps] = useCertificateSteps();

  return React.Children.map(props.children, (child: any, index: any) => {
    const isCurrentStep = +child.key === certificateSteps.currentStep;
    return React.cloneElement(child, {
      ...child.props,
      isCurrentStep
    });
  });
};

export default CertificateSteps;
