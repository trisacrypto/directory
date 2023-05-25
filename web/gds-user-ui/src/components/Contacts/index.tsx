import { Heading, HStack, Stack, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { getCurrentStep, getSteps } from 'application/store/selectors/stepper';
import { SectionStatus } from 'components/SectionStatus';
import FormLayout from 'layouts/FormLayout';
import { useSelector } from 'react-redux';
import { getStepStatus } from 'utils/utils';
import ContactsForm from './ContactsForm';

const Contacts: React.FC = () => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);

  return (
    <Stack spacing={13} mt="2rem">
      <HStack>
        <Heading size="md" pr={3} ml={2}>
          <Trans id="Section 3: Contacts">Section 3: Contacts</Trans>
        </Heading>
        {stepStatus ? <SectionStatus status={stepStatus} /> : null}
      </HStack>
      <FormLayout>
        <Text>
          <Trans id="Please supply contact information for representatives of your organization. All contacts will receive an email verification token and the contact emails must be verified before the registration can proceed.">
            Please supply contact information for representatives of your organization. All contacts
            will receive an email verification token and the contact emails must be verified before
            the registration can proceed.
          </Trans>
        </Text>
      </FormLayout>
      <ContactsForm />
    </Stack>
  );
};

export default Contacts;
