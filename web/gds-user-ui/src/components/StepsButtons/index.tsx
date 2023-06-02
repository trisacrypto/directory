import { Button, Stack } from '@chakra-ui/react';
import { StepEnum } from 'types/enums';
import { t, Trans } from '@lingui/macro';
import ConfirmationResetFormModal from 'components/Modal/ConfirmationResetFormModal';
type StepButtonsProps = {
  handlePreviousStep?: () => void;
  handleNextStep?: () => void;
  handleResetClick?: () => void;
  currentStep?: number;
  isCurrentStepLastStep?: boolean;
  handleResetForm?: () => void;
  isDefaultValue?: () => boolean;
  isFirstStep?: boolean;
  isOpened?: boolean;
  onClosed?: () => void;
  isNextButtonDisabled?: boolean;
  resetFormType?: string;
  onResetModalClose?: () => void;
  shouldShowResetFormModal?: boolean;
};

function StepButtons({
  handlePreviousStep,
  isFirstStep = false,
  isNextButtonDisabled = false,
  handleNextStep,
  isCurrentStepLastStep,
  handleResetForm,
  shouldShowResetFormModal,
  resetFormType,
  isOpened,
  onClosed,
  onResetModalClose,
  isDefaultValue = () => false
}: StepButtonsProps) {
  // const [isResetModalOpen, setIsResetModalOpen] = useState<boolean>(false);
  // const onChangeModalState = (value: boolean) => {
  //   setIsResetModalOpen(value);
  // };
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
        <Button onClick={handlePreviousStep} isDisabled={isFirstStep} data-cy="previous-bttn">
          {isCurrentStepLastStep ? t`Previous` : t`Save & Previous`}
        </Button>

        <Button onClick={handleNextStep} variant="secondary" isDisabled={isNextButtonDisabled} data-cy="next-bttn">
          {t`Save & Next`}
        </Button>

        <Button onClick={handleResetForm} isDisabled={isDefaultValue()} data-cy="clear-reset-bttn">
          {resetFormType !== StepEnum.ALL ? (
            <Trans>Clear & Reset Section</Trans>
          ) : (
            <Trans>Clear & Reset Form</Trans>
          )}
        </Button>
      </Stack>

      {shouldShowResetFormModal && (
        <ConfirmationResetFormModal
          isOpen={isOpened}
          onClose={onClosed}
          step={resetFormType}
          onReset={onResetModalClose}
          resetType={resetFormType}
        />
      )}
    </>
  );
}

export default StepButtons;
