import React, { useEffect } from 'react';
import { Heading, HStack, Stack, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { getCurrentStep, getSteps } from 'application/store/selectors/stepper';
import { SectionStatus } from 'components/SectionStatus';
import FormLayout from 'layouts/FormLayout';
import { useSelector } from 'react-redux';
import { getStepStatus } from 'utils/utils';
import ContactsForm from './ContactsForm';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { StepEnum } from 'types/enums';
import MinusLoader from 'components/Loader/MinusLoader';
const Contacts: React.FC = () => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);
  const { certificateStep, isFetchingCertificateStep, getCertificateStep } =
    useFetchCertificateStep({
      key: StepEnum.CONTACTS
    });
  const [shouldResetForm, setShouldResetForm] = React.useState(false);
  const { isStepDeleted, updateDeleteStepState } = useCertificateStepper();
  const isContactsStepDeleted = isStepDeleted(StepEnum.CONTACTS);

  useEffect(() => {
    if (isContactsStepDeleted) {
      console.log('[] isContactsStepDeleted', isContactsStepDeleted);
      const payload = {
        step: StepEnum.CONTACTS,
        isDeleted: false
      };
      updateDeleteStepState(payload);
      getCertificateStep();
      setShouldResetForm(true);
    }
  }, [
    isStepDeleted,
    updateDeleteStepState,
    getCertificateStep,
    isContactsStepDeleted,
    shouldResetForm
  ]);
  return (
    <Stack spacing={13} mt="2rem">
      <HStack>
        <Heading size="md" pr={3} ml={2} data-cy="contacts-form">
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
      {isFetchingCertificateStep ? (
        <MinusLoader />
      ) : (
        <ContactsForm
          data={certificateStep?.form}
          shouldResetForm={shouldResetForm}
          onResetFormState={setShouldResetForm}
        />
      )}
    </Stack>
  );
};

export default Contacts;
