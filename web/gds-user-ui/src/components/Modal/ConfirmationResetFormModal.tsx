import { useEffect } from 'react';
import {
  Box,
  Text,
  Flex,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  useDisclosure,
  Button
} from '@chakra-ui/react';
import { useDeleteCertificateStep } from 'hooks/useDeleteCertificateStep';
import useCertificateStepper from 'hooks/useCertificateStepper';
// import { useNavigate } from 'react-router-dom';
// import useCertificateStepper from 'hooks/useCertificateStepper';
import { Trans } from '@lingui/macro';
import { StepperType } from 'types/type';
import { getStepNumber } from 'components/BasicDetailsForm/util';
import { StepEnum } from 'types/enums';
const ConfirmationResetForm = (props: any) => {
  const { step } = props;
  const { deleteCertificateStep, isDeletingCertificateStep, wasCertificateStepDeleted } =
    useDeleteCertificateStep(step);
  const { updateStepStatusState } = useCertificateStepper();
  // const navigate = useNavigate();
  const { onClose: onAlertClose } = useDisclosure();
  // const { resetForm } = useCertificateStepper();

  const handleOnClose = () => {
    props.onClose();
    onAlertClose();
  };

  const isResetAllType = props.step === StepEnum.ALL;

  const handleResetBtn = () => {
    deleteCertificateStep({
      key: props.step as StepperType
    });

    // resetForm();
    // props.onChangeResetState(true);
    // props.onChangeState(false);

    // navigate('/dashboard/certificate/registration');
  };

  useEffect(() => {
    if (wasCertificateStepDeleted) {
      updateStepStatusState({
        step: getStepNumber(props.step),
        status: 'incomplete'
      });
      props.onClose();
    }
  }, [wasCertificateStepDeleted, props, updateStepStatusState]);

  return (
    <>
      <Flex>
        <Box w="full">
          <Modal closeOnOverlayClick={false} {...props} isOpen={props.isOpen}>
            <ModalOverlay />
            <ModalContent width={'100%'}>
              <ModalHeader data-testid="confirmation-modal-header" textAlign={'center'}>
                {isResetAllType ? (
                  <Trans>Clear & Reset Registration Form</Trans>
                ) : (
                  <Trans>Clear & Reset Section</Trans>
                )}
              </ModalHeader>

              <ModalBody pb={5}>
                <Text pb={2} fontSize={'sm'}>
                  {isResetAllType ? (
                    <Trans>
                      Click “Reset” to clear and{' '}
                      <Text as="span" fontWeight={'bold'}>
                        reset the form.
                      </Text>{' '}
                      reset the registration form. All data will be deleted and you will be
                      re-directed to the beginning of the form and you will be required to restart
                      the registration process
                    </Trans>
                  ) : (
                    <Trans>
                      Click “Reset” to clear and reset{' '}
                      <Text as="span" fontWeight={'bold'}>
                        {' '}
                        this section{' '}
                      </Text>{' '}
                      of the form. All data will be deleted and cleared from the registration form.
                      You may want to export the data first. After clearing, you will be taken to
                      the start of the form.
                    </Trans>
                  )}
                </Text>
              </ModalBody>

              <ModalFooter textAlign={'center'} justifyContent={'center'}>
                <Button
                  mr={10}
                  onClick={handleResetBtn}
                  isLoading={isDeletingCertificateStep}
                  bgColor="#23a7e0e8"
                  color="#fff"
                  _hover={{
                    bgColor: '#189fda'
                  }}>
                  <Trans>Reset</Trans>
                </Button>
                <Button
                  onClick={handleOnClose}
                  bgColor="#555151"
                  color={'#fff'}
                  _hover={{ boxShadow: '#555151', bgColor: '#555151D4' }}>
                  <Trans>Cancel</Trans>
                </Button>
              </ModalFooter>
            </ModalContent>
          </Modal>
        </Box>
      </Flex>
    </>
  );
};

export default ConfirmationResetForm;
