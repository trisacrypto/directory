import React, { useEffect, useState } from 'react';
import { Heading, HStack, Stack, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { getSteps, getCurrentStep } from 'application/store/selectors/stepper';
import { SectionStatus } from 'components/SectionStatus';
import FormLayout from 'layouts/FormLayout';
// import { useEffect } from 'react';
import { useSelector } from 'react-redux';
import { getStepStatus } from 'utils/utils';
import TrisaForm from 'components/TrisaImplementation/TrisaImplementationForm';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { StepEnum } from 'types/enums';
import MinusLoader from 'components/Loader/MinusLoader';

const TrisaImplementation: React.FC = () => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);

  const { certificateStep, isFetchingCertificateStep, getCertificateStep } =
    useFetchCertificateStep({
      key: StepEnum.TRISA
    });
  const [shouldResetForm, setShouldResetForm] = useState<boolean>(false);
  const { isStepDeleted, updateDeleteStepState } = useCertificateStepper();
  const isTrisaStepDeleted = isStepDeleted(StepEnum.TRISA);

  useEffect(() => {
    if (isTrisaStepDeleted) {
      const payload = {
        step: StepEnum.TRISA,
        isDeleted: false
      };
      updateDeleteStepState(payload);
      getCertificateStep();
      setShouldResetForm(true);
    }
    return () => {
      setShouldResetForm(false);
    };
  }, [
    isStepDeleted,
    updateDeleteStepState,
    getCertificateStep,
    isTrisaStepDeleted,
    shouldResetForm
  ]);

  // rerender this view everytime user land on this page
  useEffect(() => {
    getCertificateStep();
  }, [getCertificateStep]);

  return (
    <Stack spacing={7} pt={8}>
      <HStack>
        <Heading size="md" pr={3} ml={2} data-cy="trisa-form">
          <Trans id="Section 4: TRISA Implementation">Section 4: TRISA Implementation</Trans>
        </Heading>
        {stepStatus ? <SectionStatus status={stepStatus} /> : null}
      </HStack>
      <FormLayout>
        <Text>
          <Trans id="Each VASP is required to establish a TRISA endpoint for inter-VASP communication. Please specify the details of your endpoint for certificate issuance. Please specify the TestNet endpoint and the MainNet endpoint. The TestNet endpoint and the MainNet endpoint must be different.">
            Each VASP is required to establish a TRISA endpoint for inter-VASP communication. Please
            specify the details of your endpoint for certificate issuance. Please specify the
            TestNet endpoint and the MainNet endpoint. The TestNet endpoint and the MainNet endpoint
            must be different.
          </Trans>
        </Text>
      </FormLayout>
      {isFetchingCertificateStep ? (
        <MinusLoader />
      ) : (
        <TrisaForm
          data={certificateStep?.form}
          shouldResetForm={shouldResetForm}
          onResetFormState={setShouldResetForm}
        />
      )}
    </Stack>
  );
};

export default TrisaImplementation;
