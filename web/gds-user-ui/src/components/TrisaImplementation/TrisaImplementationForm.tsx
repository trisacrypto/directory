import React, { useEffect, useState, useRef } from 'react';
import { chakra, useDisclosure } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import FormLayout from 'layouts/FormLayout';
import { FormProvider, useForm } from 'react-hook-form';
import { trisaImplementationValidationSchema } from 'modules/dashboard/certificate/lib/trisaImplementationValidationSchema';
import { yupResolver } from '@hookform/resolvers/yup';
import useCertificateStepper from 'hooks/useCertificateStepper';
import StepButtons from 'components/StepsButtons';
import MinusLoader from 'components/Loader/MinusLoader';
import { StepEnum } from 'types/enums';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
import { useUpdateCertificateStep } from 'hooks/useUpdateCertificateStep';
import TrisaImplementationForm from './TrisaImplementationForm/index';
import { StepsIndexes } from 'constants/steps';
import { isProdEnv } from 'application/config';
import { DevTool } from '@hookform/devtools';
const TrisaForm: React.FC = () => {
  const { isOpen, onClose, onOpen } = useDisclosure();
  const [shouldShowResetFormModal, setShouldShowResetFormModal] = useState(false);
  const { previousStep, nextStep, currentState, updateIsDirty } = useCertificateStepper();
  const { certificateStep, isFetchingCertificateStep } = useFetchCertificateStep({
    key: StepEnum.TRISA
  });
  const {
    updateCertificateStep,
    updatedCertificateStep,
    wasCertificateStepUpdated,
    isUpdatingCertificateStep,
    reset: resetMutation
  } = useUpdateCertificateStep();
  const previousStepRef = useRef<any>(false);
  const nextStepRef = useRef<any>(false);
  const resolver = yupResolver(trisaImplementationValidationSchema);
  const methods = useForm({
    defaultValues: certificateStep?.form,
    resolver,
    mode: 'onChange'
  });

  const {
    formState: { isDirty },
    reset: resetForm
  } = methods;

  useEffect(() => {
    updateIsDirty(isDirty, StepsIndexes.TRISA_IMPLEMENTATION);
  }, [isDirty, updateIsDirty]);

  useEffect(() => {
    if (shouldShowResetFormModal) {
      onOpen();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [shouldShowResetFormModal]);

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
    console.log('[] prev updatedCertificateStep', updatedCertificateStep);
    previousStepRef.current = false;
    previousStep(updatedCertificateStep);
  }

  const handlePreviousStepClick = () => {
    if (isDirty) {
      const payload = {
        step: StepEnum.TRISA,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      };
      updateCertificateStep(payload);
      previousStepRef.current = true;
    }
    previousStep(certificateStep);
  };

  const handleNextStepClick = () => {
    if (!isDirty) {
      nextStep(certificateStep);
    } else {
      const payload = {
        step: StepEnum.TRISA,
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
  const onCloseModalHandler = () => {
    setShouldShowResetFormModal(false);
    onClose();
  };

  return (
    <FormLayout>
      {isFetchingCertificateStep ? (
        <MinusLoader />
      ) : (
        <FormProvider {...methods}>
          <chakra.form
            onSubmit={methods.handleSubmit(handleNextStepClick)}
            data-testid="trisa-implementation-form">
            <TrisaImplementationForm
              type="TestNet"
              name="testnet"
              headerText={t`TRISA Endpoint: TestNet`}
            />
            <TrisaImplementationForm
              type="MainNet"
              name="mainnet"
              headerText={t`TRISA Endpoint: MainNet`}
            />
            <StepButtons
              handlePreviousStep={handlePreviousStepClick}
              handleNextStep={handleNextStepClick}
              onResetModalClose={handleResetClick}
              isOpened={isOpen}
              handleResetForm={handleResetForm}
              resetFormType={StepEnum.TRISA}
              onClosed={onCloseModalHandler}
              handleResetClick={handleResetClick}
              shouldShowResetFormModal={shouldShowResetFormModal}
            />
          </chakra.form>
          {!isProdEnv ? <DevTool control={methods.control} /> : null}
        </FormProvider>
      )}
    </FormLayout>
  );
};

export default TrisaForm;
