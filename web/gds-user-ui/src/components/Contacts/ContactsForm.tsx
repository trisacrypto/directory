import React, { useEffect, useState } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { t } from '@lingui/macro';
import { chakra, useDisclosure } from '@chakra-ui/react';
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
import { StepsIndexes } from 'constants/steps';
const ContactsForm: React.FC = () => {
  const { isOpen, onClose, onOpen } = useDisclosure();
  const [shouldShowResetFormModal, setShouldShowResetFormModal] = useState(false);
  const { previousStep, nextStep, currentState, updateIsDirty } = useCertificateStepper();
  const { certificateStep, isFetchingCertificateStep } = useFetchCertificateStep({
    key: StepEnum.CONTACTS
  });

  const { updateCertificateStep, updatedCertificateStep } = useUpdateCertificateStep();

  const resolver = yupResolver(contactsValidationSchema);

  const methods = useForm({
    defaultValues: certificateStep?.form,
    resolver,
    mode: 'onChange'
  });

  const {
    formState: { isDirty }
  } = methods;

  useEffect(() => {
    updateIsDirty(isDirty, StepsIndexes.CONTACTS);
  }, [isDirty, updateIsDirty]);

  const handlePreviousStepClick = () => {
    if (isDirty) {
      const payload = {
        step: StepEnum.CONTACTS,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      };
      updateCertificateStep(payload);
      previousStep(updatedCertificateStep);
    }
    previousStep(certificateStep);
  };

  useEffect(() => {
    if (shouldShowResetFormModal) {
      onOpen();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [shouldShowResetFormModal]);

  const onCloseModalHandler = () => {
    setShouldShowResetFormModal(false);
    onClose();
  };

  const handleNextStepClick = () => {
    if (!isDirty) {
      nextStep(updatedCertificateStep ?? certificateStep);
    } else {
      const payload = {
        step: StepEnum.CONTACTS,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      };

      updateCertificateStep(payload);
      nextStep(updatedCertificateStep);
    }
  };
  const handleResetForm = () => {
    setShouldShowResetFormModal(true); // this will show the modal
  };

  const handleResetClick = () => {
    setShouldShowResetFormModal(false); // this will close the modal
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
            <StepButtons
              handlePreviousStep={handlePreviousStepClick}
              handleNextStep={handleNextStepClick}
              onResetModalClose={handleResetClick}
              isOpened={isOpen}
              handleResetForm={handleResetForm}
              resetFormType={StepEnum.CONTACTS}
              onClosed={onCloseModalHandler}
              handleResetClick={handleResetClick}
              shouldShowResetFormModal={shouldShowResetFormModal}
            />
          </chakra.form>
        </FormProvider>
      )}
    </FormLayout>
  );
};

export default ContactsForm;
