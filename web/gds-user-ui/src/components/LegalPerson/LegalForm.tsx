import React, { useEffect, useState, useRef, Dispatch, SetStateAction } from 'react';
import { Box, chakra, useDisclosure } from '@chakra-ui/react';
import CountryOfRegistration from 'components/CountryOfRegistration';
import FormLayout from 'layouts/FormLayout';
import NameIdentifiers from '../NameIdentifiers';
import NationalIdentification from '../NationalIdentification';
import Address from 'components/Addresses';
import { FormProvider, useForm } from 'react-hook-form';
import StepButtons from 'components/StepsButtons';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { legalPersonValidationSchemam } from 'modules/dashboard/certificate/lib/legalPersonValidationSchema';
import { yupResolver } from '@hookform/resolvers/yup';

import { StepEnum } from 'types/enums';

import { useUpdateCertificateStep } from 'hooks/useUpdateCertificateStep';
import { StepsIndexes } from 'constants/steps';
import { isProdEnv } from 'application/config';
import { DevTool } from '@hookform/devtools';
interface LegalFormProps {
  data?: any;
  isLoading?: boolean;
  shouldResetForm?: boolean;
  onResetFormState?: Dispatch<SetStateAction<boolean>>;
  onNextClick?: () => void;
  onPreviousClick?: () => void;
}
const LegalForm: React.FC<LegalFormProps> = ({ data, shouldResetForm, onResetFormState }) => {
  const { isOpen, onClose, onOpen } = useDisclosure();
  const [shouldShowResetFormModal, setShouldShowResetFormModal] = useState(false);
  const { previousStep, nextStep, currentState, updateIsDirty } = useCertificateStepper();
  const {
    updateCertificateStep,
    updatedCertificateStep,
    reset: resetMutation,
    wasCertificateStepUpdated,
    isUpdatingCertificateStep
  } = useUpdateCertificateStep();
  const previousStepRef = useRef<any>(false);
  const nextStepRef = useRef<any>(false);
  const resolver = yupResolver(legalPersonValidationSchemam);
  const methods = useForm({
    defaultValues: data,
    resolver,
    mode: 'onChange'
  });

  const {
    formState: { isDirty },
    reset: resetForm
  } = methods;

  const onCloseModalHandler = () => {
    setShouldShowResetFormModal(false);
    onClose();
  };

  if (wasCertificateStepUpdated && nextStepRef.current) {
    resetMutation();
    // reset the form with the new values
    resetForm(updatedCertificateStep?.form, {
      keepValues: false
    });
    nextStep(updatedCertificateStep);
    nextStepRef.current = false;
  }

  if (wasCertificateStepUpdated && previousStepRef.current && !isUpdatingCertificateStep) {
    resetMutation();
    // reset the form with the new values
    resetForm(updatedCertificateStep?.form, {
      keepValues: false
    });
    // console.log('[] prev updatedCertificateStep', updatedCertificateStep);
    previousStepRef.current = false;
    previousStep(updatedCertificateStep);
  }

  const handleNextStepClick = () => {
    if (!isDirty) {
      nextStep({
        step: StepEnum.LEGAL,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      });
    } else {
      const payload = {
        step: StepEnum.LEGAL,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      };

      updateCertificateStep(payload);
      nextStepRef.current = true;
      // nextStep(updatedCertificateStep);
    }
  };

  const handlePreviousStepClick = () => {
    if (isDirty) {
      const payload = {
        step: StepEnum.LEGAL,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      };
      // console.log('[] isDirty  payload', payload);

      updateCertificateStep(payload);
      previousStepRef.current = true;
      // previousStep(updatedCertificateStep);
    }

    previousStep(data);
  };

  const handleResetForm = () => {
    setShouldShowResetFormModal(true); // this will show the modal
  };

  const handleResetClick = () => {
    setShouldShowResetFormModal(false); // this will close the modal
  };

  useEffect(() => {
    updateIsDirty(isDirty, StepsIndexes.LEGAL_PERSON);
  }, [isDirty, updateIsDirty]);

  useEffect(() => {
    if (shouldShowResetFormModal) {
      onOpen();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [shouldShowResetFormModal]);

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
          <NameIdentifiers />
          <Address />
          <CountryOfRegistration />
          <NationalIdentification />
          <Box pt={5}>
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
          </Box>
        </chakra.form>
        {!isProdEnv ? <DevTool control={methods.control} /> : null}
      </FormProvider>
    </FormLayout>
  );
};

export default LegalForm;
