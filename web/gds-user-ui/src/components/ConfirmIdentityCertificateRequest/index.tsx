import {
  Button,
  ButtonProps,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Stack,
  Text,
  useDisclosure
} from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import CheckboxFormControl from 'components/ui/CheckboxFormControl';
import { ReactNode } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Link } from 'react-router-dom';

export type ConfirmIdentityCertificateProps = { children: ReactNode } & ButtonProps;

function ConfirmIdentityCertificateModal({ children, ...rest }: ConfirmIdentityCertificateProps) {
  const { onClose, onOpen, isOpen } = useDisclosure();
  const methods = useForm({
    defaultValues: {
      agreed: false
    },
    mode: 'all'
  });
  const { register, getValues, watch } = methods;
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const agreed = watch('agreed');

  return (
    <FormProvider {...methods}>
      <form>
        <Modal isOpen={isOpen} onClose={onClose}>
          <ModalOverlay />
          <ModalContent border="1px solid">
            <ModalHeader mt={3} pb={1}>
              <Trans id="New X.509 Identity Certificate Request">
                New X.509 Identity Certificate Request
              </Trans>
            </ModalHeader>
            <ModalBody display="flex" flexDirection="column" gap={[4, 6]}>
              <Text>
                <Trans id="Requesting a new X.509 Identity Certificate will invalidate and revoke your current X.509 Identity Certificate.">
                  Requesting a new X.509 Identity Certificate will invalidate and revoke your
                  current X.509 Identity Certificate.
                </Trans>
              </Text>
              <Stack>
                <CheckboxFormControl controlId="agreed" {...register('agreed')} colorScheme="gray">
                  <Trans id="I acknowledge that requesting a new X.509 Identity Certificate will invalidate and revoke my organization’s current X.509 Identity Certificate.">
                    I acknowledge that requesting a new X.509 Identity Certificate will invalidate
                    and revoke my organization’s current X.509 Identity Certificate.
                  </Trans>
                </CheckboxFormControl>
              </Stack>
              <Text>
                You are required to re-confirm your organization’s profile with TRISA. Click next to
                proceed. You can cancel later.
              </Text>
            </ModalBody>

            <ModalFooter display="flex" flexDir="column" justifyContent="center" gap={2}>
              <Link to="/dashboard/certificate/registration">
                <Button
                  bg="orange"
                  _hover={{ bg: 'orange' }}
                  type="submit"
                  minW={150}
                  isDisabled={getValues().agreed === false}>
                  <Trans id="Next">Next</Trans>
                </Button>
              </Link>
              <Button variant="ghost" onClick={onClose}>
                <Trans id="Cancel">Cancel</Trans>
              </Button>
            </ModalFooter>
          </ModalContent>
        </Modal>
      </form>
      <Button bg="#55ACD8" color="#fff" onClick={onOpen} {...rest}>
        {children}
      </Button>
    </FormProvider>
  );
}

export default ConfirmIdentityCertificateModal;
