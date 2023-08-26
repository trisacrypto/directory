import React, { useState, useEffect } from 'react';
import {
  Stack,
  HStack,
  Heading,
  Text,
  Box,
  Button,
  useColorModeValue,
  useDisclosure
} from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import FormLayout from 'layouts/FormLayout';
import BasicDetailsReview from './BasicDetailsReview';
import ContactsReview from './ContactsReview';
import LegalPersonReview from './LegalPersonReview';
import TrisaImplementationReview from './TrisaImplementationReview';
import TrixoReview from './TrixoReview';
import { CgExport } from 'react-icons/cg';
import StepButtons from 'components/StepsButtons';
import { downloadRegistrationData } from 'modules/dashboard/registration/utils';
import { handleError } from 'utils/utils';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { StepEnum } from 'types/enums';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
import MinusLoader from 'components/Loader/MinusLoader';
const ReviewsSummary: React.FC = () => {
  const { isOpen, onClose, onOpen } = useDisclosure();
  // const [shouldReload, setShouldReload] = useState(false);
  const [shouldShowResetFormModal, setShouldShowResetFormModal] = useState(false);
  const { previousStep, updateHasReachSubmitStep, updateCurrentStepState } =
    useCertificateStepper();
  const [isLoadingExport, setIsLoadingExport] = useState(false);
  const { certificateStep, getCertificateStep, isFetchingCertificateStep } =
    useFetchCertificateStep({
      key: StepEnum.ALL
    });

  const handleExport = () => {
    const downloadData = async () => {
      try {
        setIsLoadingExport(true);
        await downloadRegistrationData();
      } catch (error) {
        handleError(error, 'Error while downloading registration data');
      } finally {
        setIsLoadingExport(false);
      }
    };
    downloadData();
  };

  const handleNextStep = () => {
    updateHasReachSubmitStep(true);
  };

  const handlePreviousStep = () => {
    previousStep();
  };

  const isNextButtonDisabled = certificateStep?.errors?.length > 0;

  const handleResetForm = () => {
    setShouldShowResetFormModal(true); // this will show the modal
  };

  const handleResetClick = () => {
    setShouldShowResetFormModal(false); // this will close the modal
  };

  const onCloseModalHandler = () => {
    setShouldShowResetFormModal(false);
    onClose();
    updateCurrentStepState(StepEnum.BASIC);
  };

  // if certificate step is 6 then reload the page

  useEffect(() => {
    if (shouldShowResetFormModal) {
      onOpen();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [shouldShowResetFormModal]);

  // fetch certificate step everytime user land on this page
  useEffect(() => {
    getCertificateStep();
  }, [getCertificateStep]);

  return (
    <Stack spacing={7}>
      <HStack pt={10} justifyContent={'space-between'}>
        <Heading size="md" data-testid="review" data-cy="review-page">
          <Trans id="Review">Review</Trans>
        </Heading>
        <Box>
          <Button
            bg={useColorModeValue('black', 'white')}
            _hover={{
              bg: useColorModeValue('black', 'white')
            }}
            color={useColorModeValue('white', 'black')}
            onClick={handleExport}
            isLoading={isLoadingExport}
            leftIcon={<CgExport />}>
            <Trans id="Export Data">Export Data</Trans>
          </Button>
        </Box>
      </HStack>
      <FormLayout>
        <Text>
          <Trans id="Please review the information provided, edit as needed, and submit to complete the registration form. After the information is reviewed, you will be contacted to verify details. Once verified, your TestNet certificate will be issued.">
            Please review the information provided, edit as needed, and submit to complete the
            registration form. After the information is reviewed, you will be contacted to verify
            details. Once verified, your TestNet certificate will be issued.
          </Trans>
        </Text>
      </FormLayout>
      {isFetchingCertificateStep ? (
        <Box>
          <MinusLoader />
        </Box>
      ) : (
        <>
          <BasicDetailsReview />
          <LegalPersonReview />
          <ContactsReview />
          <TrisaImplementationReview />
          <TrixoReview />

          <StepButtons
            handleNextStep={handleNextStep}
            handlePreviousStep={handlePreviousStep}
            isNextButtonDisabled={isNextButtonDisabled}
            onResetModalClose={handleResetClick}
            isOpened={isOpen}
            handleResetForm={handleResetForm}
            resetFormType={StepEnum.ALL}
            onClosed={onCloseModalHandler}
            handleResetClick={handleResetClick}
            shouldShowResetFormModal={shouldShowResetFormModal}
          />
        </>
      )}
    </Stack>
  );
};

export default ReviewsSummary;
