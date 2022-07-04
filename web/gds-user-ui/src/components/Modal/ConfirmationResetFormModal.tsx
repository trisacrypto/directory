import React from 'react';
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
const ConfirmationResetForm = (props: any) => {
  const { isOpen: isAlertOpen, onOpen: onAlertOpen, onClose: onAlertClose } = useDisclosure();
  const { resetForm } = useCertificateStepper();
  const handleOnClose = () => {
    props.onClose();
    onAlertClose();
  };
  const handleResetBtn = () => {
    resetForm();
    props.onClose();
    onAlertClose();
  };
  return (
    <>
      <Flex>
        <Box w="full">
          <Modal closeOnOverlayClick={false} {...props}>
            <ModalOverlay />
            <ModalContent width={'100%'}>
              <ModalHeader data-testid="confirmation-modal-header" textAlign={'center'}>
                Clear & Reset Registration Form
              </ModalHeader>

              <ModalBody pb={5}>
                <Text pb={2} fontSize={'sm'}>
                  Click “Reset” to clear and reset the registration form. All data will be deleted
                  and you will be re-directed to the beginning of the form and you will be required
                  to restart the registration process
                </Text>
              </ModalBody>

              <ModalFooter textAlign={'center'} justifyContent={'center'}>
                <Button
                  mr={10}
                  onClick={handleResetBtn}
                  bgColor="#23a7e0e8"
                  color="#fff"
                  _hover={{
                    bgColor: '#189fda'
                  }}>
                  Reset
                </Button>
                <Button onClick={handleOnClose} bgColor="#555151" color={'#fff'}>
                  Cancel
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
