import React, { useEffect, useState } from 'react';
import { chakra, useDisclosure } from '@chakra-ui/react';
import CountryOfRegistration from 'components/CountryOfRegistration';
import FormLayout from 'layouts/FormLayout';
import NameIdentifiers from '../NameIdentifiers';
import NationalIdentification from '../NameIdentification';
import Address from 'components/Addresses';
import { FormProvider, useForm } from 'react-hook-form';
import StepButtons from 'components/StepsButtons';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { legalPersonValidationSchemam } from 'modules/dashboard/certificate/lib/legalPersonValidationSchema';
import { yupResolver } from '@hookform/resolvers/yup';
import MinusLoader from 'components/Loader/MinusLoader';
import { StepEnum } from 'types/enums';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
import { useUpdateCertificateStep } from 'hooks/useUpdateCertificateStep';
import { StepsIndexes } from 'constants/steps';
const LegalForm: React.FC = () => {
  const { isOpen, onClose, onOpen } = useDisclosure();
  const [shouldShowResetFormModal, setShouldShowResetFormModal] = useState(false);
  const { previousStep, nextStep, currentState, updateIsDirty } = useCertificateStepper();
  const { certificateStep, isFetchingCertificateStep } = useFetchCertificateStep({
    key: StepEnum.LEGAL
  });
  const { updateCertificateStep, updatedCertificateStep } = useUpdateCertificateStep();

  const resolver = yupResolver(legalPersonValidationSchemam);
  const methods = useForm({
    defaultValues: certificateStep?.form,
    resolver,
    mode: 'onChange'
  });

  const {
    formState: { isDirty }
  } = methods;

  useEffect(() => {
    updateIsDirty(isDirty, StepsIndexes.LEGAL_PERSON);
  }, [isDirty, updateIsDirty]);

  const onCloseModalHandler = () => {
    setShouldShowResetFormModal(false);
    onClose();
  };

  const handleNextStepClick = () => {
    if (!isDirty) {
      nextStep(updatedCertificateStep ?? certificateStep);
    } else {
      const payload = {
        step: StepEnum.LEGAL,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      };

      updateCertificateStep(payload);
      nextStep(updatedCertificateStep);
    }
  };
  useEffect(() => {
    if (shouldShowResetFormModal) {
      onOpen();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [shouldShowResetFormModal]);

  const handlePreviousStepClick = () => {
    if (isDirty) {
      const payload = {
        step: StepEnum.LEGAL,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      };
      console.log('[] isDirty  payload', payload);

      updateCertificateStep(payload);
      previousStep(updatedCertificateStep);
    }
    previousStep(certificateStep);
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
            <NameIdentifiers />
            <Address />
            <CountryOfRegistration />
            <NationalIdentification />
            <StepButtons
              handlePreviousStep={handlePreviousStepClick}
              handleNextStep={handleNextStepClick}
              onResetModalClose={handleResetClick}
              isOpened={isOpen}
              handleResetForm={handleResetForm}
              resetFormType={StepEnum.LEGAL}
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

export default LegalForm;
