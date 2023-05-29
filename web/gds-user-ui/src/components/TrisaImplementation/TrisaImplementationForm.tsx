import React, { useEffect } from 'react';
import { chakra } from '@chakra-ui/react';
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
  const { previousStep, nextStep, currentState, updateIsDirty } = useCertificateStepper();
  const { certificateStep, isFetchingCertificateStep } = useFetchCertificateStep({
    key: StepEnum.TRISA
  });
  const { updateCertificateStep, updatedCertificateStep } = useUpdateCertificateStep();

  const resolver = yupResolver(trisaImplementationValidationSchema);
  const methods = useForm({
    defaultValues: certificateStep?.form,
    resolver,
    mode: 'onChange'
  });

  const {
    formState: { isDirty }
  } = methods;

  useEffect(() => {
    updateIsDirty(isDirty, StepsIndexes.TRISA_IMPLEMENTATION);
  }, [isDirty, updateIsDirty]);

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
      previousStep(updatedCertificateStep);
    }
    previousStep(certificateStep);
  };

  const handleNextStepClick = () => {
    console.log('[] handleNextStep', methods.getValues());
    // if the form is dirty, then we need to save the data and move to the next step
    console.log('[] isDirty', isDirty);
    if (!isDirty) {
      console.log('[] is not Dirty', isDirty);
      nextStep(updatedCertificateStep ?? certificateStep);
    } else {
      const payload = {
        step: StepEnum.TRISA,
        form: {
          ...methods.getValues(),
          state: currentState()
        } as any
      };
      console.log('[] isDirty  payload', payload);

      updateCertificateStep(payload);
      console.log('[] isDirty 3 (not)', updatedCertificateStep);
      nextStep(updatedCertificateStep);
    }
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
            />
          </chakra.form>
          {!isProdEnv ? <DevTool control={methods.control} /> : null}
        </FormProvider>
      )}
    </FormLayout>
  );
};

export default TrisaForm;
