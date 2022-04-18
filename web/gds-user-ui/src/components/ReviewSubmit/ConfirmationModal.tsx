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
  ModalCloseButton,
  useDisclosure,
  Button
} from '@chakra-ui/react';
import ModalAlert from 'components/ui/ModalAlert';
interface ConfirmationModalProps {}
const ConfirmationModal = (props: any) => {
  const { onClose } = useDisclosure();
  const { hasCopied, onCopy } = useClipboard(props.pkcs12password);
  const [isAlerted, setIsAlerted] = React.useState(false);
  const handleOnClose = () => {
    // should ask user to confirm before closing
    onClose();

    // const result = window.prompt('Copy to clipboard: Ctrl+C, Enter', pkcs12password);
  };
  return (
    <>
      <Flex>
        <Box w="full">
          <Modal closeOnOverlayClick={false} {...props}>
            <ModalOverlay />
            <ModalContent width={'100%'}>
              <ModalHeader textAlign={'center'}>TRISA Registration Request Submitted!</ModalHeader>
              <ModalCloseButton />
              <ModalBody pb={6}>
                <Text pb={5} fontSize={'sm'}>
                  Your registration request has been successfully received by the Directory Service.
                  Verification emails have been sent to all contacts listed. Once your contact
                  information has been verified, the registration form will be sent to the TRISA
                  review board to verify your membership in the TRISA network.
                </Text>
                <Text pb={2} fontSize={'sm'}>
                  When you are verified you will be issued PKCS12 encrypted identity certificates
                  for use in mTLS authentication between TRISA members. The password to decrypt
                  those certificates is shown below:
                </Text>
                <Text>
                  <Flex mb={2} fontSize={'sm'}>
                    <Input
                      value={props.pkcs12password}
                      isReadOnly
                      bg={!hasCopied ? 'yellow.100' : 'green.200'}
                    />
                    <Button onClick={onCopy} ml={2}>
                      {hasCopied ? 'Copied' : 'Copy'}
                    </Button>
                  </Flex>
                </Text>
                <Text py={2} color={'orange.500'} fontSize={'sm'}>
                  This is the only time the PKCS12 password is shown during the registration
                  process.
                  <br />
                  Please copy and paste this password and store somewhere safe!
                </Text>
                <Box py={2} fontSize={'sm'}>
                  <chakra.tr>
                    <chakra.td>
                      <Text fontWeight={'semibold'}>ID :</Text>
                    </chakra.td>
                    <chakra.td>
                      <Text pl={2}>{props.id}</Text>
                    </chakra.td>
                  </chakra.tr>
                  <chakra.tr>
                    <chakra.td>
                      <Text fontWeight={'semibold'}>Verification Status : </Text>
                    </chakra.td>
                    <chakra.td>
                      <Text pl={2}>{props.status}</Text>
                    </chakra.td>
                  </chakra.tr>
                </Box>
                <Text>
                  <Text as={'span'} fontWeight={'semibold'}>
                    Message from server:
                  </Text>{' '}
                  {props.message?.toUpperCase()}
                </Text>
              </ModalBody>

              <ModalFooter>
                <Button onClick={onClose}>Understood</Button>
              </ModalFooter>
            </ModalContent>
          </Modal>
        </Box>
      </Flex>
    </>
  );
};

export default ConfirmationModal;
