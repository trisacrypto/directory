import { Button, Stack } from '@chakra-ui/react';
import { t, Trans } from '@lingui/macro';

type StepButtonsProps = {
  handlePreviousStep?: () => void;
  handleNextStep?: () => void;
  currentStep?: number;
  isCurrentStepLastStep?: boolean;
  handleResetForm?: () => void;
  isDefaultValue?: () => boolean;
  isFirstStep?: boolean;
  isNextButtonDisabled?: boolean;
};

function StepButtons({
  handlePreviousStep,
  isFirstStep = false,
  isNextButtonDisabled = false,
  handleNextStep,
  isCurrentStepLastStep,
  handleResetForm,
  isDefaultValue = () => false
}: StepButtonsProps) {
  return (
    <>
      <Stack
        width="100%"
        direction={'row'}
        spacing={8}
        justifyContent={'center'}
        py={6}
        wrap="wrap"
        data-testid="step-buttons"
        rowGap={2}>
        <Button onClick={handlePreviousStep} isDisabled={isFirstStep}>
          {isCurrentStepLastStep ? t`Previous` : t`Save & Previous`}
        </Button>

        <Button onClick={handleNextStep} variant="secondary" isDisabled={isNextButtonDisabled}>
          {t`Save & Next`}
        </Button>

        <Button onClick={handleResetForm} isDisabled={isDefaultValue()}>
          <Trans id="Clear & Reset Form">Clear & Reset Form</Trans>
        </Button>
      </Stack>
    </>
  );
}

export default StepButtons;
