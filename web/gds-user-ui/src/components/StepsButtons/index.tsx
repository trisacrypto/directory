import { Button, Stack, useDisclosure } from '@chakra-ui/react';

import { t, Trans } from '@lingui/macro';
import ConfirmationResetFormModal from 'components/Modal/ConfirmationResetFormModal';
type StepButtonsProps = {
  handlePreviousStep?: () => void;
  handleNextStep?: () => void;
  currentStep?: number;
  isCurrentStepLastStep?: boolean;
  handleResetForm?: () => void;
  isDefaultValue?: () => boolean;
  isFirstStep?: boolean;
  isNextButtonDisabled?: boolean;
  resetFormType?: string;
  shouldShowResetFormModal?: boolean;
};

function StepButtons({
  handlePreviousStep,
  isFirstStep = false,
  isNextButtonDisabled = false,
  handleNextStep,
  isCurrentStepLastStep,
  handleResetForm,
  shouldShowResetFormModal = false,
  // resetFormType = 'section',
  isDefaultValue = () => false
}: StepButtonsProps) {
  const { isOpen, onClose } = useDisclosure();
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
        <Button onClick={handlePreviousStep} isDisabled={isFirstStep}>
          {isCurrentStepLastStep ? t`Previous` : t`Save & Previous`}
        </Button>

        <Button onClick={handleNextStep} variant="secondary" isDisabled={isNextButtonDisabled}>
          {t`Save & Next`}
        </Button>

        <Button onClick={handleResetForm} isDisabled={isDefaultValue()}>
          <Trans id="Clear & Reset Form">Clear & Reset Form</Trans>
        </Button>
        {shouldShowResetFormModal && (
          <ConfirmationResetFormModal isOpen={isOpen} onClose={onClose} />
        )}
      </Stack>
    </>
  );
}

export default StepButtons;
