import React, { useEffect } from 'react';
import { Heading, HStack, Stack, Text, Link } from '@chakra-ui/react';

import FormLayout from 'layouts/FormLayout';
import LegalForm from './LegalForm';
import { useSelector } from 'react-redux';
import { getCurrentStep, getSteps } from 'application/store/selectors/stepper';
import { getStepStatus } from 'utils/utils';
import { SectionStatus } from 'components/SectionStatus';
import { Trans } from '@lingui/react';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { StepEnum } from 'types/enums';
import MinusLoader from 'components/Loader/MinusLoader';
const LegalPerson: React.FC = () => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const [shouldResetForm, setShouldResetForm] = React.useState(false);
  const stepStatus = getStepStatus(steps, currentStep);
  const { certificateStep, isFetchingCertificateStep, getCertificateStep } =
    useFetchCertificateStep({
      key: StepEnum.LEGAL
    });
  const { isStepDeleted, updateDeleteStepState } = useCertificateStepper();
  const isLegalStepDeleted = isStepDeleted(StepEnum.LEGAL);

  useEffect(() => {
    if (isLegalStepDeleted) {
      const payload = {
        step: StepEnum.LEGAL,
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
    isLegalStepDeleted,
    shouldResetForm
  ]);

  // rerender this view everytime user land on this page
  useEffect(() => {
    getCertificateStep();
  }, [getCertificateStep]);

  return (
    <Stack spacing={7} mt="2rem">
      <HStack>
        <Heading size="md" pr={3} ml={2} data-cy="legal-person-form">
          <Trans id={'Section 2: Legal Person'}>Section 2: Legal Person</Trans>
        </Heading>
        {stepStatus ? <SectionStatus status={stepStatus} /> : null}
      </HStack>
      <FormLayout>
        <Text>
          <Trans id="Please enter the information that identifies your organization as a Legal Person. This form represents the">
            Please enter the information that identifies your organization as a Legal Person. This
            form represents the
          </Trans>{' '}
          <Link isExternal href="https://intervasp.org/" color={'blue'} fontWeight={'bold'}>
            {' '}
            <Trans id="IVMS 101 data structure">IVMS 101 data structure</Trans>
          </Link>{' '}
          <Trans
            id={
              'for legal persons and is strongly suggested for use as KYC or CDD (Know your Counterparty) information exchanged in TRISA transfers.'
            }>
            for legal persons and is strongly suggested for use as KYC or CDD (Know your
            Counterparty) information exchanged in TRISA transfers.
          </Trans>
        </Text>
      </FormLayout>
      {isFetchingCertificateStep ? (
        <MinusLoader />
      ) : (
        <LegalForm
          data={certificateStep?.form}
          shouldResetForm={shouldResetForm}
          onResetFormState={setShouldResetForm}
        />
      )}
    </Stack>
  );
};

export default LegalPerson;
