import React, { FC, useEffect } from 'react';
import { Collapse } from '@chakra-ui/transition';
interface StepLabelProps {
  key: string;
  component: JSX.Element;
  isCurrentStep?: boolean;
  isLast?: boolean;
}

const CertificateStepContainer: FC<StepLabelProps> = (props) => {
  return (
    <>
      <Collapse in={props.isCurrentStep}>{props.component}</Collapse>
    </>
  );
}; // ProgressBar

export default CertificateStepContainer;
