import { InfoIcon } from '@chakra-ui/icons';
import { Box, Heading, HStack, Icon, Stack, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { getSteps, getCurrentStep } from 'application/store/selectors/stepper';
import { SectionStatus } from 'components/SectionStatus';
import TrixoQuestionnaireForm from 'components/TrixoQuestionnaireForm';
import FormLayout from 'layouts/FormLayout';
import { useSelector } from 'react-redux';
import { getStepStatus } from 'utils/utils';

const TrixoQuestionnaire: React.FC = () => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);

  return (
    <Stack spacing={4} mt="2rem">
      <HStack>
        <Heading size="md">Section 5: TRIXO Questionnaire</Heading>
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
      <TrixoQuestionnaireForm />
    </Stack>
  );
};

export default TrixoQuestionnaire;
