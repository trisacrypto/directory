import React, { useEffect, useState, useRef, SetStateAction, Dispatch } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { t } from '@lingui/macro';
import { Box, chakra, useDisclosure } from '@chakra-ui/react';
import StepButtons from 'components/StepsButtons';
import ContactForm from 'components/Contacts/ContactForm';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { yupResolver } from '@hookform/resolvers/yup';
import { contactsValidationSchema } from 'modules/dashboard/certificate/lib/contactsValidationSchema';
import { StepEnum } from 'types/enums';

import { useUpdateCertificateStep } from 'hooks/useUpdateCertificateStep';
import FormLayout from 'layouts/FormLayout';

import { StepsIndexes } from 'constants/steps';
import { isProdEnv } from 'application/config';
import { DevTool } from '@hookform/devtools';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
interface ContactsFormProps {
  data?: any;
  isLoading?: boolean;
  shouldResetForm?: boolean;
  onResetFormState?: Dispatch<SetStateAction<boolean>>;
}
const ContactsForm: React.FC<ContactsFormProps> = ({ data, shouldResetForm, onResetFormState }) => {
  const { isOpen, onClose, onOpen } = useDisclosure();
  const [shouldShowResetFormModal, setShouldShowResetFormModal] = useState(false);
  const { previousStep, nextStep, currentState, updateIsDirty } = useCertificateStepper();
  const { certificateStep } = useFetchCertificateStep({
    key: StepEnum.CONTACTS
  });

  const {
    updateCertificateStep,
    updatedCertificateStep,
    wasCertificateStepUpdated,
    isUpdatingCertificateStep,
    reset
  } = useUpdateCertificateStep();
  const previousStepRef = useRef<any>(false);
  const nextStepRef = useRef<any>(false);
  const resolver = yupResolver(contactsValidationSchema);

  const methods = useForm({
    defaultValues: data,
    resolver,
    mode: 'onChange'
  });

  const {
    formState: { isDirty },
    reset: resetForm
  } = methods;

  useEffect(() => {
    updateIsDirty(isDirty, StepsIndexes.CONTACTS);
  }, [isDirty, updateIsDirty]);

  if (wasCertificateStepUpdated && nextStepRef.current) {
    reset();
    // reset the form with the new values
    resetForm(updatedCertificateStep?.form, {
      keepValues: false
    });
    nextStep(updatedCertificateStep);
    nextStepRef.current = false;
  }

  if (wasCertificateStepUpdated && previousStepRef.current && !isUpdatingCertificateStep) {
    reset();
    // reset the form with the new values
    resetForm(updatedCertificateStep?.form, {
      keepValues: false
    });
    previousStepRef.current = false;
    previousStep(updatedCertificateStep);
  }

  const handlePreviousStepClick = () => {
    if (!isDirty) {
      previousStep(certificateStep);
    } else {
      const payload = {
        step: StepEnum.CONTACTS,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      };
      updateCertificateStep(payload);
      previousStepRef.current = true;
    }
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
      nextStep(certificateStep);
    } else {
      const payload = {
        step: StepEnum.CONTACTS,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      };

      updateCertificateStep(payload);
      nextStepRef.current = true;
    }
  };
  const handleResetForm = () => {
    setShouldShowResetFormModal(true); // this will show the modal
  };

  const handleResetClick = () => {
    setShouldShowResetFormModal(false); // this will close the modal
  };

  // reset the form from the parent component
  useEffect(() => {
    if (shouldResetForm && onResetFormState) {
      resetForm(data);
      onResetFormState(false);
      window.location.reload();
    }
  }, [shouldResetForm, resetForm, data, onResetFormState]);

  return (
    <FormLayout>
      <FormProvider {...methods}>
        <chakra.form onSubmit={methods.handleSubmit(handleNextStepClick)}>
          <ContactForm
            name={`contacts.legal`}
            title={t`Legal/ Compliance Contact (required)`}
            description={t`Compliance officer or legal contact for requests about the compliance requirements and legal status of your organization. A business phone number is required to complete physical verification for MainNet registration. Please provide a phone number where the Legal/ Compliance contact can be contacted.`}
          />
          <Box pt={5}>
            <ContactForm
            name="contacts.technical"
            title={t`Technical Contact (required)`}
            description={t`Primary contact for handling technical queries about the operation and status of your service participating in the TRISA network. Can be a group or admin email.`}
            />
          </Box>
          <Box pt={5}>
            <ContactForm
            name="contacts.administrative"
            title={t`Administrative Contact (optional)`}
            description={t`Administrative or executive contact for your organization to field high-level requests or queries. (Strongly recommended)`}
            />
          </Box>
          <Box pt={5}>
            <ContactForm
            name="contacts.billing"
            title={t`Billing Contact (optional)`}
            description={t`Billing contact for your organization to handle account and invoice requests or queries relating to the operation of the TRISA network.`}
            />
          </Box>
          <Box pt={5}>
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
          </Box>
        </chakra.form>
        {!isProdEnv ? <DevTool control={methods.control} /> : null}
      </FormProvider>
    </FormLayout>
  );
};

export default ContactsForm;
