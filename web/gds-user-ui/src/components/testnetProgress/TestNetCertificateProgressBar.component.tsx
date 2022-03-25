import React, { FC, useEffect } from 'react';
import { HStack, Box, Icon, Text, Heading, Stack, Grid, Button } from '@chakra-ui/react';

import { Collapse } from '@chakra-ui/transition';
import useCertificateStepper from 'hooks/useCertificateStepper';
import {
  CertificateStepContainer,
  CertificateStepLabel,
  CertificateSteps
} from './CertificateStepper';

const ProgressBar = () => {
  const { nextStep, previousStep } = useCertificateStepper();

  return (
    <>
      <CertificateSteps>
        <CertificateStepLabel />
        <CertificateStepContainer
          key="1"
          status="progress"
          component={<Text> component 1 </Text>}
        />
        <CertificateStepContainer
          key="2"
          status="progress"
          component={<Text> component 2 </Text>}
        />
        <CertificateStepContainer
          key="3"
          status="progress"
          component={<Text> component 3 </Text>}
        />
        <CertificateStepContainer
          key="4"
          status="progress"
          component={<Text> component 4 </Text>}
        />
      </CertificateSteps>

      <HStack>
        <Button onClick={() => nextStep()}>Next</Button>
        <Button onClick={() => previousStep()}>Previous</Button>
      </HStack>
    </>
    // </CertificateStepsProvider>
  );
};

export default ProgressBar;
