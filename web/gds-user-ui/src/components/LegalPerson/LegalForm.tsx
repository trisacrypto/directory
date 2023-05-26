import { chakra } from '@chakra-ui/react';
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
const LegalForm: React.FC = () => {
  const { previousStep, nextStep, currentState } = useCertificateStepper();
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

  const handleNextStepClick = () => {
    console.log('[] handleNextStep', methods.getValues());
    // if the form is dirty, then we need to save the data and move to the next step
    console.log('[] isDirty', isDirty);
    if (!isDirty) {
      console.log('[] is not Dirty', isDirty);
      nextStep(updatedCertificateStep?.errors ?? certificateStep?.errors);
    } else {
      const payload = {
        step: StepEnum.LEGAL,
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
      previousStep(updatedCertificateStep?.errors);
    }
    previousStep();
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
            />
          </chakra.form>
        </FormProvider>
      )}
    </FormLayout>
  );
};

export default LegalForm;
