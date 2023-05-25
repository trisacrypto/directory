import { FormProvider, useForm } from 'react-hook-form';
import { t } from '@lingui/macro';
import { chakra } from '@chakra-ui/react';
import StepButtons from 'components/StepsButtons';
import ContactForm from 'components/Contacts/ContactForm';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { yupResolver } from '@hookform/resolvers/yup';
import { contactsValidationSchema } from 'modules/dashboard/certificate/lib/contactsValidationSchema';
import { StepEnum } from 'types/enums';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
import { useUpdateCertificateStep } from 'hooks/useUpdateCertificateStep';
import FormLayout from 'layouts/FormLayout';
import MinusLoader from 'components/Loader/MinusLoader';

interface ContactsFormProps {
    data?: any;
    isLoading?: boolean;
}

const ContactsForm: React.FC<ContactsFormProps> = ({data}) => {
    const { previousStep, nextStep, currentState } = useCertificateStepper();
    const { certificateStep, isFetchingCertificateStep } = useFetchCertificateStep({
    key: StepEnum.CONTACTS
    });

    const { updateCertificateStep, updatedCertificateStep } = useUpdateCertificateStep();

    const resolver = yupResolver(contactsValidationSchema);

    const methods = useForm({
    defaultValues: certificateStep?.form || data,
    resolver,
    mode: 'onChange'
  });

    const {
    formState: { isDirty }
    } = methods;

    const handleNextStepClick = () => {
        console.log('[] handleNextStep', methods.getValues());
        // if the form is dirty, then we need to save the data and move to the next step
        console.log('[] isDirty', isDirty);
        if (!isDirty) {
          console.log('[] is not Dirty', isDirty);
          nextStep(updatedCertificateStep?.errors ?? certificateStep?.errors);
        } else {
          const payload = {
            step: StepEnum.CONTACTS,
            form: {
              ...methods.getValues(),
              state: currentState()
            } as any
          };
          console.log('[] isDirty  payload', payload);
    
          updateCertificateStep(payload);
          console.log('[] isDirty 3 (not)', updatedCertificateStep);
          nextStep(updatedCertificateStep?.errors);
        }
      };
    

    return (
        <FormLayout>
            {isFetchingCertificateStep ? (
            <MinusLoader />
            ) : (
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
    )}
       </FormLayout>
    );
};

export default ContactsForm;
