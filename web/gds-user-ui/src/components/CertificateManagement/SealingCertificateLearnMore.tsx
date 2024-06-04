import {
  useDisclosure,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
  Text,
  ButtonProps
} from '@chakra-ui/react';
import { ReactNode } from 'react';

export type SealingCertificateLearnMoreProps = {
  children: ReactNode;
} & ButtonProps;

function SealingCertificateLearnMore({ children, ...rest }: SealingCertificateLearnMoreProps) {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <>
      <Modal blockScrollOnMount={false} isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader textAlign="center">Sealing Certificates</ModalHeader>
          <ModalBody>
            <Text mb="1rem">
              TRISA will soon have the ability to issue sealing certificates for members. Sealing
              certificates are useful for the long-term data storage of Secure Envelopes.
            </Text>
            <Text>
              If your organization requires a sealing certificate now, TRISA can issue them
              manually. Please contact support@rotational.io to learn more.
            </Text>
          </ModalBody>

          <ModalFooter>
            <Button bg="#555151D4" _hover={{ boxShadow: '#555151D9' }} onClick={onClose} mx="auto">
              Close
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
      <Button onClick={onOpen} {...rest}>
        {children}
      </Button>
    </>
  );
}

export default SealingCertificateLearnMore;
