import { Button } from '@chakra-ui/react';
import { t, Trans } from '@lingui/macro';

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
  return (
    <>
      <Button onClick={handlePreviousStep} isDisabled={isFirstStep}>
        {isCurrentStepLastStep ? t`Previous` : t`Save & Previous`}
      </Button>

      <Button onClick={handleNextStep} variant="secondary">
        {t`Save & Next`}
      </Button>

      <Button onClick={handleResetForm} isDisabled={isDefaultValue()}>
        <Trans id="Clear & Reset Form">Clear & Reset Form</Trans>
      </Button>
    </>
  );
}

export default StepButtons;
