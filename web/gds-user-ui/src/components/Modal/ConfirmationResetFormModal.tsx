import React, { useState, useEffect } from 'react';
import {
  Box,
  chakra,
  Heading,
  Input,
  Text,
  Flex,
  useClipboard,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  useDisclosure,
  Button
} from '@chakra-ui/react';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { Trans } from '@lingui/react';
const ConfirmationResetForm = (props: any) => {
  const { isOpen: isAlertOpen, onOpen: onAlertOpen, onClose: onAlertClose } = useDisclosure();
  const { resetForm } = useCertificateStepper();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const handleOnClose = () => {
    props.onClose();
    onAlertClose();
    props.onChangeState(false);
  };
  const handleResetBtn = () => {
    setIsLoading(true);
    // props.onReset(loadDefaultValueFromLocalStorage);
    resetForm();
    props.onChangeResetState(true);
    props.onChangeState(false);
    setIsLoading(false);
    props.onClose();
    onAlertClose();
    // props.onRefeshState();
  };
  return (
    <>
      <Flex>
        <Box w="full">
          <Modal closeOnOverlayClick={false} {...props}>
            <ModalOverlay />
            <ModalContent width={'100%'}>
              <ModalHeader data-testid="confirmation-modal-header" textAlign={'center'}>
                <Trans id="Clear & Reset Registration Form">Clear & Reset Registration Form</Trans>
              </ModalHeader>

              <ModalBody pb={5}>
                <Text pb={2} fontSize={'sm'}>
                  <Trans id="Click “Reset” to clear and reset the registration form. All data will be deleted and you will be re-directed to the beginning of the form and you will be required to restart the registration process">
                    Click “Reset” to clear and reset the registration form. All data will be deleted
                    and you will be re-directed to the beginning of the form and you will be
                    required to restart the registration process
                  </Trans>
                </Text>
              </ModalBody>

              <ModalFooter textAlign={'center'} justifyContent={'center'}>
                <Button
                  mr={10}
                  onClick={handleResetBtn}
                  isLoading={isLoading}
                  bgColor="#23a7e0e8"
                  color="#fff"
                  _hover={{
                    bgColor: '#189fda'
                  }}>
                  <Trans id="Reset">Reset</Trans>
                </Button>
                <Button onClick={handleOnClose} bgColor="#555151" color={'#fff'}>
                  <Trans id="Cancel">Cancel</Trans>
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
