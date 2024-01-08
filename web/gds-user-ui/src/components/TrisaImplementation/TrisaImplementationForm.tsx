import React, { useEffect, useState, useRef, Dispatch, SetStateAction } from 'react';
import { Box, chakra, useDisclosure } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import FormLayout from 'layouts/FormLayout';
import { FormProvider, useForm } from 'react-hook-form';
import { trisaImplementationValidationSchema } from 'modules/dashboard/certificate/lib/trisaImplementationValidationSchema';
import { yupResolver } from '@hookform/resolvers/yup';
import useCertificateStepper from 'hooks/useCertificateStepper';
import StepButtons from 'components/StepsButtons';
import { StepEnum } from 'types/enums';
import { useUpdateCertificateStep } from 'hooks/useUpdateCertificateStep';
import TrisaImplementationForm from './TrisaImplementationForm/index';
import { StepsIndexes } from 'constants/steps';
import { isProdEnv } from 'application/config';
import { DevTool } from '@hookform/devtools';
interface TrisaFormProps {
  data?: any;
  isLoading?: boolean;
  shouldResetForm?: boolean;
  onResetFormState?: Dispatch<SetStateAction<boolean>>;
}
const TrisaForm: React.FC<TrisaFormProps> = ({ data, shouldResetForm, onResetFormState }) => {
  const { isOpen, onClose, onOpen } = useDisclosure();
  const [shouldShowResetFormModal, setShouldShowResetFormModal] = useState(false);
  const { previousStep, nextStep, currentState, updateIsDirty } = useCertificateStepper();
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
    defaultValues: data,
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
    resetForm(updatedCertificateStep?.form);
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
    previousStep(data);
  };

  const handleNextStepClick = () => {
    if (!isDirty) {
      nextStep({
        step: StepEnum.TRISA,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      });
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
        <chakra.form
          onSubmit={methods.handleSubmit(handleNextStepClick)}
          data-testid="trisa-implementation-form">
          <TrisaImplementationForm
            type="TestNet"
            name="testnet"
            headerText={t`TRISA Endpoint: TestNet`}
          />
          <Box pt={5}>
            <TrisaImplementationForm
            type="MainNet"
            name="mainnet"
            headerText={t`TRISA Endpoint: MainNet`}
            />
          </Box>
          <Box pt={5}>
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
          </Box>
        </chakra.form>
        {!isProdEnv ? <DevTool control={methods.control} /> : null}
      </FormProvider>
    </FormLayout>
  );
};

export default TrisaForm;
