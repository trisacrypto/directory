import { Heading, HStack, Stack, Text, chakra } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';
import { getCurrentStep, getSteps } from 'application/store/selectors/stepper';
import ContactForm from 'components/Contacts/ContactForm';
import { SectionStatus } from 'components/SectionStatus';
import FormLayout from 'layouts/FormLayout';
import { useSelector } from 'react-redux';
import { getStepStatus } from 'utils/utils';
import { contactsValidationSchema } from 'modules/dashboard/certificate/lib/contactsValidationSchema';
import StepButtons from 'components/StepsButtons';
// import { ContactsForm } from './ContactsForm'; // group all contact forms in one component
import { yupResolver } from '@hookform/resolvers/yup';
import { FormProvider, useForm } from 'react-hook-form';
import useCertificateStepper from 'hooks/useCertificateStepper';

const Contacts: React.FC = () => {
  const steps = useSelector(getSteps);
  const currentStep = useSelector(getCurrentStep);
  const stepStatus = getStepStatus(steps, currentStep);

  const { previousStep, nextStep } = useCertificateStepper();

  const resolver = yupResolver(contactsValidationSchema);
  const methods = useForm({
    defaultValues: {},
    resolver,
    mode: 'onChange'
  });

  const handleNextStepClick = () => {
    console.log('[] handleNextStep', methods.getValues());
    nextStep();
  };

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
      <FormProvider {...methods}>
        <chakra.form onSubmit={methods.handleSubmit(handleNextStepClick)}>
          <ContactForm
            name={`contacts.legal`}
            title={t`Legal/ Compliance Contact (required)`}
            description={t`Compliance officer or legal contact for requests about the compliance requirements and legal status of your organization. A business phone number is required to complete physical verification for MainNet registration. Please provide a phone number where the Legal/ Compliance contact can be contacted.`}
          />
          <ContactForm
            name="contacts.technical"
            title={t`Technical Contact (required)`}
            description={t`Primary contact for handling technical queries about the operation and status of your service participating in the TRISA network. Can be a group or admin email.`}
          />
          <ContactForm
            name="contacts.administrative"
            title={t`Administrative Contact (optional)`}
            description={t`Administrative or executive contact for your organization to field high-level requests or queries. (Strongly recommended)`}
          />
          <ContactForm
            name="contacts.billing"
            title={t`Billing Contact (optional)`}
            description={t`Billing contact for your organization to handle account and invoice requests or queries relating to the operation of the TRISA network.`}
          />
          <StepButtons handlePreviousStep={previousStep} handleNextStep={handleNextStepClick} />
        </chakra.form>
      </FormProvider>
    </Stack>
  );
};

export default Contacts;
