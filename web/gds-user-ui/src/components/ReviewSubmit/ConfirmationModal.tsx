import React from 'react';
import {
  Box,
  chakra,
  Alert,
  Tag,
  AlertTitle,
  AlertDescription,
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
import ModalAlert from 'components/ReviewSubmit/ModalAlert';
import AlertContent from './AlertContent';
import { upperCaseFirstLetter } from 'utils/utils';
import { Trans } from '@lingui/react';
const ConfirmationModal = (props: any) => {
  const { isOpen: isAlertOpen, onOpen: onAlertOpen, onClose: onAlertClose } = useDisclosure();
  const { hasCopied, onCopy } = useClipboard(props.pkcs12password);
  const [isAlerted, setIsAlerted] = React.useState(false);
  const handleOnClose = () => {
    // should ask user to confirm before closing
    setIsAlerted(true);
    onAlertOpen();

    // const result = window.prompt('Copy to clipboard: Ctrl+C, Enter', pkcs12password);
  };
  const handleYesBtn = () => {
    // should ask user to confirm before closing
    props.onClose();
    onAlertClose();

    // const result = window.prompt('Copy to clipboard: Ctrl+C, Enter', pkcs12password);
  };
  return (
    <>
      <Flex>
        <Box w="full">
          <Modal closeOnOverlayClick={false} {...props}>
            <ModalOverlay />
            <ModalContent width={'100%'}>
              <ModalHeader data-testid="confirmation-modal-header" textAlign={'center'}>
                <Trans id="TRISA Registration Request Submitted!">
                  TRISA Registration Request Submitted!
                </Trans>
              </ModalHeader>

              <ModalBody pb={6}>
                <Text pb={5} fontSize={'sm'}>
                  <Trans id="Your registration request has been successfully received by the Directory Service. Verification emails have been sent to all contacts listed. Once your contact information has been verified, the registration form will be sent to the TRISA review board to verify your membership in the TRISA network.">
                    Your registration request has been successfully received by the Directory
                    Service. Verification emails have been sent to all contacts listed. Once your
                    contact information has been verified, the registration form will be sent to the
                    TRISA review board to verify your membership in the TRISA network.{' '}
                  </Trans>
                </Text>
                <Text pb={2} fontSize={'sm'}>
                  <Trans
                    id=" When you are verified you will be issued PKCS12 encrypted identity certificates
                  for use in mTLS authentication between TRISA members. The password to decrypt
                  those certificates is shown below:">
                    When you are verified you will be issued PKCS12 encrypted identity certificates
                    for use in mTLS authentication between TRISA members. The password to decrypt
                    those certificates is shown below:
                  </Trans>
                </Text>
                <Text>
                  <Flex mb={2} fontSize={'sm'}>
                    <Input
                      value={props.pkcs12password}
                      isReadOnly
                      bg={!hasCopied ? 'yellow.100' : 'green.200'}
                    />
                    <Button onClick={onCopy} ml={2}>
                      {hasCopied ? (
                        <Trans id="Copied">Copied</Trans>
                      ) : (
                        <Trans id="Copy">Copy</Trans>
                      )}
                    </Button>
                  </Flex>
                </Text>
                <Text py={2} color={'orange.500'} fontSize={'sm'}>
                  <Trans id="This is the only time the PKCS12 password is shown during the registration process.">
                    This is the only time the PKCS12 password is shown during the registration
                    process.
                  </Trans>
                  <br />
                  <Trans id="Please copy and paste this password and store somewhere safe!">
                    Please copy and paste this password and store somewhere safe!
                  </Trans>
                </Text>
                <Box py={2} fontSize={'sm'}>
                  <chakra.tr>
                    <chakra.td>
                      <Text fontWeight={'semibold'}>
                        <Trans id="ID:">ID:</Trans>
                      </Text>
                    </chakra.td>
                    <chakra.td>
                      <Text ml={5}>{props.id}</Text>
                    </chakra.td>
                  </chakra.tr>
                  <chakra.tr>
                    <chakra.td>
                      <Text fontWeight={'semibold'}>
                        <Trans id="Verification Status:">Verification Status:</Trans>
                      </Text>
                    </chakra.td>
                    <chakra.td>
                      <Tag ml={5} bg={'green'} color={'white'}>
                        {props.status}
                      </Tag>
                    </chakra.td>
                  </chakra.tr>
                </Box>
                <Box mt={5}>
                  <Alert status="info">
                    <Box>
                      <AlertTitle>
                        <Trans id="Message from server:">Message from server:</Trans>
                      </AlertTitle>
                      <AlertDescription>{upperCaseFirstLetter(props.message)}</AlertDescription>
                    </Box>
                  </Alert>
                </Box>
              </ModalBody>

              <ModalFooter>
                <Button onClick={handleOnClose}>
                  <Trans id="Understood">Understood</Trans>
                </Button>
              </ModalFooter>
            </ModalContent>
          </Modal>
          {isAlerted && (
            <ModalAlert
              header={'Confirm'}
              message={<AlertContent />}
              handleYesBtn={handleYesBtn}
              isOpen={isAlertOpen}
              onOpen={onAlertOpen}
              onClose={onAlertClose}
            />
          )}
        </Box>
      </Flex>
    </>
  );
};

export default ConfirmationModal;
