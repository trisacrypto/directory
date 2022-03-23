import React, { FC, useEffect } from 'react';
import { HStack, Box, Icon, Text, Heading, Stack, Grid, Button } from '@chakra-ui/react';

import { Collapse } from '@chakra-ui/transition';
import { useCertificateSteps } from 'contexts/certificateStepsContext';
import {
  CertificateStepContainer,
  CertificateStepLabel,
  CertificateSteps
} from './CertificateStepper';

const ProgressBar = () => {
  const [certificateSteps, setCertificateSteps] = useCertificateSteps();
  useEffect(() => {
    setCertificateSteps({
      currentStep: 2,
      steps: [
        {
          key: 2,
          status: 'Progress'
        }
      ]
    });
  }, [setCertificateSteps]);

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

      {/* <HStack>
        <Button>Next</Button>
        <Button>Previous</Button>
      </HStack> */}
    </>
    // </CertificateStepsProvider>
  );
};

export default ProgressBar;
