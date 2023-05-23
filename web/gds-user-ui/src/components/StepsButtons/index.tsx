import { Button, Stack } from '@chakra-ui/react';
import { t, Trans } from '@lingui/macro';

type StepButtonsProps = {
  handlePreviousStep?: () => void;
  handleNextStep?: () => void;
  currentStep?: number;
  isCurrentStepLastStep?: boolean;
  handleResetForm?: () => void;
  isDefaultValue?: () => boolean;
};

function StepButtons({
  handlePreviousStep,
  currentStep,
  handleNextStep,
  isCurrentStepLastStep,
  handleResetForm,
  isDefaultValue = () => false
}: StepButtonsProps) {
  const isFirstStep = currentStep === 1;
  return (
    <>
      <Stack
        width="100%"
        direction={'row'}
        spacing={8}
        justifyContent={'center'}
        py={6}
        wrap="wrap"
        rowGap={2}>
        <Button onClick={handlePreviousStep} isDisabled={isFirstStep}>
          {isCurrentStepLastStep ? t`Previous` : t`Save & Previous`}
        </Button>
        {currentStep === 6 ? (
          <Button onClick={handleNextStep} variant="secondary">
            {t`Save & Next`}
          </Button>
        ) : (
          <Button type="submit" variant="secondary">
            {isCurrentStepLastStep ? t`Next` : t`Save & Next`}
          </Button>
        )}
        <Button onClick={handleResetForm} isDisabled={isDefaultValue()}>
          <Trans id="Clear & Reset Form">Clear & Reset Form</Trans>
        </Button>
      </Stack>
    </>
  );
}

export default StepButtons;
