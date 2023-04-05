import { Button } from '@chakra-ui/react';
import { t, Trans } from '@lingui/macro';
import useGetStepStatusByKey from 'components/CertificateReview/useGetStepStatusByKey';
type StepButtonsProps = {
  handlePreviousStep: () => void;
  handleNextStep?: () => void;
  currentStep: number;
  isCurrentStepLastStep: boolean;
  handleResetForm: () => void;
  isDefaultValue: () => boolean;
};

function StepButtons({
  handlePreviousStep,
  currentStep,
  handleNextStep,
  isCurrentStepLastStep,
  handleResetForm,
  isDefaultValue
}: StepButtonsProps) {
  const isFirstStep = currentStep === 1;
  const { requiredMissingFields } = useGetStepStatusByKey();
  return (
    <>
      <Button onClick={handlePreviousStep} isDisabled={isFirstStep}>
        {isCurrentStepLastStep ? t`Previous` : t`Save & Previous`}
      </Button>
      {currentStep < 6 ? (
        <Button onClick={handleNextStep} variant="secondary">
          {isCurrentStepLastStep ? t`Next` : t`Save & Next`}
        </Button>
      ) : (
        <Button type="submit" variant="secondary" disabled={!!requiredMissingFields}>
          {isCurrentStepLastStep ? t`Next` : t`Save & Next`}
        </Button>
      )}
      <Button onClick={handleResetForm} isDisabled={isDefaultValue()}>
        <Trans id="Clear & Reset Form">Clear & Reset Form</Trans>
      </Button>
    </>
  );
}

export default StepButtons;
