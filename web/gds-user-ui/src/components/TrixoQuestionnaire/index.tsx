import React, { useEffect, useState } from 'react';
import { Heading, HStack, Stack, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { getSteps, getCurrentStep } from 'application/store/selectors/stepper';
import { SectionStatus } from 'components/SectionStatus';
import TrixoQuestionnaireForm from 'components/TrixoQuestionnaireForm';
import FormLayout from 'layouts/FormLayout';
import { useSelector } from 'react-redux';
import { getStepStatus } from 'utils/utils';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { StepEnum } from 'types/enums';
import MinusLoader from 'components/Loader/MinusLoader';
const TrixoQuestionnaire: React.FC = () => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);
  const [shouldResetForm, setShouldResetForm] = useState<boolean>(false);
  const { isStepDeleted, updateDeleteStepState } = useCertificateStepper();
  const isTrixoStepDeleted = isStepDeleted(StepEnum.TRIXO);
  const { certificateStep, isFetchingCertificateStep, getCertificateStep } =
    useFetchCertificateStep({
      key: StepEnum.TRIXO
    });
  useEffect(() => {
    if (isTrixoStepDeleted) {
      const payload = {
        step: StepEnum.TRIXO,
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
    isTrixoStepDeleted,
    shouldResetForm
  ]);

  useEffect(() => {
    getCertificateStep();
  }, [getCertificateStep]);

  return (
    <Stack spacing={4} mt="2rem">
      <HStack>
        <Heading size="md" pr={3} ml={2} data-cy="trixo-form">
          <Trans id="Section 5: TRIXO Questionnaire">Section 5: TRIXO Questionnaire</Trans>
        </Heading>
        {stepStatus ? <SectionStatus status={stepStatus} /> : null}
      </HStack>
      <FormLayout>
        <Text>
          <Trans id="This questionnaire is designed to help TRISA members understand the regulatory regime of your organization. The information provided will help ensure that required compliance information exchanges are conducted correctly and safely. All verified TRISA members will have access to this information.">
            This questionnaire is designed to help TRISA members understand the regulatory regime of
            your organization. The information provided will help ensure that required compliance
            information exchanges are conducted correctly and safely. All verified TRISA members
            will have access to this information.
          </Trans>
        </Text>
      </FormLayout>
      {isFetchingCertificateStep ? (
        <MinusLoader />
      ) : (
        <TrixoQuestionnaireForm
          data={certificateStep?.form}
          shouldResetForm={shouldResetForm}
          onResetFormState={setShouldResetForm}
        />
      )}
    </Stack>
  );
};

export default TrixoQuestionnaire;
